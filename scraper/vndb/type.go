package vndb

import (
	"izumi/internal/handler"
	"izumi/model"
)

type BaseResponse[T any] struct {
	Result         []T    `json:"results"`
	More           bool   `json:"more"`
	Count          int    `json:"count"`
	CompactFilters string `json:"compact_filters"`
}

type VN struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Titles []struct {
		Title    string `json:"title"`
		Official bool   `json:"official"`
	} `json:"titles"`
	Released  string   `json:"released"`
	Languages []string `json:"languages"`
	Platforms []string `json:"platforms"`
	Image     struct {
		Url string `json:"url"`
	} `json:"image"`
	LengthMinutes int     `json:"length_minutes"`
	Description   string  `json:"description"`
	Rating        float32 `json:"rating"`
	Screenshots   []struct {
		Url string `json:"url"`
	} `json:"screenshots"`
	Tags       []Tag      `json:"tags"`
	Developers []Producer `json:"developers"`
	Staff      []Staff    `json:"staff"`
	VA         []VA       `json:"va"`
	ExtLinks   []ExtLink  `json:"extlinks"`
}

type Character struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Original    string   `json:"original"`
	Aliases     []string `json:"aliases"`
	Description string   `json:"description"`
	Image       struct {
		Url string `json:"url"`
	} `json:"image"`
	BloodType string   `json:"blood_type"`
	Height    int      `json:"height"`
	Weight    int      `json:"weight"`
	Bust      int      `json:"bust"`
	Waist     int      `json:"waist"`
	Hips      int      `json:"hips"`
	Cup       string   `json:"cup"`
	Age       int      `json:"age"`
	Birthday  []int    `json:"birthday"`
	Gender    []string `json:"gender"`
	Traits    []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"traits"`
}

type Staff struct {
	ID          string    `json:"id"`
	Role        string    `json:"role"`
	Description string    `json:"description"`
	Gender      string    `json:"gender"`
	Name        string    `json:"name"`
	Original    string    `json:"original"`
	ExtLinks    []ExtLink `json:"extlinks"`
}

type Tag struct {
	Name    string   `json:"name"`
	Aliases []string `json:"aliases"`
}

type Producer struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	ExtLinks []ExtLink `json:"extlinks"`
}

type VA struct {
	Character Character `json:"character"`
	Staff     Staff     `json:"staff"`
	Note      string    `json:"note"`
}

type ExtLink struct {
	Url   string `json:"url"`
	Label string `json:"label"`
}

func TransfGender(gender string) model.Gender {
	switch gender {
	case "f":
		return model.Female
	case "m":
		return model.Male
	case "o":
		return model.Futa
	default:
		return model.UnKnown
	}
}

func TransfRelation(role string) model.PersonRelation {
	switch role {
	case "director", "editor":
		return model.PRelationWriter
	case "chardesign", "art":
		return model.PRelationPainter
	case "music", "songs":
		return model.PRelationMusic
	default:
		return model.PRelationUnknown
	}
}

func Staff2Staff(staff Staff) handler.StaffVo {
	return handler.StaffVo{
		VNDBID: staff.ID,
		Name: func() string {
			if staff.Original != "" {
				return staff.Original
			}
			return staff.Name
		}(),
		Alias: func() []string {
			if staff.Original == "" {
				return []string{}
			}
			return []string{staff.Name}
		}(),
		Summary: staff.Description,
		Gender: func() model.Gender {
			return TransfGender(staff.Gender)
		}(),
		Relation: func() []model.PersonRelation {
			return []model.PersonRelation{TransfRelation(staff.Role)}
		}(),
	}
}

func Link2Link(extLink ExtLink) model.Link {
	return model.Link{
		Name: extLink.Label,
		Url:  extLink.Url,
	}
}
