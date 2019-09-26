package cmd

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/moorada/neferpitool/pkg/changes"
	"github.com/moorada/neferpitool/pkg/configuration"
	"github.com/moorada/neferpitool/pkg/console"
	"github.com/moorada/neferpitool/pkg/constants"
	"github.com/moorada/neferpitool/pkg/db"
	"github.com/moorada/neferpitool/pkg/domains"
	"github.com/moorada/neferpitool/pkg/log"
	"github.com/moorada/neferpitool/pkg/notification"
	"github.com/moorada/neferpitool/pkg/stats"
)

var changesToSend changes.ChangeList

func MonitorCmd() {

	for true {
		prompt := promptui.Select{
			Label: "Monitor mode",
			Items: []string{"Domains", "Show all Typodomains in Expiration", "Check for changes", "Active changes in background", "Stats"},
		}

		_, result, err := prompt.Run()

		if err != nil {
			log.Error("Prompt failed %v\n", err)
			return
		}
		switch result {
		case "Domains":
			manageDomains()
		case "Show all Typodomains in Expiration":
			showTypoDomainsInExpiration()
		case "Check for changes":
			checkForChanges()
		case "Active changes in background":
			background()
		case "Stats":
			stats.PrintWhoisStats()
		}
	}
}

func showTypoDomainsInExpiration() {
	tdsEx := getTypoDomainsInExpiration()
	if len(tdsEx) != 0 {
		console.PrintTableTypoDomains(tdsEx)
	} else {
		log.Error("No expiry typodomains in the Database")
	}

}
func manageDomains() {
	mds := db.GetMainDomainListFromDB()
	var mdsnames []string
	for _, d := range mds {
		mdsnames = append(mdsnames, d.Name)
	}

	prompt := promptui.Select{
		Label: "Domains",
		Items: mdsnames,
	}

	_, result, err := prompt.Run()

	if err != nil {
		log.Error("Prompt failed %v\n", err)
		return
	}

	log.Info("Manage %s", result)
	SingleDomainMode(result)

}

func checkForChanges() {

	mds := db.GetMainDomainListFromDB()
	var mdsnames []string

	mdsnames = append(mdsnames, "All")

	for _, d := range mds {
		mdsnames = append(mdsnames, d.Name)
	}

	prompt := promptui.Select{
		Label: "Check for changes",
		Items: mdsnames,
	}

	_, result, err := prompt.Run()

	if err != nil {
		log.Error("Prompt failed %v\n", err)
		return
	}
	if result == "All" {
		checkChangesOfAll()
	} else {
		log.Info("Checking changes about %s", result)
		checkChanges(db.GetTypoDomainListFromDB(result))
	}
}

func checkChangesOfAll() {
	mds := db.GetMainDomainListFromDB()
	for i, d := range mds {
		log.Info("Checking changes about %s, %v di %v", d.Name, i+1, len(mds))
		checkChanges(db.GetTypoDomainListFromDB(d.Name))
	}
}

func background() {
	timeToSleepBackground := time.Duration(configuration.GetConf().HOURSLEEPBACKGROUNDMONITORING) * time.Hour
	mds := db.GetMainDomainListFromDB()
	if len(mds) != 0 {
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
	log.RemoveDebugLog()
	log.ActiveDebugLog()
	start := time.Now()
	checkChangesOfAll()
	prepareAndSendEmail()
	elapsed := time.Since(start)
	elapsedMin := int(elapsed / time.Minute)
	log.Debug("Time of full scansion: %v", elapsed)
	if elapsedMin != 0 {
		log.Debug("Typodomains scanned per minute: %v", totaltd/elapsedMin)
	}
}

func checkChanges(tds domains.TypoList) bool {
	tdsChanged, changes := iterateCheckGetChanges(tds)

	if changes != nil {
		db.AddTypoListToDB(tdsChanged)
		console.PrintChanges(changes)
		changesToSend = append(changesToSend, changes...)
		return true
	} else {
		log.Info("%s", "no changes")
	}
	return false
}

func iterateCheckGetChanges(tds domains.TypoList) (tdsReliable []domains.TypoDomain, changesReliable []changes.Change) {

	tdsNew := tds.GetUnfilledCopy()
	errs := UpdateTypoDomainsWithProgressBar(tdsNew)
	if len(errs) != 0 {
		console.PrintTableErrs(errs)
	}
	console.PrintTableErrs(errs)
	tdsOldCh, tdsNewCh, chs := changes.MakeChangeList(tds, tdsNew)
	for i := 0; i < 2 && len(tdsOldCh) > 0; i++ {
		log.Info("Checking reliability about %v changes", len(chs))
		s := spinner.New(spinner.CharSets[26], 200*time.Millisecond) // Build our new spinner
		s.Prefix = "Sleeping "
		s.Start()
		time.Sleep(1 * time.Minute)
		s.Stop()
		errs := UpdateTypoDomainsWithProgressBar(tdsNewCh)
		if len(errs) != 0 {
			console.PrintTableErrs(errs)
		}
		tdsOldChNext, tdsNewChNext, chsNext := changes.MakeChangeList(tdsOldCh, tdsNewCh)
		tdsReliable, changesReliable = chsNext.FilterReliableWithPrev(chs, tdsNewCh, tdsNewChNext)
		tdsOldCh, tdsNewCh, chs = tdsOldChNext, tdsReliable, changesReliable
	}

	return tdsReliable, changesReliable
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
		}
	} else {
		log.Info("No email to send")
	}
	changesToSend = []changes.Change{}
}

func getTypoDomainsInExpiration() domains.TypoList {

	var totalTdInExpiration domains.TypoList

	ds := db.GetMainDomainListFromDB()

	for _, d := range ds {
		exp := configuration.GetConf().EXPIRATIONTIME
		tdt := db.GetTypoDomainListWithStatusFromDB(d.Name, []int{constants.INACTIVE, constants.ACTIVE, constants.ALIAS})
		tds := tdt.FilterInExpiration(exp)
		totalTdInExpiration = append(totalTdInExpiration, tds...)
	}

	return totalTdInExpiration

}
