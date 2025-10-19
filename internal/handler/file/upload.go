package file

import (
	"crypto/sha256"
	"fmt"
	"hermes/config"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dokidokikoi/go-common/core"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h Handler) Upload(ctx *gin.Context) {
	logger := zaplog.From(ctx)
	var err error
	var fileName string
	defer func() {
		if err != nil {
			logger.With(zap.Error(err)).Error("")
			core.WriteResponse(ctx, err, nil)
		} else {
			core.WriteResponse(ctx, nil, gin.H{"path": fileName})
		}
	}()
	file, err := ctx.FormFile("file")
	if err != nil {
		return
	}
	f, err := file.Open()
	if err != nil {
		return
	}
	defer f.Close()

	fi, err := os.CreateTemp(config.Dir, "image_*"+filepath.Ext(file.Filename))
	if err != nil {
		return
	}
	defer fi.Close()

	hash := sha256.New()

	reader := io.TeeReader(f, hash)
	_, err = io.Copy(fi, reader)
	if err != nil {
		return
	}
	fileName = filepath.Join(config.Dir, fmt.Sprintf("%X", hash.Sum(nil))+filepath.Ext(file.Filename))
	err = os.Rename(fi.Name(), fileName)
	if err != nil {
		return
	}
	fileName = strings.TrimPrefix(fileName, config.Dir)
}
