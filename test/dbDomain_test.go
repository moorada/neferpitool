package test

import (
	"os"
	"testing"

	"github.com/moorada/neferpitool/pkg/db"
	"github.com/moorada/neferpitool/pkg/domains"
	"github.com/moorada/neferpitool/pkg/log"
)

func TestAddDomainToDBAndGetDomainFromDBAndRemoveDomainFromDB(t *testing.T) {

	log.ActiveConsoleLog()
	os.Remove(dbFile)

	db.InitDB(dbName)
	d := domains.NewLegitDomain(googleDomain)
	d.Update()

	db.AddLegitDomainToDB(d)
	dGet, _ := db.GetLegitDomainFromDB(googleDomain)
	if !d.IsEqual(dGet) {
		t.Errorf("Expected that DomainName added to database and DomainName got are egual but: %v", d.IsEqual(dGet))
	}

	db.RemoveLegitDomainFromDB(googleDomain)

	dGet, _ = db.GetLegitDomainFromDB(googleDomain)
	if dGet != (domains.LegitDomain{}) {
		t.Errorf("Expected that DomainName is removed but it is in db, name: %s", d.Name)
	}

	_, err := db.GetLegitDomainFromDB(stackoverflowDomain)
	if err == nil {
		t.Errorf("Expected that GetLegitDomainFromDB return error because %s is not in db", stackoverflowDomain)
	}

	db.AddLegitDomainToDB(d)
	db.AddTypoDomainToDB(domains.NewTypoDomain(missingDomain, googleDomain, algorithm))
	db.RemoveLegitDomainFromDB(googleDomain)

	tds := db.GetTypoDomainListFromDB(googleDomain)
	if len(tds) != 0 {
		t.Errorf("Expected that TypoDomains about %s are removed but they are in db", googleDomain)
	}

	if db.GetTypoDomainFromDB(missingDomain) != (domains.TypoDomain{}) {
		t.Errorf("Expected that TypoDomain is removed but it is in DB, name: %s", d.Name)
	}

	_ = os.Remove(dbFile)
}
