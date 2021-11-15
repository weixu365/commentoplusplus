package handler

import (
	"simple-commenting/repository"
	"simple-commenting/test"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CommentNewTestSuite struct {
	suite.Suite
	commentRepo   repository.CommentRepository
	commenterRepo repository.CommenterRepository
}

func TestCommentNewTestSuite(t *testing.T) {
	suite.Run(t, new(CommentNewTestSuite))
}

func (suite *CommentNewTestSuite) SetupTest() {
	test.FailTestOnError(suite.T(), test.SetupTestEnv())

	suite.commentRepo = repository.Repo.CommentRepository
	suite.commenterRepo = repository.Repo.CommenterRepository
}

func (suite *CommentNewTestSuite) TestCommentNewBasics() {
	_, err := commentNew("temp-commenter-hex", "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC())

	suite.Require().NoError(err)
}

func (suite *CommentNewTestSuite) TestCommentNewEmpty() {
	_, err := commentNew("temp-commenter-hex", "example.com", "", "root", "**foo**", "approved", time.Now().UTC())
	suite.Require().NoError(err, "empty path is allowed")

	_, err = commentNew("temp-commenter-hex", "", "", "root", "**foo**", "approved", time.Now().UTC())
	suite.Require().Error(err, "expected error not found creatingn new comment with empty domain")

	_, err = commentNew("", "", "", "", "", "", time.Now().UTC())
	suite.Require().Error(err, "expected error not found creatingn new comment with empty everything")
}

func (suite *CommentNewTestSuite) TestCommentNewUpvoted() {
	commentHex, err := commentNew("temp-commenter-hex", "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC())
	suite.Require().NoError(err)

	comment, err := suite.commentRepo.GetByCommentHex(commentHex)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), 0, comment.Score)
}

func (suite *CommentNewTestSuite) TestCommentNewThreadLocked() {
	repository.Repo.PageRepository.CreatePage("example.com", "/path.html")
	p, _ := pageGet("example.com", "/path.html")
	p.IsLocked = true
	pageUpdate(p)

	_, err := commentNew("temp-commenter-hex", "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC())
	suite.Require().Error(err, "expected error not found creating a new comment on a locked thread")
}

func (suite *CommentNewTestSuite) TestCommentDomainPathGetBasics() {
	commentHex, _ := commentNew("temp-commenter-hex", "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC())

	domain, path, err := repository.Repo.CommentRepository.GetCommentDomainPath(commentHex)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), domain, "example.com")
	assert.Equal(suite.T(), path, "/path.html")
}

func (suite *CommentNewTestSuite) TestCommentCountBasics() {
	commenterHex, _ := suite.commenterRepo.CreateCommenter("test@example.com", "Test", "undefined", "http://example.com/photo.jpg", "google", "")

	commentNew(commenterHex, "example.com", "/path.html", "root", "**foo**", "approved", time.Now().UTC())
	commentNew(commenterHex, "example.com", "/path.html", "root", "**bar**", "approved", time.Now().UTC())
	commentNew(commenterHex, "example.com", "/path.html", "root", "**baz**", "unapproved", time.Now().UTC())

	counts, err := suite.commentRepo.GetCommentsCount("example.com", []string{"/path.html"})
	suite.Require().NoError(err)

	assert.Equal(suite.T(), 3, counts["/path.html"])
}

func (suite *CommentNewTestSuite) TestCommentCountNewPage() {
	counts, err := suite.commentRepo.GetCommentsCount("example.com", []string{"/path.html"})
	suite.Require().NoError(err)

	assert.Equal(suite.T(), 0, counts["/path.html"])
}

func (suite *CommentNewTestSuite) TestCommentCountEmpty() {
	_, err := suite.commentRepo.GetCommentsCount("example.com", []string{""})
	suite.Require().NoError(err)

	_, err = suite.commentRepo.GetCommentsCount("", []string{""})
	suite.Require().Error(err, "expected error not found counting comments with empty everything")
}
