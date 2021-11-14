package repository

import (
	"simple-commenting/app"
	"simple-commenting/model"
	"simple-commenting/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CommentRepoTestSuite struct {
	suite.Suite
	repo          CommentRepository
	commenterRepo CommenterRepository
}

func TestCommentRepoTestSuite(t *testing.T) {
	suite.Run(t, new(CommentRepoTestSuite))
}

func (suite *CommentRepoTestSuite) SetupTest() {
	SetupTestRepo()
	suite.repo = Repo.CommentRepository
	suite.commenterRepo = Repo.CommenterRepository
}

func (suite *CommentRepoTestSuite) CreateComment(commentHex, commenterHex, parentHex, state string) *model.Comment {
	commentHex1, _ := util.RandomHex(32)
	return &model.Comment{
		CommentHex:   commentHex1,
		CommenterHex: commenterHex,
		Domain:       "example.com",
		Path:         "/path.html",
		ParentHex:    parentHex,
		Markdown:     "**foo**",
		Html:         "html text",
		State:        "unapproved",
		CreationDate: time.Now().UTC(),
	}
}
func (suite *CommentRepoTestSuite) NewComment(commenterHex string, domainName string, path string, parentHex string, markdown, html string, state string, creationDate time.Time) *model.Comment {
	commentHex1, _ := util.RandomHex(32)
	return &model.Comment{
		CommentHex:   commentHex1,
		CommenterHex: commenterHex,
		Domain:       domainName,
		Path:         path,
		ParentHex:    parentHex,
		Markdown:     markdown,
		Html:         html,
		State:        state,
		CreationDate: time.Now().UTC(),
	}
}
func (suite *CommentRepoTestSuite) Test_approve_comment() {
	commenterHex, _ := suite.commenterRepo.CreateCommenter("test@example.com", "Test", "undefined", "https://example.com/photo.jpg", "google", "")
	commentHex, err := util.RandomHex(32)
	suite.Require().NoError(err)

	_, err = suite.repo.CreateComment(&model.Comment{
		CommentHex:   commentHex,
		CommenterHex: commenterHex,
		Domain:       "example.com",
		Path:         "/path.html",
		ParentHex:    "root",
		Markdown:     "**foo**",
		Html:         "html text",
		State:        "unapproved",
		CreationDate: time.Now().UTC(),
	})
	suite.Require().NoError(err)

	err = suite.repo.ApproveComment(commentHex, "/path.html")
	suite.Require().NoError(err)

	comments, _, err := suite.repo.ListComments("anonymous", "example.com", "/path.html", true)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), comments[0].State, "approved")
}

func (suite *CommentRepoTestSuite) Test_ApproveComment_return_error_when_commentHex_is_empty() {
	err := suite.repo.ApproveComment("", "/any-path.html")

	suite.Require().Error(err)
}

func (suite *CommentRepoTestSuite) Test_Delete_Comment_also_delete_replies() {
	commenterHex := "temp-commenter-hex"
	commentHex1, _ := util.RandomHex(32)
	commentHex2, _ := util.RandomHex(32)
	suite.repo.CreateComment(suite.CreateComment(commentHex1, commenterHex, "root", "approved"))
	suite.repo.CreateComment(suite.CreateComment(commentHex2, commenterHex, commentHex1, "approved"))

	err := suite.repo.DeleteComment(commentHex1, commenterHex, "example.com", "/path.html")
	suite.Require().NoError(err)

	comments, _, _ := suite.repo.ListComments(commenterHex, "example.com", "/path.html", false)

	assert.Equal(suite.T(), len(comments), 0)
}

func (suite *CommentRepoTestSuite) Test_DeleteComment_return_error_when_comment_hex_is_empty() {
	err := suite.repo.DeleteComment("", "test-commenter-hex", "anydomain.com", "/path.html")

	suite.Require().Error(err)
}

func (suite *CommentRepoTestSuite) TestCommentVoteBasics() {
	commenter1, _ := suite.commenterRepo.CreateCommenter("test1@example.com", "Test1", "undefined", "http://example.com/photo.jpg", "google", "")
	commenter2, _ := suite.commenterRepo.CreateCommenter("test2@example.com", "Test2", "undefined", "http://example.com/photo.jpg", "google", "")
	commenter3, _ := suite.commenterRepo.CreateCommenter("test3@example.com", "Test3", "undefined", "http://example.com/photo.jpg", "google", "")

	comment, _ := suite.repo.CreateComment(suite.NewComment(commenter1, "example.com", "/path.html", "root", "**foo**", "html", "approved", time.Now().UTC()))

	if err := suite.repo.VoteComment(commenter1, comment.CommentHex, 1, "example.com/path.html"); err != app.ErrorSelfVote {
		suite.Require().NoError(err)
	}

	if c, _, _ := suite.repo.ListComments("temp", "example.com", "/path.html", false); c[0].Score != 0 {
		suite.T().Errorf("expected c[0].Score = 0 got c[0].Score = %d", c[0].Score)
		return
	}

	if err := suite.repo.VoteComment(commenter2, comment.CommentHex, -1, "example.com/path.html"); err != nil {
		suite.T().Errorf("unexpected error voting: %v", err)
		return
	}

	if err := suite.repo.VoteComment(commenter3, comment.CommentHex, -1, "example.com/path.html"); err != nil {
		suite.T().Errorf("unexpected error voting: %v", err)
		return
	}

	if c, _, _ := suite.repo.ListComments("temp", "example.com", "/path.html", false); c[0].Score != -2 {
		suite.T().Errorf("expected c[0].Score = -2 got c[0].Score = %d", c[0].Score)
		return
	}

	if err := suite.repo.VoteComment(commenter2, comment.CommentHex, -1, "example.com/path.html"); err != nil {
		suite.T().Errorf("unexpected error voting: %v", err)
		return
	}

	if c, _, _ := suite.repo.ListComments("temp", "example.com", "/path.html", false); c[0].Score != -2 {
		suite.T().Errorf("expected c[0].Score = -2 got c[0].Score = %d", c[0].Score)
		return
	}

	if err := suite.repo.VoteComment(commenter2, comment.CommentHex, 0, "example.com/path.html"); err != nil {
		suite.T().Errorf("unexpected error voting: %v", err)
		return
	}

	if c, _, _ := suite.repo.ListComments("temp", "example.com", "/path.html", false); c[0].Score != -1 {
		suite.T().Errorf("expected c[0].Score = -1 got c[0].Score = %d", c[0].Score)
		return
	}
}

func (suite *CommentRepoTestSuite) TestCommentDomainGetEmpty() {
	_, _, err := suite.repo.GetCommentDomainPath("")

	suite.Require().Error(err)
}

func (suite *CommentRepoTestSuite) TestCommentListBasics() {
	commenterHex, _ := suite.commenterRepo.CreateCommenter("test1@example.com", "Test", "undefined", "http://example.com/photo.jpg", "google", "")

	suite.repo.CreateComment(suite.NewComment(commenterHex, "example.com", "/path.html", "root", "**foo**", "html text", "approved", time.Now().UTC()))
	suite.repo.CreateComment(suite.NewComment(commenterHex, "example.com", "/path.html", "root", "**bar**", "bar", "approved", time.Now().UTC()))

	comments, _, err := suite.repo.ListComments("temp-commenter-hex", "example.com", "/path.html", false)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), 2, len(comments))
	assert.Equal(suite.T(), 0, comments[0].Direction)
	assert.Equal(suite.T(), "html text", comments[0].Html)

	comments, _, err = suite.repo.ListComments(commenterHex, "example.com", "/path.html", false)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), len(comments), 2)
	assert.Equal(suite.T(), comments[0].Direction, 0)
}

func (suite *CommentRepoTestSuite) TestCommentListEmpty_should_return_error() {
	_, _, err := suite.repo.ListComments("temp-commenter-hex", "", "/path.html", false)

	suite.Require().Error(err)
}

func (suite *CommentRepoTestSuite) TestCommentListSelfUnapproved() {
	commentHex, _ := util.RandomHex(32)
	commenterHex, _ := suite.commenterRepo.CreateCommenter("test@example.com", "Test", "undefined", "http://example.com/photo.jpg", "google", "")
	suite.repo.CreateComment(suite.CreateComment(commentHex, commenterHex, "root", "unapproved"))

	comments, _, _ := suite.repo.ListComments("temp-commenter-hex", "example.com", "/path.html", false)
	assert.Equal(suite.T(), len(comments), 0, "expected user to not see unapproved comment")

	comments, _, _ = suite.repo.ListComments(commenterHex, "example.com", "/path.html", false)
	assert.Equal(suite.T(), len(comments), 1, "expected user to see self comments")
}

func (suite *CommentRepoTestSuite) TestCommentListAnonymousUnapproved() {
	commentHex, _ := util.RandomHex(32)
	suite.repo.CreateComment(suite.CreateComment(commentHex, "anonymous", "root", "unapproved"))

	comments, _, _ := suite.repo.ListComments("anonymous", "example.com", "/path.html", false)

	assert.Equal(suite.T(), len(comments), 0, "expected user to not see unapproved anonymous comment as anonymous")
}

func (suite *CommentRepoTestSuite) TestCommentListIncludeUnapproved() {
	commentHex, _ := util.RandomHex(32)
	suite.repo.CreateComment(suite.CreateComment(commentHex, "anonymous", "root", "unapproved"))

	comments, _, _ := suite.repo.ListComments("anonymous", "example.com", "/path.html", true)

	assert.Equal(suite.T(), len(comments), 1, "expected to see unapproved comments when includeUnapproved is true")
}

func (suite *CommentRepoTestSuite) TestCommentListDifferentPaths() {
	suite.repo.CreateComment(suite.NewComment("anonymous", "example.com", "/path1.html", "root", "**foo**", "html", "unapproved", time.Now().UTC()))
	suite.repo.CreateComment(suite.NewComment("anonymous", "example.com", "/path1.html", "root", "**foo**", "html", "unapproved", time.Now().UTC()))
	suite.repo.CreateComment(suite.NewComment("anonymous", "example.com", "/path2.html", "root", "**foo**", "html", "unapproved", time.Now().UTC()))

	comments, _, _ := suite.repo.ListComments("anonymous", "example.com", "/path1.html", true)
	assert.Equal(suite.T(), len(comments), 2)

	comments, _, _ = suite.repo.ListComments("anonymous", "example.com", "/path2.html", true)
	assert.Equal(suite.T(), len(comments), 1)
}

func (suite *CommentRepoTestSuite) TestCommentListDifferentDomains() {
	suite.repo.CreateComment(suite.NewComment("anonymous", "example1.com", "/path.html", "root", "**foo**", "html", "unapproved", time.Now().UTC()))
	suite.repo.CreateComment(suite.NewComment("anonymous", "example2.com", "/path.html", "root", "**foo**", "html", "unapproved", time.Now().UTC()))

	comments, _, _ := suite.repo.ListComments("anonymous", "example1.com", "/path.html", true)
	assert.Equal(suite.T(), len(comments), 1)

	comments, _, _ = suite.repo.ListComments("anonymous", "example2.com", "/path.html", true)
	assert.Equal(suite.T(), len(comments), 1)
}

func (suite *CommentRepoTestSuite) TestCommentNewBasics() {
	_, err := suite.repo.CreateComment(suite.CreateComment("temp-commenter-hex", "anonymous", "root", "approved"))

	suite.Require().NoError(err)
}
