package domains

import (
	"sort"
	"time"

	"github.com/moorada/neferpitool/pkg/format"
	"github.com/moorada/neferpitool/pkg/log"
)

type TypoList []TypoDomain

type LegitList []LegitDomain

func (a TypoList) Len() int {
	return len(a)
}
func (a TypoList) Less(i, j int) bool {

	date1, err := a[i].GetExpiryDate()
	if err != nil {
		log.Error(err.Error())
	}
	date2, err := a[j].GetExpiryDate()
	if err != nil {
		log.Error(err.Error())
	}
	return date1.Before(date2)
}
func (a TypoList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (tdl TypoList) ToExpiryTable() ([]string, [][]string) {

	sort.Sort(tdl)

	var headers []string
	headers = append(headers, "Typo Domain")
	headers = append(headers, "Expiry Date")

	var data [][]string
	for _, d := range tdl {
		//w, _ := d.GetWhois()
		w := d.GetWhois()

		timeString := ""

		expiration, err := format.StringToTime(w.Parsed.Registrar.ExpirationDate)
		if err != nil {
			log.Error("Error to parse data of %s, Data: %s, Error: %s", d.Name, w.Parsed.Registrar.ExpirationDate, err.Error())
		}
		timeString = format.TimeToStringConsole(expiration)

		data = append(data, []string{d.Name, timeString})
	}
	return headers, data
}

func (tdl TypoList) GetUnfilledCopy() TypoList {
	var tdlnew TypoList
	for _, td := range tdl {
		tdlnew = append(tdlnew, NewTypoDomain(td.Name, td.LegitDomain, td.Algorithm))
	}
	return tdlnew
}

func (tdl TypoList) FilterInExpiration(d int) TypoList {

	var tdtInExpiration []TypoDomain

	now := time.Now()

	for _, td := range tdl {

		//if td.Whois != "" {
		/*w, err := td.GetWhois()
		if err != nil {
			log.Error(err.Error())
		}
		*/
		w := td.Whois
		if w.Parsed.Registrar.ExpirationDate != "" {
			expiryDate, err := format.StringToTime(w.Parsed.Registrar.ExpirationDate)
			if err != nil {
				log.Error("%s %s %s", td.Name, w.Parsed.Registrar.ExpirationDate, err)
			} else {
				diff := expiryDate.Sub(now)

				days := int(diff.Hours() / 24)
				if days < d && days > 0 {
					tdtInExpiration = append(tdtInExpiration, td)
				}
			}
		}

		//}

	}

	return tdtInExpiration
}

func (tds TypoList) ToMap() map[string]TypoDomain {
	tdsMap := make(map[string]TypoDomain)
	for _, d := range tds {
		tdsMap[d.Name] = d
	}
	return tdsMap
}

func (ds LegitList) ToMap() map[string]LegitDomain {
	dsMap := make(map[string]LegitDomain)
	for _, d := range ds {
		dsMap[d.Name] = d
	}
	return dsMap
}
