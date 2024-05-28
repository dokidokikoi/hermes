package tools

import (
	"crypto/sha256"
	"fmt"
	"hermes/config"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func SaveTmpFile(ext string, data io.Reader) (string, error) {
	path := filepath.Join(config.TmpDir, fmt.Sprintf("%s_%d%s", time.Now().Format("20060102150405"), rand.Intn(100000), ext))
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	r := io.TeeReader(data, f)
	h := sha256.New()
	_, err = io.Copy(h, r)
	if err != nil {
		return "", err
	}
	newPath := filepath.Join(config.TmpDir, fmt.Sprintf("%X%s", h.Sum(nil), ext))
	err = os.Rename(path, newPath)
	if err != nil {
		return "", err
	}
	return newPath, nil
}
