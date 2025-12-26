package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

type Server struct {
	userSvc *UserService
	taskSvc *TaskService
	router  *chi.Mux
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (s *Server) LoginHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Error decoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.userSvc.AuthenticateUser(ctx, req.Name, req.Password)
	if err != nil {
		log.Println("Error authenticating user: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := GenerateJWT(user.Id, user.Role)
	if err != nil {
		log.Println("Error generating JWT: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := LoginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = EncodeJSONhelper(w, response)
	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewServer(userSvc *UserService, taskSvc *TaskService) *Server {
	s := &Server{
		userSvc: userSvc,
		taskSvc: taskSvc,
		router:  chi.NewRouter(),
	}
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.StripSlashes)

	s.Routes()
	return s
}

func (s *Server) Routes() {

	//s.router.Post("/setup-admin", s.CreateNewUserHTTP)

	s.router.Post("/login", s.LoginHTTP)

	s.router.Group(func(r chi.Router) {
		r.Use(JWTmiddleware)
		r.Route("/admin", func(r chi.Router) {
			r.Use(AdminOnly)
			// admin can see all users and do these actions with them
			r.Route("/users", func(r chi.Router) {
				r.Get("/", s.GetAllUsersHTTP)
				r.Get("/{id}", s.GetUserByIdHTTP)
				r.Post("/", s.CreateNewUserHTTP)
				r.Patch("/{id}/rename", s.RenameUserHTTP)
				r.Patch("/{id}/password", s.ChangeUserPasswordHTTP)
				r.Delete("/{id}", s.DeleteUserHTTP)
				r.Post("/{id}/tasks", s.CreateNewTaskHTTP)  // создать таск данному пользователю
				r.Get("/{id}/tasks", s.GetTaskByUserIdHTTP) // получить таски данного пользователя
				r.Patch("/{id}/role", s.UpdateRoleHTTP)
			})
			// admin can see all tasks and do these actions with them, as well as with users
			r.Route("/tasks", func(r chi.Router) {
				r.Get("/", s.GetAllTasksHTTP)
				r.Delete("/{id}", s.DeleteTaskHTTP)
				r.Patch("/{id}/title", s.UpdateTaskTitleHTTP)
				r.Patch("/{id}/description", s.UpdateTaskDescriptionHTTP)
				r.Patch("/{id}/switch", s.SwitchTaskStatusHTTP)
			})
		})
		// this is for the users, /me means that they are logged in and can do actions with their own account ONLY
		r.Route("/me", func(r chi.Router) {
			r.Get("/", s.GetMyProfileHTTP)
			r.Get("/tasks", s.GetMyTasksHTTP)
			r.Post("/tasks", s.CreateMyTaskHTTP)
			// TODO : add :
			//  1. switch status task of a current user
			//  2. delete task of a current user
			//  3. update task title of a current user
			//  4. update task description of a current user
			//  5. update current user's password
			//  6. update current user's name
			//  7. <Optional> delete current user's account, but this is not compulsory
		})
	})

}
