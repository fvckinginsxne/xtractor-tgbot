package coding

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func EncodeUsernameAndTitle(username, title string) string {
	combined := fmt.Sprintf("%s_%s", username, title)

	hash := sha256.Sum256([]byte(combined))

	hashString := hex.EncodeToString(hash[:])

	shortHash := hashString[:20]

	return shortHash
}
