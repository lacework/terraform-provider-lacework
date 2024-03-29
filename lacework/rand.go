package lacework

import (
	"math/rand"
	"time"
)

var (
	charset              = "abcdefghijklmnopqrstuvwxyz0123456789"
	randomSeed           = rand.New(rand.NewSource(time.Now().UnixNano()))
	charsetAWSExternalID = "=,.@:/-abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

func randomString(length int) string {
	return stringFromCharset(length, charset)
}

func randomStringExternalID(length int) string {
	return stringFromCharset(length, charsetAWSExternalID)
}

func stringFromCharset(length int, charset string) string {
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = charset[randomSeed.Intn(len(charset))]
	}
	return string(bytes)
}
