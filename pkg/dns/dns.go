package dns

import (
	"errors"
	"net"
	"time"

	"github.com/miekg/dns"

	"github.com/moorada/neferpitool/pkg/configuration"
	"github.com/moorada/neferpitool/pkg/constants"
	"github.com/moorada/neferpitool/pkg/log"
)

type Dns struct {
	SOA   string
	CNAME string
	A     string
	AAAA  string
	MX    string
	NS    string
}

func (r Dns) IsEqual(r2 Dns) bool {

	soa := r.SOA == r2.SOA
	cname := r.CNAME == r2.CNAME
	a := r.A == r2.A
	aaaa := r.AAAA == r2.AAAA
	mx := r.MX == r2.MX
	ns := r.NS == r2.NS

	if soa && cname && a && aaaa && mx && ns {
		return true
	} else {
		return false
	}
}
func CheckDNS(d string) (result int, info Dns, err error, duration time.Duration) {
	start := time.Now()
	result = constants.UNKNOWN

	mas := configuration.GetConf().MAXATTEMPTSSOA
	tss := configuration.GetConf().TIMETOSLEEPSOA
	info = Dns{}

	errorsMap := make(map[string]error)

	respSOA, soa, errSoa := iterateRequest(d, dns.TypeSOA, mas, tss)
	info.SOA = soa
	errorsMap["soa"] = errSoa

	respNS, ns, errNS := iterateRequest(d, dns.TypeNS, mas, tss)
	info.NS = ns
	errorsMap["ns"] = errNS

	if errSoa == nil || errNS == nil {
		if respSOA || respNS {
			result = constants.INACTIVE
		} else {
			result = constants.AVAILABLE
		}
	}
	respCNAME, cname, errCNAME := iterateRequest(d, dns.TypeCNAME, mas, tss)
	info.CNAME = cname
	errorsMap["cname"] = errCNAME
	if errCNAME == nil {
		if respCNAME {
			result = constants.ALIAS
		}
	}

	respA, a, errA := iterateRequest(d, dns.TypeA, mas, tss)
	info.A = a
	errorsMap["a"] = errA
	respAAAA, aaaa, errAAAA := iterateRequest(d, dns.TypeAAAA, mas, tss)
	info.AAAA = aaaa
	errorsMap["aaaa"] = errAAAA
	respMX, mx, errMX := iterateRequest(d, dns.TypeMX, mas, tss)
	info.MX = mx
	errorsMap["mx"] = errMX

	if respA || respAAAA || respMX {
		result = constants.ACTIVE
	}

	var errString string
	for k, e := range errorsMap {
		if e != nil {
			errString += k + ": " + e.Error() + " | "
		}
	}
	if errString != "" {
		err = errors.New(errString)
	}

	duration = time.Since(start)
	return
}

func iterateRequest(domain string, recordType uint16, maxAttempts int, timesleep int) (bool, string, error) {
	resp, vl, er := isThereRecord(domain, recordType)
	for i := 0; i < maxAttempts && er != nil; i++ {
		resp, vl, er = isThereRecord(domain, recordType)
		time.Sleep(time.Duration(timesleep) * time.Millisecond)
	}
	return resp, vl, er
}

func isThereRecord(d string, recordType uint16) (ok bool, recordValue string, err error) {
	resolver := configuration.GetConf().PATHRESOLVER
	config, err := dns.ClientConfigFromFile(resolver)
	if err != nil {
		config, err = dns.ClientConfigFromFile("/etc/resolv.conf")
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(d), recordType)
	m.RecursionDesired = true
	r, _, err := c.Exchange(m, net.JoinHostPort(config.Servers[0], config.Port))
	if err != nil {
		log.Debug("Error Answer RecordType: %v,  about %s: %s", recordType, d, err.Error())
		return false, "", err
	} else {
		if r.Answer != nil {
			log.Debug("Answer RecordType: %v,  about %s ok ", recordType, d)
			return true, r.Answer[0].String(), nil
		} else {
			log.Debug("Answer RecordType: %v, empty about %s", recordType, d)
			return false, "", nil
		}

	}
}
