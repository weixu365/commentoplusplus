package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type OwnerRepoTestSuite struct {
	suite.Suite
	repo OwnerRepository
}

func TestOwnerRepoTestSuite(t *testing.T) {
	suite.Run(t, new(OwnerRepoTestSuite))
}

func (suite *OwnerRepoTestSuite) SetupTest() {
	SetupTestRepo()
	suite.repo = Repo.OwnerRepository
}

func (suite *OwnerRepoTestSuite) Test_get_page_by_path() {
	suite.repo.CreatePage("example.com", "/path.html")

	page, err := suite.repo.GetPageByPath("example.com", "/path.html")
	suite.Require().Nil(err)

	assert.False(suite.T(), page.IsLocked)
}

func (suite *OwnerRepoTestSuite) TestOwnerGetByEmailBasics(t *testing.T) {
	ownerHex, _ := ownerNew("test@example.com", "Test", "hunter2")

	o, err := suite.repo.GetByEmail("test@example.com")
	if err != nil {
		t.Errorf("unexpected error on ownerGetByEmail: %v", err)
		return
	}

	if o.OwnerHex != ownerHex {
		t.Errorf("expected ownerHex=%s got ownerHex=%s", ownerHex, o.OwnerHex)
		return
	}
}

func (suite *OwnerRepoTestSuite) TestOwnerGetByEmailDNE(t *testing.T) {
	if _, err := suite.repo.GetByEmail("invalid@example.com"); err == nil {
		t.Errorf("expected error not found on ownerGetByEmail before creating an account")
		return
	}
}

func (suite *OwnerRepoTestSuite) TestOwnerGetByOwnerTokenBasics(t *testing.T) {
	ownerHex, _ := ownerNew("test@example.com", "Test", "hunter2")

	ownerToken, _ := ownerLogin("test@example.com", "hunter2")

	o, err := ownerGetByOwnerToken(ownerToken)
	if err != nil {
		t.Errorf("unexpected error on ownerGetByOwnerToken: %v", err)
		return
	}

	if o.OwnerHex != ownerHex {
		t.Errorf("expected ownerHex=%s got ownerHex=%s", ownerHex, o.OwnerHex)
		return
	}
}

func (suite *OwnerRepoTestSuite) TestOwnerGetByOwnerTokenDNE(t *testing.T) {
	if _, err := ownerGetByOwnerToken("does-not-exist"); err == nil {
		t.Errorf("expected error not found on ownerGetByOwnerToken before creating an account")
		return
	}
}
