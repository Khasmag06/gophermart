package handlers

import (
	"errors"
	"github.com/Khasmag06/gophermart/internal/middlewares"
	"github.com/Khasmag06/gophermart/internal/repository"
	"github.com/theplant/luhn"
	"io"
	"net/http"
	"strconv"
)

func (s *Service) NewOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(int)
	order, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderNum, err := strconv.Atoi(string(order))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if valid := luhn.Valid(orderNum); !valid {
		http.Error(w, repository.ErrOrderInvalidNum.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := s.repo.AddOrder(userID, string(order)); err != nil {
		if errors.Is(err, repository.ErrOrderUploadedByUser) {
			http.Error(w, err.Error(), http.StatusOK)
			return
		}
		if errors.Is(err, repository.ErrOrderUploadedByOtherUser) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
