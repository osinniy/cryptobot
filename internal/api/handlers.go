package api

import (
	"net/http"
	"osinniy/cryptobot/internal/store"
	"strconv"
)

// Endpoint: /
func (s Server) root(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("ok"))
	s.handleWrErr(err, &w)
}

// Endpoint: /users/total
//
// Returns total numbers of active usersTotal
func (s *Server) usersTotal(w http.ResponseWriter, _ *http.Request) {
	totalUsers, err := s.store.Users().UsersTotal()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		_, err = w.Write([]byte(strconv.Itoa(totalUsers)))
		s.handleWrErr(err, &w)
	}
}

// Endpoint: /users/lang
func (s *Server) usersLang(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userIdQ := r.Form.Get("user")
	if userIdQ == "" {
		s.e(w, errMissingUserId)
		return
	}

	userId, err := strconv.ParseInt(userIdQ, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lang, err := s.store.Users().Lang(userId)
	if err != nil {
		if err == store.ErrUserNotFound {
			s.e(w, errUserNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(lang))
	s.handleWrErr(err, &w)
}
