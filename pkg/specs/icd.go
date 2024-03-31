/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package specs

import (
	"encoding/json"
)

func NewICDJson(data []byte) (*ICDJson, error) {
	ans := &ICDJson{}
	if err := json.Unmarshal(data, &ans); err != nil {
		return nil, err
	}

	return ans, nil
}
