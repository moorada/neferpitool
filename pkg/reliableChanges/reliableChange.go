package reliableChanges

import (
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/moorada/neferpitool/pkg/log"
	"golang.org/x/net/idna"
)

type ReliableChange struct {
	gorm.Model
	TypoDomain string
	Field      string
	Before     string
	After      string
}

type ReliableChangeList []ReliableChange

func (tdcs ReliableChangeList) ToTables() (headersAvailability []string, datasStatus [][]string, headersWhois []string, datasWhois [][]string) {

	headersAvailability = []string{
		"TypoDomain",
		"Before",
		"After",
	}

	var tdcsWhois ReliableChangeList

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

func (tdcs ReliableChangeList) changesWhoisToTable() ([]string, [][]string) {

	s := reflect.ValueOf(ReliableChange{})
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
	var dataWithoutGormModel [][]string
	for _, r := range data {
		dataWithoutGormModel = append(dataWithoutGormModel, r[1:])
	}

	return headers, dataWithoutGormModel
}
