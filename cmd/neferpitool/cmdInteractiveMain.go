package cmd

import (
	"github.com/manifoldco/promptui"
	"github.com/moorada/neferpitool/pkg/changes"
	"github.com/moorada/neferpitool/pkg/db"
	"github.com/moorada/neferpitool/pkg/log"
	"github.com/moorada/neferpitool/pkg/stats"
)

var changesToSend changes.ChangeList

func MonitorCmd() {

	for true {
		prompt := promptui.Select{
			Label: "Monitor mode",
			Items: []string{"Domains", "Show all Typodomains in Expiration", "Check for changes", "Active changes in background", "Stats", "Test"},
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
		case "Test":
			test()
		}
	}
}

func test() {
	log.Info("Empty")

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
