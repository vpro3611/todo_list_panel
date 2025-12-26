package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	userSvc *UserService
	taskSvc *TaskService
	router  *chi.Mux
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
	s.router.Route("/users", func(r chi.Router) {
		r.Get("/", s.GetAllUsersHTTP)
		r.Get("/{id}", s.GetUserByIdHTTP)
		r.Post("/", s.CreateNewUserHTTP)
		r.Patch("/{id}/rename", s.RenameUserHTTP)
		r.Patch("/{id}/password", s.ChangeUserPasswordHTTP)
		r.Delete("/{id}", s.DeleteUserHTTP)
		r.Post("/{id}/tasks", s.CreateNewTaskHTTP) // создать таск данному пользователю
		r.Get("/{id}/tasks", s.GetTaskByIdHTTP)    // получить таски данного пользователя

	})
	s.router.Route("/tasks", func(r chi.Router) {
		r.Get("/", s.GetAllTasksHTTP)
		r.Delete("/{id}", s.DeleteTaskHTTP)
		r.Patch("/{id}/title", s.UpdateTaskTitleHTTP)
		r.Patch("/{id}/description", s.UpdateTaskDescriptionHTTP)
		r.Patch("/{id}/switch", s.SwitchTaskStatusHTTP)
	})
}
