package handler

import (
	"clientManage/internal/domain/model"
	"clientManage/internal/service"
	"clientManage/internal/transport/http/auth"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ClientHandler struct {
	service *service.ClientService
}

func NewClientHandler(service *service.ClientService) *ClientHandler {
	return &ClientHandler{service: service}
}

func getUserRole(r *http.Request) string {
	return r.Header.Get("user_role")
}

func (h *ClientHandler) RegisterClient(w http.ResponseWriter, r *http.Request) {
	var client model.Client
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = client.SetPassword(client.PasswordHash)
	if err != nil {
		http.Error(w, "Failed to set password", http.StatusInternalServerError)
		return
	}

	err = h.service.CreateClient(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(client)
}

func (h *ClientHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client, err := h.service.GetClientByEmail(credentials.Email)
	if err != nil {
		http.Error(w, "Error retrieving client", http.StatusInternalServerError)
		return
	}
	if client == nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = client.CheckPassword(credentials.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(client.Email, client.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	role := getUserRole(r)
	if role != "client" && role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var client model.Client
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.CreateClient(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(client)
}

func (h *ClientHandler) GetClient(w http.ResponseWriter, r *http.Request) {
	role := getUserRole(r)
	if role != "client" && role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["client_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	client, err := h.service.GetClientByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if client == nil {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(client)
}

func (h *ClientHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	role := getUserRole(r)
	if role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["client_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	var client model.Client
	err = json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client.ID = id

	err = h.service.UpdateClient(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(client)
}

func (h *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	role := getUserRole(r)
	if role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["client_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteClient(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ClientHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	role := getUserRole(r)
	if role != "client" && role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	query := r.URL.Query()
	offset, _ := strconv.Atoi(query.Get("offset"))
	limit, _ := strconv.Atoi(query.Get("limit"))
	filters := make(map[string]interface{})
	for key, values := range query {
		if key != "offset" && key != "limit" && key != "sort_by" && key != "sort_order" {
			filters[key] = values[0]
		}
	}
	sortBy := query.Get("sort_by")
	if sortBy == "" {
		sortBy = "id"
	}
	sortOrder := query.Get("sort_order")
	if sortOrder == "" {
		sortOrder = "asc"
	}

	clients, err := h.service.ListClients(offset, limit, filters, sortBy, sortOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(clients)
}
