package errors

import (
	"encoding/json"
)

type Unauthorized struct{}

func (err *Unauthorized) String() string {
	return "Please login first"
}

func (err *Unauthorized) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"error": err.String(),
	})
}
