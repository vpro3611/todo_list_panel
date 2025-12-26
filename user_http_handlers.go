package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func (s *Server) GetAllUsersHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := s.userSvc.GetAllUsers(ctx)
	if err != nil {
		log.Println("Error getting all users: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = EncodeJSONhelper(w, users)
	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetUserByIdHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	idInt, err := ConvertToInt(id)
	if err != nil {
		log.Println("Error parsing id: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := s.userSvc.GetUserById(ctx, idInt)
	if err != nil {
		log.Println("Error getting user by id: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = EncodeJSONhelper(w, user)
	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) CreateNewUserHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("Error decoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := s.userSvc.CreateNewUser(ctx, user.Name, user.Password)
	if err != nil {
		log.Println("Error creating new user: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusCreated)

	response := map[string]any{
		"id":     id,
		"status": "User successfully created",
	}

	err = EncodeJSONhelper(w, response)
	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) RenameUserHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	idInt, err := ConvertToInt(id)

	if err != nil {
		log.Println("Error parsing id: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var input struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println("Error decoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.userSvc.RenameUser(ctx, idInt, input.Name)
	if err != nil {
		log.Println("Error renaming user: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]any{
		"id":     idInt,
		"status": "User successfully renamed",
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = EncodeJSONhelper(w, response)
	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) ChangeUserPasswordHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	idInt, err := ConvertToInt(id)
	if err != nil {
		log.Println("Error parsing id: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var inputPasswords struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&inputPasswords); err != nil {
		log.Println("Error decoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.userSvc.ChangeUsersPass(ctx, idInt, inputPasswords.OldPassword, inputPasswords.NewPassword)

	if err != nil {
		log.Println("Error changing user password: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	response := map[string]any{
		"id":     idInt,
		"status": "Password successfully changed",
	}
	err = EncodeJSONhelper(w, response)
	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) DeleteUserHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	idInt, err := ConvertToInt(id)

	if err != nil {
		log.Println("Error parsing id: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.userSvc.DeleteUser(ctx, idInt)

	if err != nil {
		log.Println("Error deleting user: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := map[string]any{
		"id":     idInt,
		"status": "User successfully deleted",
	}
	err = EncodeJSONhelper(w, response)
	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
