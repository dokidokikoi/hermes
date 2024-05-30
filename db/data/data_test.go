package data

import (
	"context"
	"fmt"
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
		Name: "test",
		Tags: []model.Tag{
			{
				Name: "喜剧",
			},
			{
				Name: "哲学",
			},
		},
		Series: []model.Series{
			{
				Name: "喜剧",
			},
		},
		Category: &model.Category{
			Name: "ADV",
		},
		Alias: []string{"sdafs", "sfgrwg"},
		Links: []model.Link{
			{
				Name: "res",
				Url:  "http://test.test",
			},
			{
				Name: "res",
				Url:  "http://test.test",
			},
		},
	}
	tx := GetDataFactory().Transaction().Begin()
	err := tx.Game().Create(context.Background(), &game, nil)
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

func TestGet(t *testing.T) {
	config.SetConfig("../../conf/application.yaml")
	g, err := GetDataFactory().Game().Get(context.Background(), &model.Game{}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("game: %+v", g)
}
