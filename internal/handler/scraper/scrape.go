package scraper

import (
	"context"
	"izumi/internal/handler"

	"github.com/dokidokikoi/go-common/middleware"

	"github.com/google/uuid"
)

func (h Handler) Scrape(ctx context.Context, input *handler.ScraperDetailReq, op *middleware.PreHandleOptions) (string, error) {
	if input.RequestID == "" {
		input.RequestID = uuid.New().String()
	}
	requestID := input.RequestID
	_, err := h.srv.Scrap().Scrap(ctx, requestID, input.Objs)

	return requestID, err
}

func (h Handler) AutoScrape(ctx context.Context, input *handler.ScraperDetailReq, op *middleware.PreHandleOptions) (string, error) {

}
