package tools

import (
	"bytes"
	"strings"
	"time"

	comm_tools "github.com/dokidokikoi/go-common/tools"
)

func Str2Time(str string) (t time.Time) {
	str = comm_tools.TrimBlankChar(str)
	bs := []byte(str)
	arr := make([]string, 3)
	index := 0
	buf := bytes.Buffer{}
	for _, b := range bs {
		if b >= 48 && b <= 57 {
			buf.WriteByte(b)
		} else if buf.Len() > 0 {
			arr[index] = buf.String()
			buf.Reset()
			index++
			if index >= 3 {
				break
			}
		}
	}
	if index < 3 {
		arr[index] = buf.String()
		buf.Reset()
	}
	if len(arr) < 3 || len(arr[0]) != 4 || len(arr[1]) > 2 || len(arr[1]) < 1 || len(arr[2]) > 2 || len(arr[2]) < 1 {
		return
	}
	if len(arr[1]) < 2 {
		arr[1] = "0" + arr[1]
	}
	if len(arr[2]) < 2 {
		arr[2] = "0" + arr[2]
	}
	str = strings.Join(arr, "-")

	t, _ = time.Parse("2006-01-02", str[:10])
	return
}
