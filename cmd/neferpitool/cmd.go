package cmd

//cmdroot
import (
	"bufio"
	"flag"
	"os"
	"time"

	"github.com/briandowns/spinner"

	"github.com/cheggaaa/pb/v3"

	"github.com/moorada/neferpitool/pkg/configuration"
	"github.com/moorada/neferpitool/pkg/scanner"

	"github.com/moorada/neferpitool/pkg/console"
	"github.com/moorada/neferpitool/pkg/db"
	"github.com/moorada/neferpitool/pkg/domains"
	"github.com/moorada/neferpitool/pkg/generator"
	"github.com/moorada/neferpitool/pkg/log"
)

var totaltd int

const PathConfigFolder = "./config"

func CmdRoot() {

	//init log
	if err := log.ActiveConsoleLog(); err != nil {
		panic(err)
	} else {
		defer log.Close()
	}

	err := os.MkdirAll(PathConfigFolder, os.ModePerm)
	if err != nil {
		log.Error(err.Error())
	}

	//init db
	db.InitDB("config/database")
	defer db.CloseDB()

	// init flags
	logs := flag.Bool("logs", false, "Avtive logs on file")
	singleTd := flag.String("td", "", "Manage one typodomain")
	bg := flag.Bool("bg", false, "Active monitoring in background")
	pd := flag.Bool("pd", false, "Check if domains are present")
	makeConfig := flag.Bool("mc", false, "Make config file")
	importTds := flag.String("it", "", "Import Typos from file - main domain")
	pathImportTds := flag.String("p", "", "Import Typos from file - path of the file")

	flag.Parse()

	if *logs {
		if err := log.ActiveDebugLog(); err != nil {
			panic(err)
		}
	}

	if *makeConfig {
		configuration.MakeConfigFile()
		return
	}

	logDegubInfo()

	if *bg {
		background()
		return
	}

	if *singleTd != "" {
		SingleTdMode(*singleTd)
		return
	}

	if *importTds != "" {
		if *pathImportTds != "" {
			importTypos(*importTds, *pathImportTds)
			return
		} else {
			log.Error("%s", "Specify the path of the import file: -p /path/example/ ")
		}

	}

	args := flag.Args()

	if len(args) < 1 {
		MonitorCmd()
		return
	}

	if *pd {
		presenceDomains(args)
		return
	}

	if len(args) == 1 {
		domain := args[0]
		SingleDomainMode(domain)
	} else {
		multipleDomainsMode(args)
	}

	os.Exit(0)
}

func presenceDomains(domains []string) {
	ds := db.GetMainDomainListFromDB()
	dsMap := map[string]bool{}

	for _, d := range ds {
		dsMap[d.Name] = true
	}

	for _, s := range domains {

		if _, ok := dsMap[s]; ok {
			log.Info("%s present", s)
		} else {
			log.Info("%s NOT present", s)
		}

	}
}

func multipleDomainsMode(domains []string) {
	for _, d := range domains {
		addDomainAndHisTypos(d)
		console.PrintTableTypoDomains(db.GetTypoDomainListFromDB(d))
	}
}

func logDegubInfo() {
	// debug info...
	totaltd = 0
	totald := db.GetMainDomainListFromDB()

	domains := ""

	for _, d := range totald {
		td1d := db.GetTypoDomainListFromDB(d.Name)
		domains += d.Name + ", "
		totaltd = totaltd + len(td1d)

	}
	log.Debug("%v domains in db: %s", len(totald), domains)
	log.Debug("Number of typodomains in db: %v", totaltd)
	//
}

func UpdateTypoDomainsWithProgressBar(tds domains.TypoList) map[string]error {

	c := make(chan int, len(tds))
	bar := pb.Full.Start(len(tds))

	var errs map[string]error
	go func(map[string]error) {
		errs = scanner.UpdateTypoDomains(tds, c)
	}(errs)

	for range tds {
		<-c
		bar.Increment()
	}
	bar.Finish()
	if len(tds) > 0 {
		log.Debug("Stats Scansion, Typodomains: %v, Errors: %v, Percentage of errors: %v", len(tds), len(errs), len(errs)*100/len(tds))
	} else {
		log.Debug("no typodomains")
	}
	return errs
}

func addDomainAndHisTypos(domain string) {

	s := spinner.New(spinner.CharSets[26], 200*time.Millisecond)
	s.Prefix = "Generating typodomains "
	s.Start()
	tds := generator.GetUnfilledTypoDomains(domain)
	s.Stop()
	errs := UpdateTypoDomainsWithProgressBar(tds)
	if len(errs) != 0 {
		console.PrintTableErrs(errs)
	}

	db.AddTypoListToDB(tds)
	md := domains.NewLegitDomain(domain)

	err := md.Update()
	if err != nil {
		log.Error("%s, error: %s", domain, err.Error())
	}

	db.AddLegitDomainToDB(md)

	log.Info("Typodomains added to database")

}

func importTypos(domain string, path string) {

	s := spinner.New(spinner.CharSets[26], 200*time.Millisecond)
	s.Prefix = "Importing typodomains "
	s.Start()

	var tds domains.TypoList

	file, err := os.Open(path)
	if err != nil {
		log.Error("Importing typos from file, error: %s", domain, err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var typos []string
	for scanner.Scan() {
		typos = append(typos, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Error("Importing typos from file, error: %s", domain, err.Error())
	}

	if len(typos) > 0 {
		for _, t := range typos {
			if t != "" {
				td := domains.NewTypoDomain(t, domain, "imported")
				tds = append(tds, td)
			}
		}
	} else {
		log.Error("Empty file, error: %s", domain, err.Error())
	}

	s.Stop()

	errs := UpdateTypoDomainsWithProgressBar(tds)
	if len(errs) != 0 {
		console.PrintTableErrs(errs)
	}

	db.AddTypoListToDB(tds)

	log.Info("Typodomains imported")

}
