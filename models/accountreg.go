package models

import (
	"crypto/sha256"
	"fmt"
	"time"
)


// Sha256 генерирует аунтификационный код
func Sha256() string {
	h := time.Now().String()
	sha := sha256.New()
	sha.Write([]byte(h))
	return fmt.Sprintf("%x", sha.Sum(nil))
}