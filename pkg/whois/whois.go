package whois

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	whoisparser "github.com/likexian/whois-parser-go"
	whoisRequest "github.com/moorada/whois"

	"github.com/moorada/neferpitool/pkg/configuration"
	"github.com/moorada/neferpitool/pkg/log"
)

func (w Whois) IsLike(w2 Whois) bool {
	registrantName := w.Parsed.Registrant.RegistrantName == w2.Parsed.Registrant.RegistrantName
	organization := w.Parsed.Registrant.Organization == w2.Parsed.Registrant.Organization
	creationDate := w.Parsed.Registrar.CreatedDate == w2.Parsed.Registrar.CreatedDate
	updatedDate := w.Parsed.Registrar.UpdatedDate == w2.Parsed.Registrar.UpdatedDate
	expirationDate := w.Parsed.Registrar.ExpirationDate == w2.Parsed.Registrar.ExpirationDate

	if registrantName && organization && creationDate && updatedDate && expirationDate {
		return true
	} else {
		return false
	}
}

func Get(domain string) (Whois, error, time.Duration) {

	maxAttemptsWhois := configuration.GetConf().MAXATTEMPTSWHOIS
	timeToSleepWhois := time.Duration(configuration.GetConf().TIMETOSLEEPWHOIS) * time.Millisecond
	var attempts int

	start := time.Now()

	result := Whois{}

	w, err := whoisRequest.Whois(domain)

	if err != nil {
		for attempts = 0; attempts < maxAttemptsWhois && err != nil; attempts++ {
			time.Sleep(timeToSleepWhois)
			w, err = whoisRequest.Whois(domain)
		}
		if err != nil {
			log.Debug("%s %s %s", "Error whois Request for domain: "+domain, ", Attempts: "+strconv.Itoa(attempts), err.Error())
		}
	}

	if err == nil {
		result.Raw = w
		whoisParsed, err := whoisparser.Parse(w)
		parsed := Parsed{}
		if err != nil {
			log.Debug("%s %s %s", "Error whois Parsing for domain:", domain, err.Error())

		} else {
			registrar := makeRegistrar(whoisParsed.Domain, whoisParsed.Registrar)
			registrant := makeRegistrant(whoisParsed.Registrant)
			parsed = Parsed{
				Registrant: registrant,
				Registrar:  registrar,
			}

		}
		result.Parsed = parsed
	}

	elapsed := time.Since(start)
	log.Debug("DomainName:%s,  Whois in time:%s, Attempts: %v,", domain, elapsed.String(), strconv.Itoa(attempts))
	return result, err, elapsed
}

func (w Whois) ToJson() []byte {
	whoisinfoMarshalled, err := json.Marshal(w)
	if err != nil {
		log.Error("%s", err.Error())
	}
	return whoisinfoMarshalled
}

func MakeByJson(j []byte) (Whois, error) {
	var w Whois
	err := json.Unmarshal(j, &w)
	return w, err
}

func MakeByString(s string) (Whois, error) {
	return MakeByJson([]byte(s))
}

func (w Whois) String() string {
	return string(w.ToJson())
}

func makeRegistrar(domain *whoisparser.Domain, registrarContact *whoisparser.Contact) Registrar {
	if domain == nil {
		domain = &whoisparser.Domain{}
	}
	if registrarContact == nil {
		registrarContact = &whoisparser.Contact{}
	}

	registrar := Registrar{
		RegistrarID:    registrarContact.ID,
		RegistrarName:  registrarContact.Name,
		WhoisServer:    domain.WhoisServer,
		ReferralURL:    registrarContact.ReferralURL,
		DomainId:       domain.ID,
		DomainName:     domain.Domain,
		DomainStatus:   strings.Join(domain.Status, ","),
		NameServers:    strings.Join(domain.NameServers, ","),
		DomainDNSSEC:   strconv.FormatBool(domain.DnsSec),
		CreatedDate:    domain.CreatedDate,
		UpdatedDate:    domain.UpdatedDate,
		ExpirationDate: domain.ExpirationDate,
	}
	return registrar

}

func makeRegistrant(registrantContact *whoisparser.Contact) Registrant {
	if registrantContact == nil {
		return Registrant{}
	}

	registrant := Registrant{
		RegistrantID:   registrantContact.ID,
		RegistrantName: registrantContact.Name,
		Organization:   registrantContact.Organization,
		Street:         registrantContact.Street,
		City:           registrantContact.City,
		Province:       registrantContact.Province,
		PostalCode:     registrantContact.PostalCode,
		Country:        registrantContact.Country,
		Phone:          registrantContact.Phone,
		PhoneExt:       registrantContact.PhoneExt,
		Fax:            registrantContact.Fax,
		FaxExt:         registrantContact.FaxExt,
		Email:          registrantContact.Email,
	}
	return registrant
}
