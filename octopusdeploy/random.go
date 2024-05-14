package octopusdeploy

import (
	"crypto/rand"
	"log"
	"math/big"
)

const (
	letterBytes   = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
)

func generateRandomBytes(length int) []byte {
	var randomBytes = make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal("Unable to generate random bytes")
	}
	return randomBytes
}

func generateRandomCryptoString(length int) string {
	result := make([]byte, length)
	bufferSize := int(float64(length) * 1.3)
	for i, j, randomBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			randomBytes = generateRandomBytes(bufferSize)
		}
		if idx := int(randomBytes[j%length] & letterIdxMask); idx < len(letterBytes) {
			result[i] = letterBytes[idx]
			i++
		}
	}

	return string(result)
}

func generateRandomSerialNumber() big.Int {
	random := rand.Reader
	randomSerialNumber, err := rand.Int(random, big.NewInt(9223372036854775807))
	if err != nil {
		log.Fatal("Unable to generate random serial number")
	}
	return *randomSerialNumber
}
