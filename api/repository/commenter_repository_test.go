package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CommenterRepoTestSuite struct {
	suite.Suite
	repo CommenterRepository
}

func TestCommenterRepoTestSuite(t *testing.T) {
	suite.Run(t, new(CommenterRepoTestSuite))
}

func (suite *CommenterRepoTestSuite) SetupTest() {
	SetupTestRepo()
	suite.repo = Repo.CommenterRepository
}

func (suite *CommenterRepoTestSuite) TestCommenterGetByHexBasics() {
	commenterHex, _ := commenterNew("test@example.com", "Test", "undefined", "https://example.com/photo.jpg", "google", "")

	c, err := suite.repo.GetCommenterByHex(commenterHex)

	suite.Require().Nil(err)
	assert.Equal(suite.T(), c.Name, "Test")
}

func (suite *CommenterRepoTestSuite) TestCommenterGetByHexEmpty() {
	_, err := suite.repo.GetCommenterByHex("")

	suite.Require().NotNil(err)
}

func (suite *CommenterRepoTestSuite) TestCommenterGetByCommenterToken() {
	commenterHex, _ := commenterNew("test@example.com", "Test", "undefined", "https://example.com/photo.jpg", "google", "")
	commenterToken, _ := commenterTokenNew()
	commenterSessionUpdate(commenterToken, commenterHex)

	c, err := suite.repo.GetCommenterByToken(commenterToken)

	suite.Require().Nil(err)
	assert.Equal(suite.T(), c.Name, "Test")

}

func (suite *CommenterRepoTestSuite) TestCommenterGetByCommenterTokenEmpty() {
	_, err := suite.repo.GetCommenterByToken("")
	suite.Require().NotNil(err)
}

func (suite *CommenterRepoTestSuite) TestCommenterGetByName() {
	commenterHex, _ := commenterNew("test@example.com", "Test", "undefined", "https://example.com/photo.jpg", "google", "")
	commenterToken, _ := commenterTokenNew()
	commenterSessionUpdate(commenterToken, commenterHex)

	c, err := suite.repo.GetCommenterByEmail("google", "test@example.com")

	suite.Require().Nil(err)
	assert.Equal(suite.T(), c.Name, "Test")
}

func (suite *CommenterRepoTestSuite) TestCommenterGetByNameEmpty() {
	_, err := suite.repo.GetCommenterByEmail("", "")

	suite.Require().NotNil(err)
}
