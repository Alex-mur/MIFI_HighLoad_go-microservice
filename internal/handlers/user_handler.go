package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go-microservice/internal/models"
	"go-microservice/internal/services"
)

type UserHandler struct {
	userService  *services.UserService
	auditService *services.AuditService
}

func NewUserHandler(userService *services.UserService, auditService *services.AuditService) *UserHandler {
	return &UserHandler{
		userService:  userService,
		auditService: auditService,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAll()
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		h.auditService.LogErrorAsync("GET_USERS", 0, err)
		return
	}

	h.sendJSON(w, http.StatusOK, users)
	h.auditService.LogAsync("GET_USERS", 0, "Retrieved all users")
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid user ID")
		h.auditService.LogErrorAsync("GET_USER", id, err)
		return
	}

	user, err := h.userService.GetByID(id)
	if err != nil {
		if err == services.ErrUserNotFound {
			h.sendError(w, http.StatusNotFound, err.Error())
		} else {
			h.sendError(w, http.StatusInternalServerError, err.Error())
		}
		h.auditService.LogErrorAsync("GET_USER", id, err)
		return
	}

	h.sendJSON(w, http.StatusOK, user)
	h.auditService.LogAsync("GET_USER", id, "Retrieved user")
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		h.auditService.LogErrorAsync("CREATE_USER", 0, err)
		return
	}

	user, err := h.userService.Create(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrInvalidInput {
			status = http.StatusBadRequest
		}
		h.sendError(w, status, err.Error())
		h.auditService.LogErrorAsync("CREATE_USER", 0, err)
		return
	}

	h.sendJSON(w, http.StatusCreated, user)

	// Асинхронное логирование и отправка уведомлений
	data, _ := json.Marshal(user)
	h.auditService.LogAsync("CREATE_USER", user.ID, string(data))
	go h.sendNotification("User created", user.ID)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid user ID")
		h.auditService.LogErrorAsync("UPDATE_USER", id, err)
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		h.auditService.LogErrorAsync("UPDATE_USER", id, err)
		return
	}

	user, err := h.userService.Update(id, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrUserNotFound {
			status = http.StatusNotFound
		} else if err == services.ErrInvalidInput {
			status = http.StatusBadRequest
		}
		h.sendError(w, status, err.Error())
		h.auditService.LogErrorAsync("UPDATE_USER", id, err)
		return
	}

	h.sendJSON(w, http.StatusOK, user)
	h.auditService.LogAsync("UPDATE_USER", id, "User updated")
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid user ID")
		h.auditService.LogErrorAsync("DELETE_USER", id, err)
		return
	}

	if err := h.userService.Delete(id); err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrUserNotFound {
			status = http.StatusNotFound
		}
		h.sendError(w, status, err.Error())
		h.auditService.LogErrorAsync("DELETE_USER", id, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.auditService.LogAsync("DELETE_USER", id, "User deleted")
}

func (h *UserHandler) sendNotification(message string, userID int) {
	// В реальном приложении здесь была бы отправка email, SMS,
	// сообщение в мессенджер и т.д.
	// Это выполняется асинхронно в отдельной горутине
	h.auditService.LogAsync("NOTIFICATION", userID, message)
}

func (h *UserHandler) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *UserHandler) sendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
