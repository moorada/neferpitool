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
	"github.com/moorada/neferpitool/pkg/reliableChanges"
)

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

func showTypoDomainsInExpiration() {
	tdsEx := getTypoDomainsInExpiration()
	if len(tdsEx) != 0 {
		console.PrintTableTypoDomains(tdsEx)
	} else {
		log.Error("No expiry typodomains in the Database")
	}

}

func checkChangesOfAll() {
	mds := db.GetMainDomainListFromDB()
	for i, d := range mds {
		log.Info("Checking changes about %s, %v di %v", d.Name, i+1, len(mds))
		checkChanges(db.GetTypoDomainListFromDB(d.Name))
	}
}

func checkChanges(tds domains.TypoList) bool {

	tdsChanged, changes := iterateCheckGetChanges(tds)

	var Crons = []*reliableChanges.CronExpression{}
	for _, c := range configuration.GetConf().REPORTFREQUENCY {

		cc := db.GetCronExpressionFromDB(c)
		if cc.ID != 0 {
			Crons = append(Crons, &cc)
		} else {
			Crons = append(Crons, &reliableChanges.CronExpression{Exrpression: c})
		}
	}

	if changes != nil {

		var rChanges []reliableChanges.ReliableChange

		for _, c := range changes {
			rChanges = append(rChanges, reliableChanges.ReliableChange{TypoDomain: c.TypoDomain, Field: c.Field, Before: c.Before, After: c.After, Crons: Crons})
		}
		db.AddReliableChangeListToDB(rChanges)

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
	if len(errs) > 0 {
		console.PrintTableErrs(errs)
	}
	tdsOldCh, tdsNewCh, chs := changes.MakeChangeList(tds, tdsNew)
	for i := 0; i < 2 && len(tdsOldCh) > 0; i++ {
		log.Info("Checking reliability about %v changes", len(chs))
		s := spinner.New(spinner.CharSets[26], 200*time.Millisecond) // Build our new spinner
		s.Prefix = "Sleeping "
		s.Start()
		time.Sleep(time.Duration(configuration.GetConf().CHECKRELIABILITYTIME) * time.Millisecond)
		s.Stop()
		errs := UpdateTypoDomainsWithProgressBar(tdsNewCh)
		if len(errs) > 0 {
			console.PrintTableErrs(errs)
		}
		tdsOldChNext, tdsNewChNext, chsNext := changes.MakeChangeList(tdsOldCh, tdsNewCh)
		tdsReliable, changesReliable = chsNext.FilterReliableWithPrev(chs, tdsNewCh, tdsNewChNext)
		tdsOldCh, tdsNewCh, chs = tdsOldChNext, tdsReliable, changesReliable
	}

	return tdsReliable, changesReliable
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
