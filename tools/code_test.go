package tools

import (
	"fmt"
	"testing"
)

func TestUtf82Jp(t *testing.T) {
	fmt.Println(UrlEnc(Utf82Jp([]byte("ボクの彼女はガテン系"))))
}
