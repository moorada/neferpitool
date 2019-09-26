package domains

import (
	"errors"
	"strings"
	"time"

	"github.com/evilsocket/islazy/tui"
	"github.com/moorada/neferpitool/pkg/format"
	"github.com/moorada/neferpitool/pkg/log"

	"github.com/jinzhu/gorm"
	"github.com/moorada/neferpitool/pkg/constants"

	"github.com/moorada/neferpitool/pkg/dns"
	"github.com/moorada/neferpitool/pkg/whois"
)

type Domain struct {
	gorm.Model
	Name   string
	Status int
	dns.Dns
	whois.Whois
	ErrorWhois  string
	ErrorStatus string
	TimeWhois   string
	TimeStatus  string
	Ignore      bool
	Extra       string
}

type LegitDomain struct {
	Domain
}

/*Make a new main-DomainName*/
func NewLegitDomain(nameMainDomain string) LegitDomain {
	return LegitDomain{Domain: Domain{Name: nameMainDomain, Status: constants.UNKNOWN, Ignore: false}}
}

/*Return status as a string*/
func (d Domain) StatusToString() string {
	if d.Status == constants.INACTIVE {
		return "Inactive"
	}
	if d.Status == constants.AVAILABLE {
		return "Available"
	}
	if d.Status == constants.ACTIVE {
		return "Active"
	}
	if d.Status == constants.ALIAS {
		return "Alias"
	}
	return "Unknown"
}

/*Return status as a string colored*/
func (d Domain) StatusToStringColored() string {
	if d.Status == constants.INACTIVE {
		return tui.Yellow("Inactive")
	}
	if d.Status == constants.AVAILABLE {
		return tui.Green("Available")
	}
	if d.Status == constants.ACTIVE {
		return tui.Red("Active")
	}
	if d.Status == constants.ALIAS {
		return tui.Blue("Alias")
	}
	return "Unkwnown"
}

func (d *Domain) Update() error {

	err := d.UpdateStatus()

	if d.Status == constants.INACTIVE || d.Status == constants.ACTIVE {
		err = d.UpdateWhois()
	}
	return err
}

/*Check if the DomainName is status and update the struct*/
func (d *Domain) UpdateStatus() error {

	status, rec, err, time := dns.CheckDNS(d.Name)
	d.Dns = rec
	d.Status = status
	d.TimeStatus = time.String()

	if err == nil {
		log.Debug("DomainName: %s, Status: %s, In time:%v", d.Name, d.StatusToString(), time)
	} else {
		d.ErrorStatus = err.Error()
		log.Debug("Error Update Status abour %s :%s", d.Name, err.Error())
	}
	return err
}

/*Compare two main domains and return true is their are equals*/
func (d LegitDomain) IsEqual(d2 LegitDomain) bool {
	return d.Domain.IsEqual(d2.Domain)
}

/*Compare two main domains and return true is their are equals*/
func (d Domain) IsEqual(d2 Domain) bool {

	name := d.Name == d2.Name
	status := d.Status == d2.Status
	whois := d.Whois.IsLike(d2.Whois)

	if name && status && whois {
		return true
	} else {
		return false
	}
}

func (d Domain) GetWhois() whois.Whois {
	return d.Whois
}

func (d Domain) GetExpiryDateString() string {

	w := d.Whois

	timeString := ""
	if w.Parsed.Registrar.ExpirationDate != "" {
		expiration, err := format.StringToTime(w.Parsed.Registrar.ExpirationDate)
		if err != nil {
			log.Error(err.Error())
		}
		timeString = format.TimeToStringConsole(expiration)
	}
	return timeString
}

func (d Domain) GetExpiryDate() (expiryDate time.Time, err error) {
	/*w, err := d.GetWhois()
	if err != nil {
		log.Error(err.Error())
	}*/

	w := d.Whois

	if w.Parsed.Registrar.ExpirationDate != "" {
		expiryDate, err = format.StringToTime(w.Parsed.Registrar.ExpirationDate)
		if err != nil {
			log.Error(err.Error())
		}
		return
	} else {
		err = errors.New("Domain don't have expiry date")
	}
	return
}

/*Update all the value about whois of a main DomainName */
func (d *Domain) UpdateWhois() error {

	w, err, time := whois.Get(d.Name)
	d.TimeWhois = time.String()
	//d.Whois = string(w.ToJson())*/
	d.Whois = w

	if err == nil {
		rejectionWords := []string{
			"quora exceeded",
			"limit exceeded",
			"ERROR:101",
		}
		for _, s := range rejectionWords {
			if strings.Contains(w.Raw, s) {
				log.Debug("Whois about %s contains rejection words", d.Name)
				d.ErrorWhois = "Whois contains rejection words: " + s
				break
			}
		}
	} else {
		d.ErrorWhois = err.Error()
	}

	return err

}

/*Print a short view of a DomainName*/
func (d Domain) PrintShort() {
	log.Info("%-20s ", d.Name)
}
