package system_task

import (
	"izumi/db"
	"izumi/internal/service"
)

type Handler struct {
	srv service.Iservice
}

func NewHandler() Handler {
	return Handler{srv: service.NewSrv(db.GetStore())}
}
