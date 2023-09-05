package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Khasmag06/gophermart/internal/helpers"
	"github.com/Khasmag06/gophermart/internal/models"
	"github.com/Khasmag06/gophermart/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func (s *Service) Login(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	password := user.Password

	if err := s.repo.Login(user); err != nil {
		if errors.Is(err, repository.ErrUserCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		http.Error(w, repository.ErrUserCredentials.Error(), http.StatusUnauthorized)
		return
	}

	token, err := helpers.GenerateJWT(strconv.Itoa(user.ID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "jwt",
		Value: token})
}
