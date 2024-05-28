package tools

import (
	"fmt"
	"testing"
)

func TestStr2Time(t *testing.T) {
	fmt.Println(Str2Time("2019年5月1日(金)").Format("2006-01-02"))
}
