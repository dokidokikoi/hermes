package game

import (
	"hermes/internal/handler"
	"os"
	"testing"
)

func Test_DownloadInfo(t *testing.T) {
	err := downloadInfo(handler.GameVo{}, os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
}
