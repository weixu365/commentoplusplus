package handler

import (
	"simple-commenting/test"
	"testing"
)

func TestDomainNewBasics(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	if err := domainNew("temp-owner-hex", "Example", "example.com"); err != nil {
		t.Errorf("unexpected error creating domain: %v", err)
		return
	}
}

func TestDomainNewClash(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	if err := domainNew("temp-owner-hex", "Example", "example.com"); err != nil {
		t.Errorf("unexpected error creating domain: %v", err)
		return
	}

	if err := domainNew("temp-owner-hex", "Example 2", "example.com"); err == nil {
		t.Errorf("expected error not found when creating with clashing domain")
		return
	}
}

func TestDomainNewEmpty(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	if err := domainNew("temp-owner-hex", "Example", ""); err == nil {
		t.Errorf("expected error not found when creating with emtpy domain")
		return
	}

	if err := domainNew("", "", ""); err == nil {
		t.Errorf("expected error not found when creating with emtpy everything")
		return
	}
}
