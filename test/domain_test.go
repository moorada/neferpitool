package test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/moorada/neferpitool/pkg/constants"
	"github.com/moorada/neferpitool/pkg/domains"
	"github.com/moorada/neferpitool/pkg/whois"
)

func TestUpdateStatus(t *testing.T) {
	td := domains.NewTypoDomain(googleDomain, googleDomain, algorithm)
	_ = td.UpdateStatus()
	if td.Status != constants.ACTIVE {
		t.Errorf("Expected that google.com is registered, but got %v", td.StatusToString())
	}

	td = domains.NewTypoDomain(notExistingDomain, googleDomain, algorithm)
	_ = td.UpdateStatus()
	if td.Status != constants.AVAILABLE {
		t.Errorf("Expected that google.com is available, but got %v", td.StatusToString())
	}

}

func TestCheckForChanges(t *testing.T) {
	td := domains.NewTypoDomain(googleDomain, googleDomain, algorithm)
	err := td.Update()
	if err != nil {
		t.Errorf("Error td.Update %s", err.Error())
	}
	b, _ := td.IsChanged()
	if b {
		t.Errorf("Expected that %s is not changed, but got :%v", googleDomain, b)
	}

	td.Status = constants.AVAILABLE
	b, _ = td.IsChanged()
	if !b {
		t.Errorf("Expected that %s is changed, but got :%v", googleDomain, b)
	}

}

func TestUpdateWhois(t *testing.T) {
	td := domains.NewTypoDomain(googleDomain, googleDomain, algorithm)
	_ = td.UpdateWhois()

	if cmp.Equal(td.Whois, whois.Whois{}) {
		t.Errorf("Expected that whois google.com is not a empty struct, but got empty")
	}

	td = domains.NewTypoDomain(notExistingDomain, googleDomain, algorithm)
	_ = td.UpdateWhois()
	if cmp.Equal(td.Whois, whois.Whois{}) {
		t.Errorf("Expected that whois %s is a empty struct, but got empty", notExistingDomain)
	}

}

func TestTyposAreEquals(t *testing.T) {

	td := domains.NewTypoDomain(swappingDomain, googleDomain, algorithm)
	err := td.Update()
	if err != nil {
		t.Errorf("Error td.Update %s", err.Error())

	}

	td2 := domains.NewTypoDomain(swappingDomain, googleDomain, algorithm)
	err = td2.Update()
	if err != nil {
		t.Errorf("Error td.Update %s", err.Error())

	}

	tDifferent := domains.NewTypoDomain(missingDomain, googleDomain, algorithm)
	err = tDifferent.Update()
	if err != nil {
		t.Errorf("Error td.Update %s", err.Error())

	}

	if !td.IsEqual(td2) {
		t.Errorf("Expected that typodomains structs are egual but: %v", td.IsEqual(td2))
	}

	if td.IsEqual(tDifferent) {
		t.Errorf("Expected that typodomains structs are different but: %v", td.IsEqual(tDifferent))
	}

}
