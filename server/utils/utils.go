package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const Letters = "abcdefghijklmnopqrstuvwxyzABCDEGHIJKLMNOPQRSTUVWXYZ1234567890"

func GenerateRandomString(length int) (string, error) {
	var result string
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(Letters))))
		if err != nil {
			return "", err
		}
		result += string(Letters[index.Int64()])
	}

	return result, nil
}

func GenerateRoomID(prID, userID string, length int) (string, error) {
	randomstring, err := GenerateRandomString(length)
	if err != nil {
		return "", err
	}
	roomID := fmt.Sprintf("%s-%s-%s", prID, userID, randomstring)
	return roomID, nil
}
