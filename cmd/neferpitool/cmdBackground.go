package cmd

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/moorada/neferpitool/pkg/changes"
	"github.com/moorada/neferpitool/pkg/configuration"
	"github.com/moorada/neferpitool/pkg/db"
	"github.com/moorada/neferpitool/pkg/log"
	"github.com/moorada/neferpitool/pkg/notification"
	"github.com/moorada/neferpitool/pkg/reliableChanges"
	"github.com/robfig/cron/v3"
)

func background() {
	timeToSleepBackground := time.Duration(configuration.GetConf().MINUTESLEEPBACKGROUNDMONITORING) * time.Minute
	for _, c := range configuration.GetConf().REPORTFREQUENCY {
		runCronJob(c)
	}
	mds := db.GetMainDomainListFromDB()
	if len(mds) > 0 {
		for {
			log.Info("Monotoring...")
			backgroundWork()
			s := spinner.New(spinner.CharSets[26], 200*time.Millisecond) // Build our new spinner
			s.Prefix = "Sleeping "
			s.Start()
			time.Sleep(timeToSleepBackground)
			s.Stop()
		}
	} else {
		log.Error("No domains in the Database")
	}

}

func backgroundWork() {
	start := time.Now()
	checkChangesOfAll()
	if len(changesToSend) > 0 {
		prepareAndSendEmail()
	}
	elapsed := time.Since(start)
	elapsedMin := int(elapsed / time.Minute)
	log.Debug("Time of full scansion: %v", elapsed)
	if elapsedMin != 0 {
		log.Debug("Typodomains scanned per minute: %v", totaltd/elapsedMin)
	}
}

func runCronJob(expression string) {
	c := cron.New()
	_, err := c.AddFunc(expression, func() { prepareAndSendReportEmail(expression) })

	if err != nil {
		log.Error("Cron error with cron %s: %s", expression, err)
	} else {
		log.Debug("Cron Job started with cron %s", expression)
		c.Start()
	}
}

func prepareAndSendReportEmail(expression string) {
	err, reliableChangesReportToSend := db.GetRelaibleChangesFromDBWithExpression(expression)
	if err != nil {
		log.Error(err.Error())
	}

	conf := configuration.GetConf()

	if conf.EMAIL != "" && conf.PASSWORD != "" && len(conf.EMAILTONOTIFY) != 0 {

		tdsInExpiration := getTypoDomainsInExpiration()
		headersStatus, datasStatus, headersWhois, datasWhois := reliableChangesReportToSend.ToTables()
		hExpiry, dExpiry := tdsInExpiration.ToExpiryTable()

		request := notification.Request{
			From:     conf.EMAIL,
			Password: conf.PASSWORD,
			To:       conf.EMAILTONOTIFY,
			Subject:  "domain monitoring - Report of the day",
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

		err := notification.EmailChanges(tpl, request)
		if err != nil {
			log.Error(err.Error())
		} else {
			log.Info("Report sent by email")
		}
		for _, c := range reliableChangesReportToSend {
			i := reliableChanges.Contains(c.Crons, expression)
			if i >= 0 {
				db.DeleteExprToChDB(c, *c.Crons[i])
			}
		}

	} else {
		log.Info("No email to send")
	}
}

func prepareAndSendEmail() {

	conf := configuration.GetConf()

	if conf.EMAIL != "" && conf.PASSWORD != "" && len(conf.EMAILTONOTIFY) != 0 {

		tdsInExpiration := getTypoDomainsInExpiration()
		headersStatus, datasStatus, headersWhois, datasWhois := changesToSend.ToTables()
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

		err := notification.EmailChanges(tpl, request)
		if err != nil {
			log.Error(err.Error())
		} else {
			log.Info("Changes sent by email")
		}
	} else {
		log.Info("No email to send")
	}
	changesToSend = changes.ChangeList{}
}
