package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DomainModeratorRepoTestSuite struct {
	suite.Suite
	repo DomainModeratorRepository
}

func TestDomainModeratorRepoTestSuite(t *testing.T) {
	suite.Run(t, new(DomainModeratorRepoTestSuite))
}

func (suite *DomainModeratorRepoTestSuite) SetupTest() {
	SetupTestRepo()
	suite.repo = Repo.DomainModeratorRepository
}

func (suite *DomainModeratorRepoTestSuite) TestGetModeratorsForDomain() {
	suite.repo.CreateModerator("example.com", "test@example.com")
	suite.repo.CreateModerator("example.com", "test2@example.com")

	mods, err := suite.repo.GetModeratorsForDomain("example.com")
	suite.Nil(err)

	assert.Equal(suite.T(), len(*mods), 2)
	assert.Equal(suite.T(), (*mods)[0].Email, "test@example.com")
	assert.Equal(suite.T(), (*mods)[1].Email, "test2@example.com")
}

func (suite *DomainModeratorRepoTestSuite) TestIsDomainModerators() {
	suite.repo.CreateModerator("example.com", "test@example.com")

	isMod, err := suite.repo.IsDomainModerator("example.com", "test@example.com")
	suite.Nil(err)
	assert.True(suite.T(), isMod)

	isMod, err = suite.repo.IsDomainModerator("example.com", "test2@example.com")
	suite.Nil(err)
	assert.False(suite.T(), isMod)
}
