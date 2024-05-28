package data

import (
	"context"
	"hermes/config"
	"hermes/model"
	"testing"
)

func TestConn(t *testing.T) {
	config.SetConfig("../../conf/application.yaml")
	GetDataFactory()
}

func TestTranscation(t *testing.T) {
	config.SetConfig("../../conf/application.yaml")
	game := model.Game{
		ID:   4,
		Name: "test",
		Tags: []model.Tag{
			{
				ID:   1,
				Name: "喜剧",
			},
			{
				ID:   2,
				Name: "哲学",
			},
		},
		Series: []model.Series{
			{
				ID:   1,
				Name: "喜剧",
			},
		},
		Category: &model.Category{
			ID:   3,
			Name: "ADV",
		},
	}
	tx := GetDataFactory().Transaction().Begin()
	err := tx.Game().Update(context.Background(), &game, nil)
	if err != nil {
		panic(err)
	}

	txx := GetDataFactory().Transaction().Begin()
	err = txx.Category().Create(context.Background(), &model.Category{Name: "RPGg"}, nil)
	if err != nil {
		panic(err)
	}

	tx.Transaction().Commit()
	txx.Transaction().Commit()
}
