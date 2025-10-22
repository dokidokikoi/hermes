package tag

import (
	"izumi/db/data"
	"izumi/internal/service"
)

type Handler struct {
	srv service.Iservice
}

func NewHandler() Handler {
	return Handler{srv: service.NewSrv(data.GetDataFactory())}
}
