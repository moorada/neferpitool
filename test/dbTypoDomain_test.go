package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/moorada/neferpitool/pkg/constants"
	"github.com/moorada/neferpitool/pkg/db"
	"github.com/moorada/neferpitool/pkg/domains"
	"github.com/moorada/neferpitool/pkg/whois"
)

func TestAddTypoDomainToDBAndGetTypoDomainFromDB(t *testing.T) {

	os.Remove(dbFile)

	db.InitDB(dbName)
	td := domains.NewTypoDomain(swappingDomain, googleDomain, algorithm)
	td.Status = constants.INACTIVE

	w, _, _ := whois.Get("goolge.it")

	//td.WhoisInfo = w.String()
	td.Whois = w

	db.AddTypoDomainToDB(td)

	td2 := db.GetTypoDomainFromDB(swappingDomain)

	if !td.IsEqual(td2) {
		t.Errorf("Expected that typodomain added to database and typodomain got are egual but: %v", td.IsEqual(td2))
		fmt.Println("td1", td.Name)
		fmt.Println()
		fmt.Println("td2", td2.Name)
	}

	db.RemoveTypoDomainFromDB(swappingDomain)

	if db.GetTypoDomainFromDB(swappingDomain) != (domains.TypoDomain{}) {
		t.Errorf("Expected that TypoDomain is removed but it is in DB, name: %s", td.Name)
	}

	_ = os.Remove(dbFile)

}

func TestRemoveTypoDomainFromDB(t *testing.T) {

	_ = os.Remove(dbFile)

	db.InitDB(dbName)
	td := domains.NewTypoDomain(swappingDomain, googleDomain, algorithm)
	td.Update()

	w, _, _ := whois.Get("goolge.it")

	td.Whois = w /*.String()*/

	db.AddTypoDomainToDB(td)

	td2 := db.GetTypoDomainFromDB(swappingDomain)

	if !td.IsEqual(td2) {
		t.Errorf("Expected that typodomain added to database and typodomain got are egual but: %v", td.IsEqual(td2))
		fmt.Println("td1", td.Name)
		fmt.Println()
		fmt.Println("td2", td2.Name)
	}

	_ = os.Remove(dbFile)
}
