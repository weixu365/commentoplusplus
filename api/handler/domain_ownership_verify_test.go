package handler

import (
	"simple-commenting/test"
	"testing"
)

func TestDomainVerifyOwnershipBasics(t *testing.T) {
	test.FailTestOnError(t, test.SetupTestEnv())

	ownerHex, _ := ownerNew("test@example.com", "Test", "hunter2")
	ownerLogin("test@example.com", "hunter2")

	domainNew(ownerHex, "Example", "example.com")

	isOwner, err := domainOwnershipVerify(ownerHex, "example.com")
	if err != nil {
		t.Errorf("error checking ownership: %v", err)
		return
	}

	if !isOwner {
		t.Errorf("expected isOwner=true got isOwner=false")
		return
	}

	otherOwnerHex, _ := ownerNew("test2@example.com", "Test2", "hunter2")
	ownerLogin("test2@example.com", "hunter2")

	isOwner, err = domainOwnershipVerify(otherOwnerHex, "example.com")
	if err != nil {
		t.Errorf("error checking ownership: %v", err)
		return
	}

	if isOwner {
		t.Errorf("expected isOwner=false got isOwner=true")
		return
	}
}
