package handler

import (
	"simple-commenting/repository"
	"simple-commenting/test"
	"testing"
)

func TestDomainDeleteBasics(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	domainNew("temp-owner-hex", "Example", "example.com")
	domainNew("temp-owner-hex", "Example", "example2.com")

	if err := domainDelete("example.com"); err != nil {
		t.Errorf("unexpected error deleting domain: %v", err)
		return
	}

	d, _ := repository.Repo.DomainRepository.ListDomain("temp-owner-hex")

	if len(d) != 1 {
		t.Errorf("expected number of domains to be 1 got %d", len(d))
		return
	}

	if d[0].Domain != "example2.com" {
		t.Errorf("expected first domain to be example2.com got %s", d[0].Domain)
		return
	}
}

func TestDomainDeleteEmpty(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	if err := domainDelete(""); err == nil {
		t.Errorf("expected error not found when deleting with empty domain")
		return
	}
}
