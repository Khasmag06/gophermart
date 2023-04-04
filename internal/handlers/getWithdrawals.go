package handlers

import (
	"encoding/json"
	"github.com/Khasmag06/gophermart/internal/middlewares"
	"net/http"
)

func (s *Service) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(int)

	withdrawals, err := s.repo.GetWithdrawals(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.Write([]byte("no data for return"))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(withdrawals); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
