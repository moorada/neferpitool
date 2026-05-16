package cmd

import (
	"github.com/cheggaaa/pb/v3"
	"github.com/manifoldco/promptui"
	"github.com/moorada/neferpitool/pkg/changes"
	"github.com/moorada/neferpitool/pkg/console"
	"github.com/moorada/neferpitool/pkg/db"
	"github.com/moorada/neferpitool/pkg/domains"
	"github.com/moorada/neferpitool/pkg/log"
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

	if changes != nil {
		monitorService.SaveReliableChanges(changes)
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
	errs := map[string]error{}
	var bar *pb.ProgressBar
	progress := func(done, total int) {
		if total == 0 {
			return
		}
		if done == 1 {
			bar = pb.Full.Start(total)
		}
		if bar != nil {
			bar.SetCurrent(int64(done))
			if done == total {
				bar.Finish()
				bar = nil
			}
		}
	}

	tdsReliable, changesReliable, errs = monitorService.IterateCheckGetChanges(tds, progress)
	if len(errs) > 0 {
		console.PrintTableErrs(errs)
	}
	return tdsReliable, changesReliable
}

func getTypoDomainsInExpiration() domains.TypoList {
	return monitorService.GetTypoDomainsInExpiration()
}
