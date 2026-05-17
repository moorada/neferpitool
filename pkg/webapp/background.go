package webapp

import (
	"fmt"
	"sync"
	"time"

	"github.com/moorada/neferpitool/pkg/app"
	"github.com/moorada/neferpitool/pkg/changes"
	"github.com/moorada/neferpitool/pkg/configuration"
	"github.com/moorada/neferpitool/pkg/db"
	"github.com/moorada/neferpitool/pkg/log"
	"github.com/moorada/neferpitool/pkg/notification"
	"github.com/moorada/neferpitool/pkg/reliableChanges"
	"github.com/robfig/cron/v3"
)

type backgroundStatus struct {
	Running       bool      `json:"running"`
	LastRun       time.Time `json:"last_run"`
	LastChanges   int       `json:"last_changes"`
	LastScanErrs  int       `json:"last_scan_errors"`
	LastError     string    `json:"last_error"`
	ActiveCrons   int       `json:"active_crons"`
	NextSleepMins int       `json:"next_sleep_minutes"`
}

type backgroundManager struct {
	mu      sync.Mutex
	service *app.Service
	cron    *cron.Cron
	stopCh  chan struct{}
	doneCh  chan struct{}

	running      bool
	lastRun      time.Time
	lastChanges  int
	lastScanErrs int
	lastError    string
}

func newBackgroundManager(service *app.Service) *backgroundManager {
	return &backgroundManager{service: service}
}

func (b *backgroundManager) Start() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.running {
		return nil
	}

	b.stopCh = make(chan struct{})
	b.doneCh = make(chan struct{})
	b.cron = cron.New()
	b.lastError = ""

	for _, expression := range configuration.GetConf().REPORTFREQUENCY {
		expr := expression
		_, err := b.cron.AddFunc(expr, func() {
			if err := b.sendReportEmail(expr); err != nil {
				log.Error("report email error for cron %s: %s", expr, err.Error())
			}
		})
		if err != nil {
			b.lastError = fmt.Sprintf("cron setup error (%s): %s", expr, err.Error())
		}
	}
	b.cron.Start()

	b.running = true
	go b.loop()
	return nil
}

func (b *backgroundManager) Stop() {
	b.mu.Lock()
	if !b.running {
		b.mu.Unlock()
		return
	}
	close(b.stopCh)
	doneCh := b.doneCh
	b.mu.Unlock()

	<-doneCh
}

func (b *backgroundManager) RunOnce() error {
	changesList, scanErrs, err := b.runCycle()
	b.mu.Lock()
	b.lastRun = time.Now()
	b.lastChanges = len(changesList)
	b.lastScanErrs = len(scanErrs)
	if err != nil {
		b.lastError = err.Error()
	} else if len(scanErrs) > 0 {
		b.lastError = fmt.Sprintf("%d typodomain scan errors", len(scanErrs))
	} else {
		b.lastError = ""
	}
	b.mu.Unlock()
	return err
}

func (b *backgroundManager) Status() backgroundStatus {
	b.mu.Lock()
	defer b.mu.Unlock()

	conf := configuration.GetConf()
	sleepMins := conf.MINUTESLEEPBACKGROUNDMONITORING
	if sleepMins <= 0 {
		sleepMins = 1
	}
	activeCrons := 0
	if b.cron != nil {
		activeCrons = len(b.cron.Entries())
	}

	return backgroundStatus{
		Running:       b.running,
		LastRun:       b.lastRun,
		LastChanges:   b.lastChanges,
		LastScanErrs:  b.lastScanErrs,
		LastError:     b.lastError,
		ActiveCrons:   activeCrons,
		NextSleepMins: sleepMins,
	}
}

func (b *backgroundManager) loop() {
	sleep := backgroundSleepDuration()
	ticker := time.NewTicker(sleep)
	defer ticker.Stop()

	for {
		select {
		case <-b.stopCh:
			b.mu.Lock()
			if b.cron != nil {
				ctx := b.cron.Stop()
				<-ctx.Done()
			}
			b.running = false
			close(b.doneCh)
			b.mu.Unlock()
			return
		case <-ticker.C:
			changesList, scanErrs, err := b.runCycle()

			b.mu.Lock()
			b.lastRun = time.Now()
			b.lastChanges = len(changesList)
			b.lastScanErrs = len(scanErrs)
			if err != nil {
				b.lastError = err.Error()
			} else if len(scanErrs) > 0 {
				b.lastError = fmt.Sprintf("%d typodomain scan errors", len(scanErrs))
			} else {
				b.lastError = ""
			}
			b.mu.Unlock()
		}
	}
}

func backgroundSleepDuration() time.Duration {
	conf := configuration.GetConf()
	if conf.MINUTESLEEPBACKGROUNDMONITORING <= 0 {
		return time.Minute
	}
	return time.Duration(conf.MINUTESLEEPBACKGROUNDMONITORING) * time.Minute
}

func (b *backgroundManager) runCycle() (changes.ChangeList, map[string]error, error) {
	changesList, scanErrs := b.service.CheckChangesForAll(nil)
	if len(changesList) == 0 {
		return changesList, scanErrs, nil
	}
	if err := b.sendChangesEmail(changesList); err != nil {
		return changesList, scanErrs, err
	}
	return changesList, scanErrs, nil
}

func (b *backgroundManager) sendChangesEmail(changesList changes.ChangeList) error {
	conf := configuration.GetConf()
	if conf.EMAIL == "" || conf.PASSWORD == "" || len(conf.EMAILTONOTIFY) == 0 {
		return nil
	}

	tdsInExpiration := b.service.GetTypoDomainsInExpiration()
	headersStatus, datasStatus, headersWhois, datasWhois := changesList.ToTables()
	hExpiry, dExpiry := tdsInExpiration.ToExpiryTable()

	request := notification.Request{
		From:     conf.EMAIL,
		Password: conf.PASSWORD,
		To:       conf.EMAILTONOTIFY,
		Subject:  "domain monitoring - Updates",
	}

	tpl := notification.TemplateData{
		H1:            "Domains Monitoring",
		TextStatus:    "There are status changes",
		TextWhois:     "There are whois changes",
		HeadersStatus: headersStatus,
		HeadersWhois:  headersWhois,
		DatasStatus:   datasStatus,
		DatasWhois:    datasWhois,
		TextExpiry:    "Typodomains in expiration",
		HeadersExpiry: hExpiry,
		DatasExpiry:   dExpiry,
	}

	return notification.EmailChanges(tpl, request)
}

func (b *backgroundManager) sendReportEmail(expression string) error {
	conf := configuration.GetConf()
	if conf.EMAIL == "" || conf.PASSWORD == "" || len(conf.EMAILTONOTIFY) == 0 {
		return nil
	}

	err, reliableChangesReport := db.GetRelaibleChangesFromDBWithExpression(expression)
	if err != nil {
		return err
	}

	tdsInExpiration := b.service.GetTypoDomainsInExpiration()
	headersStatus, datasStatus, headersWhois, datasWhois := reliableChangesReport.ToTables()
	hExpiry, dExpiry := tdsInExpiration.ToExpiryTable()

	request := notification.Request{
		From:     conf.EMAIL,
		Password: conf.PASSWORD,
		To:       conf.EMAILTONOTIFY,
		Subject:  "domain monitoring - Report - Expression: " + expression,
	}

	tpl := notification.TemplateData{
		H1:            "Domains Monitoring",
		TextStatus:    "There are status changes",
		TextWhois:     "There are whois changes",
		HeadersStatus: headersStatus,
		HeadersWhois:  headersWhois,
		DatasStatus:   datasStatus,
		DatasWhois:    datasWhois,
		TextExpiry:    "Typodomains in expiration",
		HeadersExpiry: hExpiry,
		DatasExpiry:   dExpiry,
	}

	if err := notification.EmailReport(tpl, request); err != nil {
		return err
	}

	for _, c := range reliableChangesReport {
		i := reliableChanges.Contains(c.Crons, expression)
		if i >= 0 {
			db.DeleteExprToChDB(c, *c.Crons[i])
		}
	}
	return nil
}
