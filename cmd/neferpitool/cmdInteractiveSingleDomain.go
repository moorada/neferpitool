package cmd

import (
	"github.com/manifoldco/promptui"

	"github.com/moorada/neferpitool/pkg/console"
	"github.com/moorada/neferpitool/pkg/db"
	"github.com/moorada/neferpitool/pkg/log"
)

func SingleDomainMode(domain string) {

	tds := db.GetTypoDomainListFromDB(domain)

	if len(tds) > 0 {

		prompt := promptui.Select{
			Label: "DomainName in the Database",
			Items: []string{"Look at his typo-domains", "Check for changes", "Remove"},
		}

		_, result, err := prompt.Run()

		if err != nil {
			log.Fatal("Prompt failed %v\n", err)
			return
		}

		switch result {
		case "Look at his typo-domains":
			console.PrintTableTypoDomains(db.GetTypoDomainListFromDB(domain))
		case "Check for changes":
			checkChanges(db.GetTypoDomainListFromDB(domain))
		case "Remove":
			db.RemoveLegitDomainFromDB(domain)
		}
		log.Info("You choose %q\n", result)
	} else {

		log.Info("This domain isn't in the database")
		addDomainAndHisTypos(domain)
		console.PrintTableTypoDomains(db.GetTypoDomainListFromDB(domain))
	}
}
