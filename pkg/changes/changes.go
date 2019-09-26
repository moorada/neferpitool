package changes

import (
	"reflect"

	"golang.org/x/net/idna"

	"github.com/moorada/neferpitool/pkg/domains"
	"github.com/moorada/neferpitool/pkg/log"
)

type Change struct {
	TypoDomain string
	Field      string
	Before     string
	After      string
}

type ChangeList []Change

const (
	NAME_REGISTRANT = "RegistrantName Registrant"
	ORGANIZAIOTN    = "Organization"
	CREATION_DATE   = "Creation Date"
	UPDATED_DATE    = "Updated Date"
	EXPIRATION_DATE = "Expiration Date"
	STATUS          = "Status"
)

/*Compare two changes and return true is their are equals*/
func (c Change) IsEqual(c2 Change) bool {
	return reflect.DeepEqual(c, c2)

}

func MakeChangeList(tdsOld domains.TypoList, tdsNew domains.TypoList) (tdsOldCh domains.TypoList, tdsNewCh domains.TypoList, tdcs ChangeList) {

	for i, _ := range tdsOld {
		tdOld := tdsOld[i]
		tdNew := tdsNew[i]

		reliability := tdNew.IsReliableAboutPrev(tdOld)
		if !reliability {
			continue
		}

		w1 := tdOld.GetWhois()
		w2 := tdNew.GetWhois()

		name := w1.Parsed.Registrant.RegistrantName == w2.Parsed.Registrant.RegistrantName
		organization := w1.Parsed.Registrant.Organization == w2.Parsed.Registrant.Organization
		expirationDate := w1.Parsed.Registrar.ExpirationDate == w2.Parsed.Registrar.ExpirationDate

		status := tdOld.Status == tdNew.Status
		if !status {
			tdcs = append(tdcs, Change{tdOld.Name, STATUS, tdOld.StatusToString(), tdNew.StatusToString()})
		} else {
			if !name {
				tdcs = append(tdcs, Change{tdOld.Name, NAME_REGISTRANT, w1.Parsed.Registrant.RegistrantName, w2.Parsed.Registrant.RegistrantName})
			}
			if !organization {
				tdcs = append(tdcs, Change{tdOld.Name, ORGANIZAIOTN, w1.Parsed.Registrant.Organization, w2.Parsed.Registrant.Organization})
			}
			if !expirationDate {
				tdcs = append(tdcs, Change{tdOld.Name, EXPIRATION_DATE, w1.Parsed.Registrar.ExpirationDate, w2.Parsed.Registrar.ExpirationDate})
			}
		}

		if !name || !organization || !expirationDate || !status {
			if !status {
				log.Debug("%s is changed about status", tdNew.Name)
			} else {
				log.Debug("%s is changed about whois", tdNew.Name)
			}
			tdsOldCh = append(tdsOldCh, tdOld)
			tdsNewCh = append(tdsNewCh, tdNew)
		}
	}
	return
}

func (tdcs ChangeList) FilterReliableWithPrev(tdcsPrev ChangeList, tdsPrev domains.TypoList, tds domains.TypoList) (tdsChecked domains.TypoList, tdcsChecked ChangeList) {
	tdsMap := tds.ToMap()
	tdsPrevMap := tdsPrev.ToMap()

	for key, val := range tdsMap {
		if tdsPrevMap[key].IsEqual(val) {
			tdsChecked = append(tdsChecked, val)
		}
	}

	tdcsMap := tdcs.ToMap()
	tdcsPrevMap := tdcsPrev.ToMap()

	for key, val := range tdcsMap {
		if reflect.DeepEqual(val, tdcsPrevMap[key]) {
			tdcsChecked = append(tdcsChecked, val)
		}
	}
	return
}

func (tdcs ChangeList) ToMap() map[string]Change {
	tdcsMap := make(map[string]Change)
	for _, c := range tdcs {
		tdcsMap[c.TypoDomain+"|"+c.Field] = c
	}
	return tdcsMap
}

func (tdcs ChangeList) ToTables() (headersAvailability []string, datasStatus [][]string, headersWhois []string, datasWhois [][]string) {

	headersAvailability = []string{
		"TypoDomain",
		"Before",
		"After",
	}

	var tdcsWhois ChangeList

	for _, ch := range tdcs {
		var p *idna.Profile
		// Raw Punycode has no restrictions and does no mappings.
		p = idna.New()
		nameUnicode, err := p.ToUnicode(ch.TypoDomain)
		if err != nil {
			log.Error("error to convert in unicode %s:", ch.TypoDomain, err.Error())
			nameUnicode = ch.TypoDomain
		}

		if ch.Field == "Status" {
			datasStatus = append(datasStatus, []string{nameUnicode, ch.Before, ch.After})
		} else {
			tdcsWhois = append(tdcsWhois, ch)
		}
	}
	hw, dw := tdcsWhois.changesWhoisToTable()
	headersWhois, datasWhois = hw, dw
	return
}

func (tdcs ChangeList) changesWhoisToTable() ([]string, [][]string) {

	s := reflect.ValueOf(Change{})
	typeOfT := s.Type()
	var headers []string
	for i := 0; i < s.NumField(); i++ {
		headers = append(headers, typeOfT.Field(i).Name)
	}

	var data [][]string
	for _, d := range tdcs {
		v := reflect.ValueOf(d)

		raw := make([]string, v.NumField())

		for i := 0; i < v.NumField(); i++ {
			raw[i] = v.Field(i).String()
		}

		data = append(data, raw)
	}
	return headers, data
}
