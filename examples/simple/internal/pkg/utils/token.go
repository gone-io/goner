package utils

import (
	"encoding/base64"
	"github.com/google/uuid"
)

func CreateToken() string {
	u := uuid.New()
	return base64.URLEncoding.EncodeToString(u[0:])
}
