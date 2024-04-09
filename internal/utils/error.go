package utils

import (
	"encoding/json"
)

// ParamToJSON converts param with requested name to JSON output format.
func ParamToJSON(name string, desc string) []byte {
	resp := map[string]string{
		name: desc,
	}
	out, _ := json.Marshal(resp)

	return out
}
