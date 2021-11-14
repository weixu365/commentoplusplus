package handler

import (
	"simple-commenting/test"
	"testing"
)

func TestDomainGetBasics(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	domainNew("temp-owner-hex", "Example", "example.com")

	domain, err := domainGet("example.com")
	if err != nil {
		t.Errorf("unexpected error getting domain: %v", err)
		return
	}

	if domain.Name != "Example" {
		t.Errorf("expected name=Example got name=%s", domain.Name)
		return
	}
}

func TestDomainGetEmpty(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	if _, err := domainGet(""); err == nil {
		t.Errorf("expected error not found when getting with empty domain")
		return
	}
}

func TestDomainGetDNE(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	if _, err := domainGet("example.com"); err == nil {
		t.Errorf("expected error not found when getting non-existant domain")
		return
	}
}
