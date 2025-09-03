package cookies

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strings"
)

var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
)

func Write(w http.ResponseWriter, cookie http.Cookie) error {
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	if len(cookie.String()) > 4096 {
		return ErrValueTooLong
	}

	http.SetCookie(w, &cookie)

	return nil
}

func WriteSigned(w http.ResponseWriter, cookie http.Cookie, secretKey []byte) error {
	log.Printf("Original cookie value: %s", cookie.Value)
	valB64 := base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(cookie.Name))
	mac.Write([]byte(valB64))
	signature := mac.Sum(nil)
	sigB64 := base64.URLEncoding.EncodeToString(signature)

	cookie.Value = sigB64 + "." + valB64
	log.Printf("Signed cookie value: %s", cookie.Value)

	return Write(w, cookie)
}

func Read(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", ErrInvalidValue
	}

	return string(value), nil
}

func ReadSigned(r *http.Request, name string, secretKey []byte) (string, error) {

	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	parts := strings.SplitN(cookie.Value, ".", 2)
	if len(parts) != 2 {
		return "", ErrInvalidValue
	}

	sigB64, valB64 := parts[0], parts[1]

	// проверяем подпись
	valueBytes, err := base64.URLEncoding.DecodeString(valB64)
	if err != nil {
		return "", ErrInvalidValue
	}

	signature, err := base64.URLEncoding.DecodeString(sigB64)
	if err != nil {
		return "", ErrInvalidValue
	}

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(cookie.Name))
	mac.Write([]byte(valB64))
	expectedSig := mac.Sum(nil)

	if !hmac.Equal(signature, expectedSig) {
		return "", ErrInvalidValue
	}

	return string(valueBytes), nil
}
