package person

import (
	"context"
	"izumi/db/data"
	"izumi/internal/handler"
	"izumi/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h Handler) Get(ctx context.Context, req *struct{}) (*handler.StaffVo, error) {
	id, err := strconv.ParseUint(ctx.(*gin.Context).Param("id"), 10, 32)
	if err != nil {
		return nil, err
	}
	c, err := data.GetDataFactory().Person().Get(ctx, &model.Person{ID: uint(id)}, nil)
	if err != nil {
		return nil, err
	}

	vo := handler.StaffVo{
		ID:        c.ID,
		Name:      c.Name,
		Alias:     c.Alias,
		Gender:    c.Gender,
		Summary:   c.Summary,
		Cover:     c.Cover,
		Images:    c.Images,
		Tags:      c.Tags,
		CreatedAt: c.CreatedAt,
	}

	return &vo, nil
}
