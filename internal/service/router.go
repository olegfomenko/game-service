package service

import (
	"github.com/go-chi/chi"
	"github.com/olegfomenko/game-service/internal/service/handlers"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log), // this line may cause compilation error but in general case `dep ensure -v` will fix it
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxConnector(s.connector),
		),
	)

	r.Route("/integrations/game-service", func(r chi.Router) {
		r.Post("/create_game", handlers.CreateGame)
	})

	return r
}
