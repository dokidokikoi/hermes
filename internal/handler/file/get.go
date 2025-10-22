package file

import (
	"io"
	"izumi/config"
	"os"
	"path/filepath"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	"github.com/gin-gonic/gin"
)

func (h Handler) Get(ctx *gin.Context) {
	fileName := ctx.Param("name")
	f, err := os.Open(filepath.Join(config.TmpDir, fileName))
	if err != nil {
		f, err = os.Open(filepath.Join(config.DataDir, fileName))
		if err != nil {
			core.WithErr(ctx, errors.ApiErrValidation)
			core.WriteResponse(ctx, err, nil)
			return
		}
	}
	io.Copy(ctx.Writer, f)
}
