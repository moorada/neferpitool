package whois

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/likexian/whois-parser-go"
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
			registrar := makeRegistrar(whoisParsed.Registrar)
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
		log.Error(err.Error())
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

func makeRegistrar(r whoisparser.Registrar) Registrar {

	registrar := Registrar{
		RegistrarID:    r.RegistrarID,
		RegistrarName:  r.RegistrarName,
		WhoisServer:    r.WhoisServer,
		ReferralURL:    r.ReferralURL,
		DomainId:       r.DomainId,
		DomainName:     r.DomainName,
		DomainStatus:   r.DomainStatus,
		NameServers:    r.NameServers,
		DomainDNSSEC:   r.DomainDNSSEC,
		CreatedDate:    r.CreatedDate,
		UpdatedDate:    r.UpdatedDate,
		ExpirationDate: r.ExpirationDate,
	}
	return registrar

}

func makeRegistrant(r whoisparser.Registrant) Registrant {

	registrant := Registrant{
		RegistrantID:   r.ID,
		RegistrantName: r.Name,
		Organization:   r.Organization,
		Street:         r.Street,
		StreetExt:      r.StreetExt,
		City:           r.City,
		Province:       r.Province,
		PostalCode:     r.PostalCode,
		Country:        r.Country,
		Phone:          r.Phone,
		PhoneExt:       r.PhoneExt,
		Fax:            r.Fax,
		FaxExt:         r.FaxExt,
		Email:          r.Email,
	}
	return registrant
}
