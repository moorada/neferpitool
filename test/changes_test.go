package test

import (
	"testing"

	"github.com/moorada/neferpitool/pkg/changes"
	"github.com/moorada/neferpitool/pkg/constants"
	"github.com/moorada/neferpitool/pkg/domains"
)

func TestCheckReliabilityWithPrev(t *testing.T) {
	td1 := domains.NewTypoDomain(missingDomain, googleDomain, algorithm)
	td1.Status = constants.INACTIVE
	tds1 := domains.TypoList{td1}
	change := changes.Change{
		TypoDomain: missingDomain,
		Field:      changes.STATUS,
		Before:     "Available",
		After:      "Registered",
	}
	tdcs1 := changes.ChangeList{change}

	td2 := domains.NewTypoDomain(missingDomain, googleDomain, algorithm)
	td1.Status = constants.INACTIVE
	tds2 := domains.TypoList{td2}
	change2 := changes.Change{
		TypoDomain: missingDomain,
		Field:      changes.STATUS,
		Before:     "Available",
		After:      "Registered",
	}
	tdcs2 := changes.ChangeList{change2}

	td3 := domains.NewTypoDomain(missingDomain, googleDomain, algorithm)
	td1.Status = constants.AVAILABLE
	tds3 := domains.TypoList{td3}

	change3 := changes.Change{
		TypoDomain: missingDomain,
		Field:      changes.STATUS,
		Before:     "Registered",
		After:      "Available",
	}

	tdcs3 := changes.ChangeList{change3}

	tdsChecked, changesChecked := tdcs1.FilterReliableWithPrev(tdcs2, tds1, tds2)
	if len(tdsChecked) != 1 && len(changesChecked) != 1 {
		t.Errorf("Expected that CheckReliability return the changes because are Reliability but it doesn't work")
	}

	tdsChecked, changesChecked = tdcs1.FilterReliableWithPrev(tdcs3, tds1, tds3)
	if len(tdsChecked) != 0 && len(changesChecked) != 0 {
		t.Errorf("Expected that CheckReliability return 0 changes because aren't Reliability but it doesn't work")
	}

}
