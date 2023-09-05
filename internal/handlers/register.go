package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Khasmag06/gophermart/internal/helpers"
	"github.com/Khasmag06/gophermart/internal/models"
	"github.com/Khasmag06/gophermart/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
)

func (s *Service) Register(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = string(passwordHash)

	if err := s.repo.AddUser(user); err != nil {
		if errors.Is(err, repository.ErrExistsData) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		} else {
			http.Error(w, "internal server error:", http.StatusInternalServerError)
			return
		}
	}
	token, err := helpers.GenerateJWT(strconv.Itoa(user.ID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "jwt",
		Value: token,
	})

}
