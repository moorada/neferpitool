package console

import (
	"os"
	"strconv"

	"golang.org/x/net/idna"

	"github.com/evilsocket/islazy/tui"

	"github.com/moorada/neferpitool/pkg/configuration"
	"github.com/moorada/neferpitool/pkg/log"
	"github.com/moorada/neferpitool/pkg/reliableChanges"

	"github.com/moorada/neferpitool/pkg/changes"
	"github.com/moorada/neferpitool/pkg/domains"
)

func PrintTableErrs(errs map[string]error) {
	if !tui.Effects() {
		log.Info("Tui effects not available on this terminal.\n")
		log.Info("%-20s | %-60s |\n", "TYPODOMAIN", "ERRORS")
		for k, v := range errs {
			log.Info("%-20s | %-60s |\n", k, v.Error())
		}
	} else {
		columns := []string{
			"Typodomain",
			"Errors",
		}

		var rows [][]string
		for k, v := range errs {
			er := v.Error()
			if len(er) > 60 {
				er = er[:60]
			}
			row := []string{tui.Bold(k), er}
			rows = append(rows, row)
		}
		tui.Table(os.Stdout, columns, rows)
	}

	log.Debug("Table logs debug.\n")
	log.Debug("%-20s | %-60s |\n", "TYPODOMAIN", "ERRORS")
	for k, v := range errs {
		log.Debug("%-20s | %-60s |\n", k, v.Error())
	}
}

/*Print a short view of all typo-domains*/
func PrintTableTypoDomains(tdt domains.TypoList) {
	if !tui.Effects() {
		log.Info("Tui effects not available on this terminal.\n")
		log.Info("%-20s | %-15s | %-35s | %-20s | %-25s | %-25s | %-25s\n", "DOMAIN NAME", "AVAILABILITY", "ALGORITHM", "MAIN DOMAIN", "REGISTRAR", "EXPIRY DATE", "UNICODE", "DATE CHECK")
		for _, t := range tdt {
			t.PrintShort()
		}
	} else {
		// ASCII tables on the terminal
		columns := []string{
			"Idna",
			"Status",
			"Algorithm",
			"Main DomainName",
			"Registrar",
			"Expiry date",
			"Date check",
			"Unicode",
		}

		var rows [][]string

		for _, td := range tdt {

			var p *idna.Profile
			// Raw Punycode has no restrictions and does no mappings.
			p = idna.New()
			nameUnicode, err := p.ToUnicode(td.Name)
			if err != nil {
				log.Error(err.Error())
				nameUnicode = td.Name
			}

			w := td.GetWhois()
			timeString := td.GetExpiryDateString()
			row := []string{tui.Bold(td.Name), td.StatusToStringColored(), td.Algorithm, td.LegitDomain, limitLen(w.Parsed.Registrar.RegistrarName, 25), timeString, td.CreatedAt.Format("02/01/2006 15:04"), nameUnicode}

			visibilityToShow := configuration.GetConf().SHOWSTATUS

			for _, v := range visibilityToShow {
				if td.StatusToString() == v {
					rows = append(rows, row)
				}
			}
		}
		tui.Table(os.Stdout, columns, rows)
	}

	log.Debug("Table logs debug.\n")
	log.Debug("%-20s |%-15s | %-15s | %-20s | %-25s | %-20s | %-25s | %-20s\n", "DOMAIN NAME", "AVAILABILITY", "ALGORITHM", "MAIN DOMAIN", "REGISTRAR", "EXPIRY DATE", "UNICODE", "DATE CHECK")
	for _, t := range tdt {
		t.PrintShortDebug()
	}

}

func limitLen(s string, i int) string {
	var result string
	if len(s) > i {
		result = s[:i-3]
		result += "..."
	}
	return result
}

func PrintChanges(changes changes.ChangeList) {

	if !tui.Effects() {
		log.Info("Tui effects not available on this terminal.\n")
		log.Info("%-20s |%-20s | %-35s | %-35s |\n", "TYPODOMAIN", "FIELD CHANGED", "BEFORE", "AFTER")
		for _, c := range changes {
			log.Info("%-20s | %-35s | %-35s |%-35s |\n", c.TypoDomain, c.Field, c.Before, c.After)
		}
	} else {
		columns := []string{
			"Typodomain",
			"Field changed",
			"Before",
			"After",
		}

		var rows [][]string
		for _, c := range changes {
			row := []string{tui.Bold(c.TypoDomain), c.Field, c.Before, c.After}
			rows = append(rows, row)
		}
		tui.Table(os.Stdout, columns, rows)
	}

	log.Debug("Table logs debug.\n")
	log.Debug("%-20s |%-20s | %-35s | %-35s |\n", "TYPODOMAIN", "FIELD CHANGED", "BEFORE", "AFTER")
	for _, c := range changes {
		log.Debug("%-20s | %-35s | %-35s |%-35s |\n", c.TypoDomain, c.Field, c.Before, c.After)
	}
}

func PrintReliableChanges(changes reliableChanges.ReliableChangeList) {

	if !tui.Effects() {
		log.Info("Tui effects not available on this terminal.\n")
		log.Info("%-20s |%-20s | %-35s | %-35s |\n", "TYPODOMAIN", "FIELD CHANGED", "BEFORE", "AFTER")
		for _, c := range changes {
			log.Info("%-20s | %-35s | %-35s |%-35s |\n", c.TypoDomain, c.Field, c.Before, c.After)
		}
	} else {
		columns := []string{
			"Typodomain",
			"Field changed",
			"Before",
			"After",
		}

		var rows [][]string
		for _, c := range changes {
			row := []string{tui.Bold(c.TypoDomain), c.Field, c.Before, c.After}
			rows = append(rows, row)
		}
		tui.Table(os.Stdout, columns, rows)
	}

	log.Debug("Table logs debug.\n")
	log.Debug("%-20s |%-20s | %-35s | %-35s |\n", "TYPODOMAIN", "FIELD CHANGED", "BEFORE", "AFTER")
	for _, c := range changes {
		log.Debug("%-20s | %-35s | %-35s |%-35s |\n", c.TypoDomain, c.Field, c.Before, c.After)
	}
}

func PrintAverageWhoisStats(averages map[string]int) {

	log.Info("Average times whois")

	if !tui.Effects() {
		log.Info("Tui effects not available on this terminal.\n")
		log.Info("%-20s |%-20s | %-35s | %-35s |\n", "TLD", "AVERAGE TIME")
		for tld, aTime := range averages {
			log.Info("%-35s | %-35s \n", tld, aTime)
		}
	} else {
		columns := []string{
			"TLD",
			"AVERAGE TIME in milliseconds",
		}

		var rows [][]string
		for tld, aTime := range averages {
			row := []string{tui.Bold(tld), strconv.Itoa(aTime)}
			rows = append(rows, row)
		}
		tui.Table(os.Stdout, columns, rows)
	}

	log.Debug("Table logs debug.\n")
	log.Debug("%-20s |%-20s | %-35s | %-35s |\n", "TLD", "AVERAGE TIME")
	for tld, aTime := range averages {
		log.Debug("%-35s | %-35s \n", tld, aTime)
	}
}
