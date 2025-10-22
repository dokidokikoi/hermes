package character

import (
	"context"
	"errors"
	"izumi/db/data"
	"izumi/model"
	"strings"

	comm_errors "github.com/dokidokikoi/go-common/errors"
)

func (h Handler) Update(ctx context.Context, input *model.Character) (any, error) {
	tx := data.GetDataFactory().Transaction().Begin()
	err := tx.CharacterTag().Delete(ctx, &model.CharacterTag{CharacterID: input.ID}, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return nil, err
	}

	tags := []*model.CharacterTag{}
	for _, tag := range input.Tags {
		tags = append(tags, &model.CharacterTag{
			CharacterID: input.ID,
			TagID:       tag.ID,
		})
	}
	err = tx.CharacterTag().Creates(ctx, tags, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return nil, err
	}

	input.Name = strings.TrimSpace(input.Name)
	if err := tx.Character().Update(ctx, input, nil); err != nil {
		if !errors.Is(err, comm_errors.ErrNoUpdateRows) {
			tx.Transaction().Rollback()
			return nil, err
		}
	}

	return nil, nil
}
