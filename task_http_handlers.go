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
	claims, ok := ctx.Value(userContextKey).(*Claims)
	if !ok || claims.Role != ADMIN {
		log.Printf("This is for admins only! Unsafe request for GetAllTasks: %s with id %d\n", claims.Role, claims.UserID)
		http.Error(w, "This is for admins only!", http.StatusForbidden)
		return
	}
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
	claims, ok := ctx.Value(userContextKey).(*Claims)
	if !ok {
		log.Println("Error getting user id from context")
		http.Error(w, "--Unauthorized", http.StatusUnauthorized)
		return
	}

	targetId, ok := ctx.Value(targetIdContextKey).(int)
	if !ok {
		log.Println("Error getting target user id from context")
		http.Error(w, "--Unauthorized", http.StatusInternalServerError)
		return
	}

	task, err := s.taskSvc.GetTaskById(ctx, targetId, claims.UserID, claims.Role)

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
	claims, ok := ctx.Value(userContextKey).(*Claims)
	if !ok {
		log.Println("Error getting user id from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	targetUserId := chi.URLParam(r, "id")

	var finalUserId int

	if targetUserId == "" {
		finalUserId = claims.UserID
	} else {
		idInt, err := ConvertToInt(targetUserId)
		if err != nil {
			log.Println("Error parsing id: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if claims.Role == ADMIN {
			finalUserId = idInt
		} else {
			finalUserId = claims.UserID
		}
	}

	var task Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Println("Error decoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	taskId, err := s.taskSvc.CreateNewTask(ctx, finalUserId, task.Title, task.Description)
	if err != nil {
		log.Println("Error creating new task: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	taskGotten, err := s.taskSvc.GetTaskByItsId(ctx, taskId, claims.UserID, claims.Role)
	if err != nil {
		log.Println("Error getting task by id: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = EncodeJSONhelper(w, taskGotten)
	if err != nil {
		log.Println("Error encoding JSON: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) DeleteTaskHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(userContextKey).(*Claims)
	if !ok {
		log.Println("Error getting user id from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	idInt, err := ConvertToInt(id)
	if err != nil {
		log.Println("Error parsing id: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.taskSvc.DeleteTask(ctx, idInt, claims.UserID, claims.Role)
	if err != nil {
		log.Println("Error deleting task: ", err)
		http.Error(w, err.Error(), http.StatusForbidden)
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
	claims, ok := ctx.Value(userContextKey).(*Claims)
	if !ok {
		log.Println("Error getting user id from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

	defer r.Body.Close()

	err = s.taskSvc.UpdateTitle(ctx, taskTitleForUpdate.Title, idInt, claims.UserID, claims.Role)
	if err != nil {
		log.Println("Error updating task title: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := s.taskSvc.GetTaskByItsId(ctx, idInt, claims.UserID, claims.Role)
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
		return
	}
}

func (s *Server) UpdateTaskDescriptionHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := ctx.Value(userContextKey).(*Claims)
	if !ok {
		log.Println("Error getting user id from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

	defer r.Body.Close()

	err = s.taskSvc.UpdateDescription(ctx, taskDescriptionForUpdate.Description, idInt, claims.UserID, claims.Role)

	if err != nil {
		log.Println("Error updating task description: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := s.taskSvc.GetTaskByItsId(ctx, idInt, claims.UserID, claims.Role)
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
		return
	}
}

func (s *Server) SwitchTaskStatusHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := ctx.Value(userContextKey).(*Claims)
	if !ok {
		log.Println("Error getting user id from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	idInt, err := ConvertToInt(id)

	if err != nil {
		log.Println("Error parsing id: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.taskSvc.SwitchTaskStatus(ctx, idInt, claims.UserID, claims.Role)

	if err != nil {
		log.Println("Error switching task status: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := s.taskSvc.GetTaskByItsId(ctx, idInt, claims.UserID, claims.Role)
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
		return
	}
}

// here we end the task http handlers for admins
