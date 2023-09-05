package client

import (
	"encoding/json"
	"fmt"
	"github.com/Khasmag06/gophermart/internal/models"
	"github.com/Khasmag06/gophermart/internal/repository"
	"io"
	"net/http"
)

type Accrual struct {
	OrderQueue           chan string
	accrualSystemAddress string
	repo                 repository.Storage
}

func NewAccrual(accrualSystemAddress string, repo repository.Storage) *Accrual {
	orderQueue := make(chan string, 10)
	return &Accrual{
		OrderQueue:           orderQueue,
		accrualSystemAddress: accrualSystemAddress,
		repo:                 repo,
	}
}

func (a *Accrual) Run() {
	go func() {
		for order := range a.OrderQueue {
			if err := a.updateOrderData(order); err != nil {
				a.OrderQueue <- order
			}
		}
	}()
}

func (a *Accrual) updateOrderData(order string) error {
	url := fmt.Sprintf("%s/api/orders/%s", a.accrualSystemAddress, order)
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case http.StatusOK:
		defer response.Body.Close()
		payload, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var accrual models.AccrualData
		if err := json.Unmarshal(payload, &accrual); err != nil {
			return err
		}

		err = a.repo.UpdateAccrual(accrual)
		if err != nil {
			return err
		}
	case http.StatusTooManyRequests:
		a.OrderQueue <- order
	case http.StatusInternalServerError:
		a.OrderQueue <- order
	}

	return nil
}
