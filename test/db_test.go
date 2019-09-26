package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/moorada/neferpitool/pkg/constants"
	"github.com/moorada/neferpitool/pkg/db"
	"github.com/moorada/neferpitool/pkg/domains"
)

func TestAddTypoDomainsListToDBAndGetTypoDomainsListFromDB(t *testing.T) {
	_ = os.Remove(dbFile)
	db.InitDB(dbName)

	td := domains.NewTypoDomain(swappingDomain, googleDomain, algorithm)
	td.Update()

	td2 := domains.NewTypoDomain(missingDomain, googleDomain, algorithm)
	td2.Update()

	tds := domains.TypoList{td, td2}

	db.AddTypoListToDB(tds)

	tds2 := db.GetTypoDomainListFromDB(googleDomain)

	tdsMap := tds.ToMap()
	tds2Map := tds2.ToMap()

	if !(tdsMap[missingDomain].IsEqual(tds2Map[missingDomain]) && tdsMap[swappingDomain].IsEqual(tds2Map[swappingDomain])) {
		fmt.Println("/////")
		fmt.Println(tdsMap[swappingDomain])
		fmt.Println("/////")
		fmt.Println(tds2Map[swappingDomain])
		t.Errorf("getTypoDomains doesn't return the right typodomains")
	}

	_ = os.Remove(dbFile)
}

func TestGetMainDomainsListFromDB(t *testing.T) {
	_ = os.Remove(dbFile)

	db.InitDB(dbName)

	md := domains.NewLegitDomain(googleDomain)
	md.Update()
	db.AddLegitDomainToDB(md)

	md2 := domains.NewLegitDomain(stackoverflowDomain)
	md2.Update()

	db.AddLegitDomainToDB(md2)

	mdt := db.GetMainDomainListFromDB()

	mdtMap := mdt.ToMap()

	if !(mdtMap[googleDomain].IsEqual(md) && mdtMap[stackoverflowDomain].IsEqual(md2)) {
		t.Errorf("getTypoDomains doesn't return the right typodomains")
	}

	_ = os.Remove(dbFile)
}

func TestGetTypoDomainsListFromDBwithAvailability(t *testing.T) {
	_ = os.Remove(dbFile)

	db.InitDB(dbName)

	td := domains.NewTypoDomain(swappingDomain, googleDomain, algorithm)
	td.Status = constants.AVAILABLE

	db.AddTypoDomainToDB(td)

	td2 := domains.NewTypoDomain(missingDomain, googleDomain, algorithm)
	td2.Status = constants.INACTIVE

	db.AddTypoDomainToDB(td2)

	tdtRegistered := db.GetTypoDomainListWithStatusFromDB(googleDomain, []int{constants.INACTIVE})

	if len(tdtRegistered) != 1 || tdtRegistered[0].Status != constants.INACTIVE {
		t.Errorf("")
	}

	os.Remove(dbFile)
}
