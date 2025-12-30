package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// this is all for admin, i.e., you can view all users, change their roles, etc.
// in the tasks section you can also do the same
// this is only for users, not for tasks and only for admins

func (s *Server) GetAllUsersHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(userContextKey).(*Claims)
	if !ok || claims.Role != ADMIN {
		log.Printf("This is for admins only! Unsafe request for GetAllUsers: %s with id %d\n", claims.Role, claims.UserID)
		http.Error(w, "This is for admins only!", http.StatusForbidden)
		return
	}
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

	claims, ok := ctx.Value(userContextKey).(*Claims)
	if !ok {
		log.Println("Error getting user id from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	targetUserId, ok := ctx.Value(targetIdContextKey).(int)
	if !ok {
		log.Println("Error getting target user id from context")
		http.Error(w, "Unauthorized", http.StatusInternalServerError)
		return
	}

	user, err := s.userSvc.GetUserById(ctx, targetUserId, claims.UserID, claims.Role)
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
	defer r.Body.Close()

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
	claims, ok := ctx.Value(userContextKey).(*Claims)
	if !ok {
		log.Println("Error getting user id from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	targetUserId, ok := ctx.Value(targetIdContextKey).(int)
	if !ok {
		log.Println("Error getting target user id from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

	defer r.Body.Close()

	err := s.userSvc.RenameUser(ctx, targetUserId, input.Name, claims.UserID, claims.Role)
	if err != nil {
		log.Println("Error renaming user: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.userSvc.GetUserById(ctx, targetUserId, claims.UserID, claims.Role)
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

func (s *Server) ChangeUserPasswordHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := ctx.Value(userContextKey).(*Claims)
	if !ok {
		log.Println("Error getting user id from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	targetUserId, ok := ctx.Value(targetIdContextKey).(int)
	if !ok {
		log.Println("Error getting target user id from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

	defer r.Body.Close()

	err := s.userSvc.ChangeUsersPass(ctx, targetUserId, inputPasswords.OldPassword, inputPasswords.NewPassword, claims.UserID, claims.Role)

	if err != nil {
		log.Println("Error changing user password: ", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	user, err := s.userSvc.GetUserById(ctx, targetUserId, claims.UserID, claims.Role)
	if err != nil {
		log.Println("Error getting user by id: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	err = EncodeJSONhelper(w, user)
	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) DeleteUserHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := ctx.Value(userContextKey).(*Claims)
	if !ok {
		log.Println("Error getting user id from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	targetUserId, ok := ctx.Value(targetIdContextKey).(int)
	if !ok {
		log.Println("Error getting target user id from context")
		http.Error(w, "Unauthorized", http.StatusInternalServerError)
		return
	}

	err := s.userSvc.DeleteUser(ctx, targetUserId, claims.UserID, claims.Role)

	if err != nil {
		log.Println("Error deleting user: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := map[string]any{
		"id":     targetUserId,
		"status": "User successfully deleted",
	}
	err = EncodeJSONhelper(w, response)
	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (s *Server) UpdateRoleHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(userContextKey).(*Claims)
	if !ok || claims.Role != ADMIN {
		log.Printf("This is for admins only! Unsafe request for UpdateRole: %s with id %d\n", claims.Role, claims.UserID)
		http.Error(w, "This is for admins only!", http.StatusForbidden)
		return
	}

	targetId, ok := ctx.Value(targetIdContextKey).(int)
	if !ok {
		log.Println("Error getting target user id from context")
		http.Error(w, "Unauthorized", http.StatusInternalServerError)
		return
	}

	err := s.userSvc.UpdateUserRole(ctx, targetId, claims.UserID, claims.Role)
	if err != nil {
		log.Println("Error updating user role: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.userSvc.GetUserById(ctx, targetId, claims.UserID, claims.Role)
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

// end of admin handlers
