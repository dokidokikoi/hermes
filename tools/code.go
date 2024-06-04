package tools

import (
	"encoding/base64"
	"net/url"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func Jp2Utf8(originBytes []byte) (string, error) {
	eucJPDecoder := japanese.EUCJP.NewDecoder()
	utf8Bytes, _, err := transform.Bytes(eucJPDecoder, originBytes)
	if err != nil {
		return string(originBytes), err
	}
	return string(utf8Bytes), err
}

func Utf82Jp(originBytes []byte) (string, error) {
	eucJPEncoder := japanese.EUCJP.NewEncoder()
	utf8Bytes, _, err := transform.Bytes(eucJPEncoder, originBytes)
	if err != nil {
		return string(originBytes), err
	}
	return string(utf8Bytes), nil
}

func Base64Enc(originBytes []byte) string {
	return base64.StdEncoding.EncodeToString(originBytes)
}

func Base64Dec(origin string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(origin)
}

func UrlEnc(origin string) string {
	return url.QueryEscape(origin)
}

func UrlDec(origin string) (string, error) {
	return url.QueryUnescape(origin)
}
