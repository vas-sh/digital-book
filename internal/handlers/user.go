package handlers

import (
	"digital-book/internal/types"
	"log"
	"net/http"
)

func (s *server) GetUsers(rw http.ResponseWriter, r *http.Request) {
	users, err := s.srv.GetUsers(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	data := struct {
		Users []types.User
	}{
		Users: users,
	}
	s.renderTemplate("html/users.html", rw, data)
}

func (s *server) CreateUser(rw http.ResponseWriter, r *http.Request) {
	log.Println("createUser", r.Method)
	ctx := r.Context()

	switch r.Method {
	case http.MethodPost:
		class := r.FormValue("class")
		login := r.FormValue("login")
		name := r.FormValue("name")
		id := r.FormValue("id")

		if id == "" {
			log.Println("new user: name", name, "class", class)
			if err := s.srv.CreateUser(ctx, &types.User{Name: name, Class: class, Login: login}); err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}
		} else {
			log.Println("update user: name", name, "class", class, "id", id)
			if err := s.srv.UpdateUser(ctx, name, class, id, login); err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}
		}

		http.Redirect(rw, r, "/users", http.StatusTemporaryRedirect)
		return

	case http.MethodGet:
		if id := r.URL.Query().Get("id"); id != "" {
			user, err := s.srv.GetUser(ctx, id)
			if err != nil {
				http.Error(rw, err.Error(), 400)
				return
			}
			s.renderTemplate("html/update-user.html", rw, struct {
				User types.User
			}{
				User: user,
			})
		} else {
			s.renderTemplate("html/create-user.html", rw, nil)
		}
	}
}

func (s *server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodGet {
		http.Error(w, "not supported", http.StatusNotImplemented)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	if err := s.srv.DeleteUser(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/users", http.StatusTemporaryRedirect)
}
