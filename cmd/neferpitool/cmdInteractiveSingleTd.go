package cmd

import (
	"github.com/manifoldco/promptui"

	"github.com/moorada/neferpitool/pkg/console"
	"github.com/moorada/neferpitool/pkg/constants"
	"github.com/moorada/neferpitool/pkg/db"
	"github.com/moorada/neferpitool/pkg/domains"
	"github.com/moorada/neferpitool/pkg/log"
)

func SingleTdMode(typodomain string) {
	tds := db.GetTypoDomainHistoryFromDB(typodomain)

	if len(tds) > 0 {

		prompt := promptui.Select{
			Label: "TypoDomain in the Database",
			Items: []string{"Look at his history", "Check changes"},
		}

		_, result, err := prompt.Run()

		if err != nil {
			log.Fatal("Prompt failed %v\n", err)
			return
		}

		switch result {
		case "Look at his history":
			console.PrintTableTypoDomains(db.GetTypoDomainHistoryFromDB(typodomain))
		case "Check changes":
			checkChanges([]domains.TypoDomain{db.GetTypoDomainFromDB(typodomain)})
		}
		log.Info("You choose %q\n", result)
	} else {
		log.Info("This typodomain isn't in the database")
		updateTypo(typodomain)
	}
}

func updateTypo(typoDomain string) {

	td := db.GetTypoDomainFromDB(typoDomain)
	tdNew := domains.NewTypoDomain(td.Name, td.LegitDomain, td.Algorithm)
	err := tdNew.UpdateStatus()
	if err != nil {
		log.Error("%s %s", typoDomain, err.Error())
	}
	if tdNew.Status == constants.INACTIVE || tdNew.Status == constants.ACTIVE {
		err = tdNew.UpdateWhois()
		if err != nil {
			log.Error("%s %s", typoDomain, err.Error())
		}
	}

	db.AddTypoDomainToDB(tdNew)

	console.PrintTableTypoDomains(db.GetTypoDomainHistoryFromDB(typoDomain))

}
