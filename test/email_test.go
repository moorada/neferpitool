package test

import (
	"testing"

	"github.com/moorada/neferpitool/pkg/changes"
	"github.com/moorada/neferpitool/pkg/configuration"
	"github.com/moorada/neferpitool/pkg/domains"
	"github.com/moorada/neferpitool/pkg/notification"
)

func TestEmailChanges(t *testing.T) {
	conf := configuration.GetConf()

	if conf.EMAIL != "" && conf.PASSWORD != "" && len(conf.EMAILTONOTIFY) != 0 {

		var chs changes.ChangeList
		chs = append(chs, changes.Change{TypoDomain: swappingDomain, Field: "Status", Before: "Available", After: "Active"})
		chs = append(chs, changes.Change{TypoDomain: swappingDomain, Field: "ExpiryDate", Before: "08/12/2000", After: "02/10/2008"})
		chs = append(chs, changes.Change{TypoDomain: missingDomain, Field: "UpdateDate", Before: "08/10/2000", After: "02/10/2009"})
		chs = append(chs, changes.Change{TypoDomain: "example.com", Field: "CreateDate", Before: "08/04/2000", After: "02/10/2008"})
		chs = append(chs, changes.Change{TypoDomain: "example.com", Field: "Status", Before: "Inactive", After: "Active"})
		chs = append(chs, changes.Change{TypoDomain: "xn--yhoo-loa.com", Field: "Status", Before: "Inactive", After: "Alias"})

		td1 := domains.NewTypoDomain(swappingDomain, googleDomain, algorithm)
		td2 := domains.NewTypoDomain(missingDomain, googleDomain, algorithm)
		td3 := domains.NewTypoDomain(wtldDomain, googleDomain, algorithm)
		td1.Update()
		td2.Update()
		td3.Update()

		headersAva, datasAva, headersWhois, datasWhois := chs.ToTables()

		hExpiry, dExpiry := domains.TypoList([]domains.TypoDomain{td1, td2, td3}).ToExpiryTable()

		tpl := notification.TemplateData{
			H1:            "Domains Monitoring",
			TextStatus:    "There are status changes",
			TextWhois:     "There are whois changes",
			HeadersStatus: headersAva,
			HeadersWhois:  headersWhois,
			DatasStatus:   datasAva,
			DatasWhois:    datasWhois,
			TextExpiry:    "Typodomains in expiration",
			HeadersExpiry: hExpiry,
			DatasExpiry:   dExpiry,
		}

		conf := configuration.GetConf()
		request := notification.Request{
			From:     conf.EMAIL,
			Password: conf.PASSWORD,
			To:       conf.EMAILTONOTIFY,
			Subject:  "TEST - domain monitoring",
		}

		err := notification.EmailChanges(tpl, request)

		if err != nil {
			t.Errorf("Error to send Email: %s", err.Error())
		}
	}
}
