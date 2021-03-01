package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type User struct {
	Username string    `json:"username"`
	Password string    `json:"password"`
	ID       uuid.UUID `json:"id"`
	Age      int       `json:"age"`
}

type Server struct {
	*mux.Router

	loggedUsers []User
}

func NewServer() *Server {
	s := &Server{
		Router:      mux.NewRouter(),
		loggedUsers: []User{},
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.HandleFunc("/logged-users", s.listUsers()).Methods("GET")
	s.HandleFunc("/logged-users", s.createUser()).Methods("POST")
	s.HandleFunc("/logged-users/{id}", s.removeUser()).Methods("DELETE")
}

func (s *Server) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u.ID = uuid.New()
		s.loggedUsers = append(s.loggedUsers, u)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) listUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.loggedUsers); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func (s Server) removeUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, _ := mux.Vars(r)["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		for i, user := range s.loggedUsers {
			if user.ID == id {
				s.loggedUsers = append(s.loggedUsers[:i], s.loggedUsers[i+1:]...)
				break
			}
		}
	}
}
