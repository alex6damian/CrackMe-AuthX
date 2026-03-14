package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// GeneratePredictableResetToken creates a predictable reset token based on the username and the current date.
// This is intentionally insecure for the purpose of the CrackMe challenge.
// The token is the MD5 hash of the username concatenated with a fixed secret and the current date (YYYY-MM-DD).
func GeneratePredictableResetToken(username string) string {
	// Concatenate the predictable parts.
	dataToHash := username

	// Hash the data using MD5.
	hasher := md5.New()
	hasher.Write([]byte(dataToHash))
	hashBytes := hasher.Sum(nil)

	// Return the hash as a hex string.
	return hex.EncodeToString(hashBytes)
}
