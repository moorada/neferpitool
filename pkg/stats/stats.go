package stats

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/moorada/neferpitool/pkg/console"
	"github.com/moorada/neferpitool/pkg/db"
	"github.com/moorada/neferpitool/pkg/log"
	"github.com/moorada/neferpitool/pkg/whois"
)

func PrintWhoisStats() {
	//statTimeWhois := make(map[string]int)

	tds := db.GetAllTypoDomainListFromDB()
	fmt.Println(len(tds))
	i := 0
	for _, td := range tds {
		if td.ErrorWhois != "" || td.ErrorStatus != "" {
			i++
		}
	}
	fmt.Println("number of typodomains", len(tds))
	fmt.Println("number of typodomains with error", i)

	if len(tds) != 0 {
		var data = map[string]map[string]int{}
		for _, td := range tds {
			timeTD, err := time.ParseDuration(td.TimeWhois)
			if !td.Whois.IsLike(whois.Whois{}) && td.ErrorWhois == "" {
				if err != nil {
					log.Error("Err:%s, time:%s, td:%s", err.Error(), td.TimeWhois, td.Name)
				} else {
					tMS := int(timeTD / time.Millisecond)
					if err != nil {
						log.Error(err.Error())
					} else {
						tld := getTLD(td.Name)
						if data[tld] == nil {
							data[tld] = map[string]int{
								td.Name: tMS,
							}
						} else {
							data[tld][td.Name] = tMS
						}
					}
				}
			}
		}

		var averageList []int

		averageTable := make(map[string]int)

		for tld, data := range data {
			var timeTLDTotal int
			for _, time := range data {
				timeTLDTotal += time
			}
			average := timeTLDTotal / len(data)
			averageList = append(averageList, average)
			averageTable[tld] = average
			log.Debug("Average time whois about %s: %v milliseconds", tld, average)
		}

		averageTableSorted := make(map[string]int)

		keys := make([]string, 0, len(averageTable))
		for k := range averageTable {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			averageTableSorted[k] = averageTable[k]
		}

		var timeTotalAverage int
		for _, avg := range averageList {
			timeTotalAverage += avg
		}

		averageTotal := timeTotalAverage / len(averageList)
		console.PrintAverageWhoisStats(averageTableSorted)
		log.Info("Average time TOTAL whois about: %v milliseconds", averageTotal)

	} else {
		log.Error("No domains in the Database")
	}

}

func getTLD(domain string) string {
	domains := strings.Split(domain, ".")
	if len(domains) < 2 {
		log.Fatal("DomainName %s is not valid", domain)
	}
	return domains[len(domains)-1]

}
