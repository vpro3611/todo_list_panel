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

func (s *Server) GetAllTasksHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks, err := s.taskSvc.GetAllTasks(ctx)
	if err != nil {
		log.Println("Error getting all tasks: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = EncodeJSONhelper(w, tasks)
	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetTaskByIdHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	idInt, err := ConvertToInt(id)

	if err != nil {
		log.Println("Error parsing id: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := s.taskSvc.GetTaskById(ctx, idInt)
	if err != nil {
		log.Println("Error getting task by id: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = EncodeJSONhelper(w, task)

	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) CreateNewTaskHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := chi.URLParam(r, "id")
	userIdInt, err := ConvertToInt(userId)
	if err != nil {
		log.Println("Error parsing id: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var task Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Println("Error decoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	taskId, err := s.taskSvc.CreateNewTask(ctx, userIdInt, task.Title, task.Description)
	if err != nil {
		log.Println("Error creating new task: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	var response map[string]any = map[string]any{
		"id":     taskId,
		"status": "Task successfully created",
	}

	err = EncodeJSONhelper(w, response)
	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) DeleteTaskHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	idInt, err := ConvertToInt(id)
	if err != nil {
		log.Println("Error parsing id: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.taskSvc.DeleteTask(ctx, idInt)
	if err != nil {
		log.Println("Error deleting task: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := map[string]any{
		"id":     idInt,
		"status": "Task successfully deleted",
	}
	err = EncodeJSONhelper(w, response)
	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) UpdateTaskTitleHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	idInt, err := ConvertToInt(id)
	if err != nil {
		log.Println("Error parsing id: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var taskTitleForUpdate struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&taskTitleForUpdate); err != nil {
		log.Println("Error decoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.taskSvc.UpdateTitle(ctx, taskTitleForUpdate.Title, idInt)
	if err != nil {
		log.Println("Error updating task title: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]any{
		"id":     idInt,
		"status": "Task's title successfully updated",
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = EncodeJSONhelper(w, response)
	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) UpdateTaskDescriptionHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	idInt, err := ConvertToInt(id)

	if err != nil {
		log.Println("Error parsing id: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var taskDescriptionForUpdate struct {
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&taskDescriptionForUpdate); err != nil {
		log.Println("Error decoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.taskSvc.UpdateDescription(ctx, taskDescriptionForUpdate.Description, idInt)

	if err != nil {
		log.Println("Error updating task description: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]any{
		"id":     idInt,
		"status": "Task's description successfully updated",
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = EncodeJSONhelper(w, response)
	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) SwitchTaskStatusHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	idInt, err := ConvertToInt(id)

	if err != nil {
		log.Println("Error parsing id: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.taskSvc.SwitchTaskStatus(ctx, idInt)

	if err != nil {
		log.Println("Error switching task status: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]any{
		"id":     idInt,
		"status": "Task's status successfully switched",
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = EncodeJSONhelper(w, response)
	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
