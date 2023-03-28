package handlers

import (
	"github.com/Khasmag06/gophermart/config"
	mw "github.com/Khasmag06/gophermart/internal/middlewares"
	"github.com/Khasmag06/gophermart/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Service struct {
	cfg  config.Config
	repo repository.Storage
}

func NewService(cfg config.Config, repo repository.Storage) *Service {
	return &Service{cfg, repo}
}

func (s *Service) Route() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(mw.GzipHandle)

	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/register", s.Register)
			r.Post("/login", s.Login)

			r.Group(func(r chi.Router) {
				r.Use(mw.CheckToken)
				r.Post("/orders", s.NewOrder)
				r.Get("/orders", s.GetOrders)
				r.Get("/balance", s.GetBalance)
				r.Post("/balance/withdraw", s.WithdrawPoints)
				r.Get("/withdrawals", s.GetWithdrawals)
			})

		})
	})
	return r
}
