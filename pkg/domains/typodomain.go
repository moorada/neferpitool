package domains

import (
	"golang.org/x/net/idna"

	"github.com/moorada/neferpitool/pkg/constants"
	"github.com/moorada/neferpitool/pkg/format"
	"github.com/moorada/neferpitool/pkg/log"
	"github.com/moorada/neferpitool/pkg/whois"
)

type TypoDomain struct {
	Domain
	Algorithm   string
	LegitDomain string
}

/*Make a new typo-DomainName*/
func NewTypoDomain(nameTypoDomain string, mainDomain string, algorithm string) TypoDomain {
	return TypoDomain{Domain: Domain{Name: nameTypoDomain, Status: constants.UNKNOWN, Ignore: false}, LegitDomain: mainDomain, Algorithm: algorithm}
}

/*Check if the typodomain is changed and return new typodomain updated*/
func (td TypoDomain) IsChanged() (bool, TypoDomain) {
	tdNew := td
	err := tdNew.Update()
	if err != nil {
		log.Error("IsChanged TypoDomain:%s, Error:%s", td.Name, err.Error())
	}
	if td.IsEqual(tdNew) {
		return false, tdNew
	} else {
		return true, tdNew
	}
}

/*Compare two main domains and return true is they are equals*/
func (td TypoDomain) IsEqual(td2 TypoDomain) bool {
	return td.LegitDomain == td2.LegitDomain && td.Domain.IsEqual(td2.Domain)
}

/*Print a short view of a typo-DomainName*/
func (td TypoDomain) PrintShort() {
	//w, err := td.GetWhois()
	/*if err != nil {
		log.Error(err.Error())
	} else {*/
	w := td.GetWhois()
	timeString := ""
	if w.Parsed.Registrar.ExpirationDate != "" {
		expiration, err := format.StringToTime(w.Parsed.Registrar.ExpirationDate)
		if err != nil {
			log.Error(err.Error())
		}
		timeString = format.TimeToStringConsole(expiration)
	}
	var p *idna.Profile
	// Raw Punycode has no restrictions and does no mappings.
	p = idna.New()
	nameUnicode, err := p.ToUnicode(td.Name)
	if err != nil {
		log.Error(err.Error())
		nameUnicode = td.Name
	}
	log.Info("%-20s |%-15s | %-15s | %-20s | %-25s | %-20s | %-25s | %-20s", td.Name, td.StatusToString(), td.Algorithm, td.LegitDomain, w.Parsed.Registrar.RegistrarName, timeString, nameUnicode, td.CreatedAt.Format("02/01/2006 15:04"))
	//}

}

/*Print a short view of a typo-DomainName*/
func (td TypoDomain) PrintShortDebug() {
	/*w, err := td.GetWhois()
	if err != nil {
		log.Error(err.Error())
	} else {*/
	w := td.GetWhois()
	timeString := ""
	if w.Parsed.Registrar.ExpirationDate != "" {
		expiration, err := format.StringToTime(w.Parsed.Registrar.ExpirationDate)
		if err != nil {
			log.Error(err.Error())
		}
		timeString = format.TimeToStringConsole(expiration)
	}
	var p *idna.Profile
	// Raw Punycode has no restrictions and does no mappings.
	p = idna.New()
	nameUnicode, err := p.ToUnicode(td.Name)
	if err != nil {
		log.Error(err.Error())
		nameUnicode = td.Name
	}
	log.Debug("%-20s |%-15s | %-15s | %-20s | %-25s | %-20s | %-25s | %-20s", td.Name, td.StatusToString(), td.Algorithm, td.LegitDomain, w.Parsed.Registrar.RegistrarName, timeString, nameUnicode, td.CreatedAt.Format("02/01/2006 15:04"))
	//}

}

func (tdNew TypoDomain) IsReliableAboutPrev(tdOld TypoDomain) (reliability bool) {
	reliability = false

	if tdNew.ErrorStatus != "" {
		log.Debug("Insufficient Reliability about Domain %s, tdNew contains error in Availability", tdNew.Name)
		return
	}

	if tdNew.Status == tdOld.Status {
		if tdNew.ErrorWhois != "" {
			log.Debug("Insufficient Reliability about Domain %s, tdNew contains error in Whois", tdNew.Name)
			return
		}

		wOld := tdOld.GetWhois()
		wNew := tdNew.GetWhois()

		wEmpty := whois.Whois{}

		if wOld.Parsed != wEmpty.Parsed && wNew.Parsed == wEmpty.Parsed {
			log.Debug("Insufficient Reliability about Domain %s, loss of information in whois parsed of tdNew", tdNew.Name)
			return
		}
	}
	reliability = true
	return
}
