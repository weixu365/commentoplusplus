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
	commenterHex, _ := suite.repo.CreateCommenter("test@example.com", "Test", "undefined", "https://example.com/photo.jpg", "google", "")

	c, err := suite.repo.GetCommenterByHex(commenterHex)

	suite.Require().Nil(err)
	assert.Equal(suite.T(), c.Name, "Test")
}

func (suite *CommenterRepoTestSuite) TestCommenterGetByHexEmpty() {
	_, err := suite.repo.GetCommenterByHex("")

	suite.Require().NotNil(err)
}

func (suite *CommenterRepoTestSuite) TestCommenterGetByCommenterToken() {
	commenterHex, _ := suite.repo.CreateCommenter("test@example.com", "Test", "undefined", "https://example.com/photo.jpg", "google", "")
	commenterToken, _ := suite.repo.CreateCommenterToken()
	suite.repo.UpdateCommenterSession(commenterToken, commenterHex)

	c, err := suite.repo.GetCommenterByToken(commenterToken)

	suite.Require().Nil(err)
	assert.Equal(suite.T(), c.Name, "Test")
}

func (suite *CommenterRepoTestSuite) TestCommenterGetByCommenterTokenEmpty() {
	_, err := suite.repo.GetCommenterByToken("")
	suite.Require().NotNil(err)
}

func (suite *CommenterRepoTestSuite) TestCommenterGetByName() {
	commenterHex, _ := suite.repo.CreateCommenter("test@example.com", "Test", "undefined", "https://example.com/photo.jpg", "google", "")
	commenterToken, _ := suite.repo.CreateCommenterToken()
	suite.repo.UpdateCommenterSession(commenterToken, commenterHex)

	c, err := suite.repo.GetCommenterByEmail("google", "test@example.com")

	suite.Require().Nil(err)
	assert.Equal(suite.T(), c.Name, "Test")
}

func (suite *CommenterRepoTestSuite) TestCommenterGetByNameEmpty() {
	_, err := suite.repo.GetCommenterByEmail("", "")

	suite.Require().NotNil(err)
}

func (suite *CommenterRepoTestSuite) TestCommenterTokenNewBasics() {
	_, err := suite.repo.CreateCommenterToken()

	suite.Require().NoError(err)
}

func (suite *CommenterRepoTestSuite) TestCommenterSessionUpdateBasics() {
	commenterToken, _ := suite.repo.CreateCommenterToken()

	err := suite.repo.UpdateCommenterSession(commenterToken, "temp-commenter-hex")
	suite.Require().NoError(err)

	commenterHex, err := suite.repo.GetCommenterHex(commenterToken)

	suite.Require().NoError(err)
	assert.Equal(suite.T(), commenterHex, "temp-commenter-hex")
}

func (suite *CommenterRepoTestSuite) TestCommenterSessionUpdateEmpty() {
	err := suite.repo.UpdateCommenterSession("", "temp-commenter-hex")

	suite.Require().Error(err)
}
