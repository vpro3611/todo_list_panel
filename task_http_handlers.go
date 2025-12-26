package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

// this is all for tasks, in this case admin can view all tasks, change their status, etc.
// in the users section you can also do the same
// this is only for tasks, not for users and only for admins

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

func (s *Server) GetTaskByUserIdHTTP(w http.ResponseWriter, r *http.Request) {
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

// here we end the task http handlers for admins
// and immediately start the user (current client) handlers for tasks

func (s *Server) GetMyTasksHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(userContextKey).(*Claims)
	if !ok {
		log.Println("Error getting user id from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id := claims.UserID

	tasks, err := s.taskSvc.GetTaskById(ctx, id)
	if err != nil {
		log.Println("Error getting tasks by id: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = EncodeJSONhelper(w, tasks)
	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) CreateMyTaskHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(userContextKey).(*Claims)
	if !ok {
		log.Println("Error getting user id from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userId := claims.UserID

	var inputTask struct {
		UserId      int    `json:"user_id,omitempty"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&inputTask); err != nil {
		log.Println("Error decoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	idTaskCreate, err := s.taskSvc.CreateNewTask(ctx, userId, inputTask.Title, inputTask.Description)
	if err != nil {
		log.Println("Error creating new task: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := map[string]any{
		"id":     idTaskCreate,
		"status": "Task successfully created",
	}
	err = EncodeJSONhelper(w, response)
	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
