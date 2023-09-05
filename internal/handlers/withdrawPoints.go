package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Khasmag06/gophermart/internal/middlewares"
	"github.com/Khasmag06/gophermart/internal/models"
	"github.com/Khasmag06/gophermart/internal/repository"
	"github.com/theplant/luhn"
	"net/http"
	"strconv"
)

func (s *Service) WithdrawPoints(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(int)

	var withdraws *models.Withdraws
	if err := json.NewDecoder(r.Body).Decode(&withdraws); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderNum, err := strconv.Atoi(withdraws.Order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if valid := luhn.Valid(orderNum); !valid {
		http.Error(w, repository.ErrOrderInvalidNum.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := s.repo.NewWithdrawal(userID, withdraws); err != nil {
		if errors.Is(err, repository.ErrBalanceNotEnoughPoints) {
			http.Error(w, err.Error(), http.StatusPaymentRequired)
			return
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}
