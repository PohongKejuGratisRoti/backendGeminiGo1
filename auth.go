package main

import (
	"net/http"
)

func checkSession(req *http.Request) (int, bool) {
	sessionContext := req.Context()

	if !sessionManager.Exists(sessionContext, "userID") {
		return 0, false
	}

	UserID := sessionManager.GetInt(sessionContext, "userID")
	return UserID, true
}

