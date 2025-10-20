package bangumi

import (
	"hermes/model"
	"hermes/tools"
	"maps"

	"github.com/pkg/errors"

	comm_tools "github.com/dokidokikoi/go-common/tools"
)

func (b *Bangumi) DoReq(method, uri string, header map[string]string, body any) ([]byte, error) {
	h := map[string]string{}
	b.RLock()
	maps.Copy(h, b.Headers)
	b.RUnlock()
	maps.Copy(h, header)

	query := comm_tools.GenQueryParams(body)
	if query != "" {
		uri += "?" + query
	}

	rsp, err := tools.ReqWithProxy(method, uri, body, b.Proxy, tools.SetHeadersWithOption(h))
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode()/100 != 2 {
		return nil, errors.Errorf("status code: %d, body: %s", rsp.StatusCode(), rsp.String())
	}
	return rsp.Bytes(), nil
}

func PersonRelation(career string) model.PersonRelation {
	switch career {
	case "artist":
		return model.PRelationMusic
	case "seiyu":
		return model.PRelationCV
	case "writer":
		return model.PRelationWriter
	case "illustrator":
		return model.PRelationPainter
	default:
		return model.PRelationUnknown
	}
}

func Gender(gender string) model.Gender {
	switch gender {
	case "male":
		return model.Male
	case "female":
		return model.Female
	case "futanari", "futa":
		return model.Futa
	default:
		return model.UnKnown
	}
}
