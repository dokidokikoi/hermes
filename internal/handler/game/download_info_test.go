package game

import (
	"izumi/internal/handler"
	"os"
	"testing"
)

func Test_DownloadInfo(t *testing.T) {
	err := downloadInfo(handler.GameVo{}, os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
}
