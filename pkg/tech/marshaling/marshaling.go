package marshaling

import (
	"encoding/json"

	"bot/pkg/tech/e"
)

func DataToJSON(data map[string]interface{}) ([]byte, error) {
	replyMarkupJSON, err := json.Marshal(data)
	if err != nil {
		return nil, e.Wrap("can't marshal data to JSON", err)
	}

	return replyMarkupJSON, nil
}
