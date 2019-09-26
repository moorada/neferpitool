package test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/moorada/neferpitool/pkg/whois"
)

func TestMakeByStringAndString(t *testing.T) {

	w, err, _ := whois.Get(googleDomain)
	if err != nil {
		t.Errorf("Whois.Get doesn't work: %s", err.Error())
	}

	s := w.String()

	w2, err := whois.MakeByString(s)
	if err != nil {
		t.Errorf("MakeByString: %s", err.Error())
	}

	if !cmp.Equal(w, w2) {
		t.Errorf("Expected that whois structs are egual but not")
	}

	wDifferent, err, _ := whois.Get(stackoverflowDomain)
	if err != nil {
		t.Errorf("MakeByString: %s", err.Error())
	}
	if cmp.Equal(w, wDifferent) {
		t.Errorf("Expected that whois structs are different but not")
	}

	wEmpty := ""
	if cmp.Equal(w, wEmpty) {
		t.Errorf("Expected that whois structs are different but not")
	}

}
