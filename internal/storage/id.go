package storage

import (
	"math/rand"
	"time"
)

const (
	IDChars  = "abcdefghijklmnopqrstuvwxyz0123456789"
	IDLength = 7
)

func GenerateRawID() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, IDLength)
	for i := range b {
		b[i] = IDChars[rand.Intn(len(IDChars))]
	}
	return string(b)
}

func GetFullID(storageID, rawID string) string {
	// The PHP implementation seems to mix the storage ID into the raw ID
	// Let's simplify for now: storageID + rawID
	return storageID + rawID
}
