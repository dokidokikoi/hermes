package developer

import (
	"hermes/db/data"
	"hermes/internal/service"
)

type Handler struct {
	srv service.Iservice
}

func NewHandler() Handler {
	return Handler{srv: service.NewSrv(data.GetDataFactory())}
}
