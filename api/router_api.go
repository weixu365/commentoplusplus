package main

import (
	"net/http"
	"simple-commenting/handler"
	"simple-commenting/notification"

	"github.com/gorilla/mux"
)

func apiRouterInit(router *mux.Router) error {
	router.HandleFunc("/api/owner/new", handler.OwnerNewHandler).Methods("POST")
	router.HandleFunc("/api/owner/confirm-hex", handler.OwnerConfirmHexHandler).Methods("GET")
	router.HandleFunc("/api/owner/login", handler.OwnerLoginHandler).Methods("POST")
	router.HandleFunc("/api/owner/self", handler.OwnerSelfHandler).Methods("POST")
	router.HandleFunc("/api/owner/delete", handler.OwnerDeleteHandler).Methods("POST")

	router.HandleFunc("/api/domain/new", handler.DomainNewHandler).Methods("POST")
	router.HandleFunc("/api/domain/delete", handler.DomainDeleteHandler).Methods("POST")
	router.HandleFunc("/api/domain/clear", handler.DomainClearHandler).Methods("POST")
	router.HandleFunc("/api/domain/sso/new", handler.DomainSsoSecretNewHandler).Methods("POST")
	router.HandleFunc("/api/domain/list", handler.DomainListHandler).Methods("POST")
	router.HandleFunc("/api/domain/update", handler.DomainUpdateHandler).Methods("POST")
	router.HandleFunc("/api/domain/moderator/new", handler.DomainModeratorNewHandler).Methods("POST")
	router.HandleFunc("/api/domain/moderator/delete", handler.DomainModeratorDeleteHandler).Methods("POST")
	router.HandleFunc("/api/domain/statistics", handler.DomainStatisticsHandler).Methods("POST")
	router.HandleFunc("/api/domain/import/disqus", handler.DomainImportDisqusHandler).Methods("POST")
	router.HandleFunc("/api/domain/import/commento", handler.DomainImportCommentoHandler).Methods("POST")
	router.HandleFunc("/api/domain/export/begin", handler.DomainExportBeginHandler).Methods("POST")
	router.HandleFunc("/api/domain/export/download", handler.DomainExportDownloadHandler).Methods("GET")

	router.HandleFunc("/api/commenter/token/new", handler.CommenterTokenNewHandler).Methods("GET")
	router.HandleFunc("/api/commenter/new", handler.CommenterNewHandler).Methods("POST")
	router.HandleFunc("/api/commenter/login", handler.CommenterLoginHandler).Methods("POST")
	router.HandleFunc("/api/commenter/self", handler.CommenterSelfHandler).Methods("POST")
	router.HandleFunc("/api/commenter/update", handler.CommenterUpdateHandler).Methods("POST")
	router.HandleFunc("/api/commenter/photo", handler.CommenterPhotoHandler).Methods("GET")
	router.HandleFunc("/api/commenter/delete", handler.CommenterDeleteHandler).Methods("POST")

	router.HandleFunc("/api/forgot", handler.ForgotHandler).Methods("POST")
	router.HandleFunc("/api/reset", handler.ResetHandler).Methods("POST")

	router.HandleFunc("/api/email/get", handler.EmailGetHandler).Methods("POST")
	router.HandleFunc("/api/email/update", handler.EmailUpdateHandler).Methods("POST")
	router.HandleFunc("/api/email/moderate", handler.EmailModerateHandler).Methods("GET")

	router.HandleFunc("/api/oauth/google/redirect", handler.GoogleRedirectHandler).Methods("GET")
	router.HandleFunc("/api/oauth/google/callback", handler.GoogleCallbackHandler).Methods("GET")

	router.HandleFunc("/api/oauth/github/redirect", handler.GithubRedirectHandler).Methods("GET")
	router.HandleFunc("/api/oauth/github/callback", handler.GithubCallbackHandler).Methods("GET")

	router.HandleFunc("/api/oauth/twitter/redirect", handler.TwitterRedirectHandler).Methods("GET")
	router.HandleFunc("/api/oauth/twitter/callback", handler.TwitterCallbackHandler).Methods("GET")

	router.HandleFunc("/api/oauth/gitlab/redirect", handler.GitlabRedirectHandler).Methods("GET")
	router.HandleFunc("/api/oauth/gitlab/callback", handler.GitlabCallbackHandler).Methods("GET")

	router.HandleFunc("/api/oauth/sso/redirect", handler.SsoRedirectHandler).Methods("GET")
	router.HandleFunc("/api/oauth/sso/callback", handler.SsoCallbackHandler).Methods("GET")

	router.HandleFunc("/api/comment/new", handler.CommentNewHandler).Methods("POST")
	router.HandleFunc("/api/comment/edit", handler.CommentEditHandler).Methods("POST")
	router.HandleFunc("/api/comment/list", handler.CommentListHandler).Methods("POST")
	router.HandleFunc("/api/comment/count", handler.CommentCountHandler).Methods("POST")
	router.HandleFunc("/api/comment/vote", handler.CommentVoteHandler).Methods("POST")
	router.HandleFunc("/api/comment/approve", handler.CommentApproveHandler).Methods("POST")
	router.HandleFunc("/api/comment/delete", handler.CommentDeleteHandler).Methods("POST")
	router.HandleFunc("/api/comment/owner/list", handler.CommentListApprovalsHandler).Methods("POST")
	router.HandleFunc("/api/comment/owner/listAll", handler.CommentListAllHandler).Methods("POST")
	router.HandleFunc("/api/comment/owner/approve", handler.CommentOwnerApproveHandler).Methods("POST")
	router.HandleFunc("/api/comment/owner/delete", handler.CommentOwnerDeleteHandler).Methods("POST")

	router.HandleFunc("/api/page/update", handler.PageUpdateHandler).Methods("POST")

	notification.NotificationHub = notification.NewHub()
	go notification.NotificationHub.Run()
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		notification.ServeWs(notification.NotificationHub, w, r)
	})

	return nil
}
