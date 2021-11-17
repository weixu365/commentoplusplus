package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DomainRepoTestSuite struct {
	suite.Suite
	repo DomainRepository
}

func TestDomainRepoTestSuite(t *testing.T) {
	suite.Run(t, new(DomainRepoTestSuite))
}

func (suite *DomainRepoTestSuite) SetupTest() {
	SetupTestRepo()
	suite.repo = Repo.DomainRepository
}

func (suite *DomainRepoTestSuite) TestDomainUpdateBasics() {
	Repo.DomainRepository.CreateDomain("temp-owner-hex", "Example", "example.com")

	domain, _ := Repo.DomainRepository.ListDomain("temp-owner-hex")

	domain[0].Name = "Example2"
	err := Repo.DomainRepository.UpdateDomain(domain[0])
	suite.Require().NoError(err)

	domain, _ = Repo.DomainRepository.ListDomain("temp-owner-hex")
	assert.Equal(suite.T(), "Example2", domain[0].Name)
}
