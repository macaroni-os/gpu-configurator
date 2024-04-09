/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package specs

import (
	"encoding/json"
	"strings"
)

func NewJsonFile(n string, f *ICDJson) *JsonFile {
	ans := &JsonFile{
		Name:     n,
		File:     f,
		Disabled: false,
	}

	if strings.HasSuffix(n, ".disabled") {
		ans.Name = ans.Name[0 : len(ans.Name)-9]
		ans.Disabled = true
	}

	return ans
}

func NewICDJson(data []byte) (*ICDJson, error) {
	ans := &ICDJson{}
	if err := json.Unmarshal(data, &ans); err != nil {
		return nil, err
	}

	return ans, nil
}
