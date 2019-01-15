package server

import (
	"net/http"

	log "github.com/duo-labs/webauthn.io/logger"
	"github.com/duo-labs/webauthn.io/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// CreateUser adds a new user to the database
func (ws *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	if username == "" {
		jsonResponse(w, "No username specified", http.StatusBadRequest)
		return
	}
	if email == "" {
		jsonResponse(w, "No email specified", http.StatusBadRequest)
		return
	}
	user, err := models.GetUserByUsername(email)
	if err != gorm.ErrRecordNotFound {
		log.Errorf("user already exists: %s", email)
		jsonResponse(w, user, http.StatusOK)
		return
	}
	u := models.User{
		Username:    email,
		DisplayName: username,
		Icon:        models.PlaceholderUserIcon,
	}
	err = models.PutUser(&u)
	if err != nil {
		jsonResponse(w, "Error Creating User", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, u, http.StatusCreated)
}

// GetUser retrieves a user from the database
func (ws *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["name"]
	u, err := models.GetUserByUsername(username)
	if err != nil {
		log.Errorf("user not found: %s: %s", username, err)
		jsonResponse(w, "User not found, try registering one first!", http.StatusNotFound)
		return
	}
	jsonResponse(w, u, http.StatusOK)
}
