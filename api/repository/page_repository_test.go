package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PageRepoTestSuite struct {
	suite.Suite
	repo PageRepository
}

func (suite *PageRepoTestSuite) SetupTest() {
	SetupTestRepo()
	suite.repo = Repo.PageRepository
}

func (suite *PageRepoTestSuite) Test_get_page_by_path() {
	suite.repo.CreatePage("example.com", "/path.html")

	page, err := suite.repo.GetPageByPath("example.com", "/path.html")
	suite.Require().Nil(err)

	assert.False(suite.T(), page.IsLocked)
}

func (suite *PageRepoTestSuite) Test_fail_to_get_page_when_domain_is_empty() {
	suite.repo.CreatePage("example.com", "")

	_, err := suite.repo.GetPageByPath("", "/path.html")
	suite.Require().NotNil(err)
}

func (suite *PageRepoTestSuite) Test_able_to_get_page_when_path_is_empty() {
	suite.repo.CreatePage("example.com", "")

	_, err := suite.repo.GetPageByPath("example.com", "")
	suite.Nil(err)
}

func (suite *PageRepoTestSuite) Test_able_to_get_page_when_path_does_not_exist() {
	suite.repo.CreatePage("example.com", "")

	_, err := suite.repo.GetPageByPath("example.com", "/non-exists-path.html")
	suite.Nil(err)
}

func (suite *PageRepoTestSuite) Test_fail_to_save_page_when_domain_is_empty() {
	err := suite.repo.CreatePage("", "/path.html")
	suite.Require().NotNil(err)
}

func (suite *PageRepoTestSuite) Test_skip_saving_page_when_the_same_path_exists() {
	err := suite.repo.CreatePage("example.com", "/path.html")
	suite.Require().Nil(err)

	err = suite.repo.CreatePage("example.com", "/path.html")
	suite.Require().Nil(err)
}

func TestPageRepoTestSuite(t *testing.T) {
	suite.Run(t, new(PageRepoTestSuite))
}
