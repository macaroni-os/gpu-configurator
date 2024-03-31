/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package pci

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v2"
)

func (s *SystemDevices) Yaml() ([]byte, error) {
	return yaml.Marshal(s)
}

func (s *SystemDevices) Json() ([]byte, error) {
	return json.Marshal(s)
}

func (s *SystemDevices) GetVGADevices() *[]*PCIDevice {
	ans := []*PCIDevice{}

	for _, device := range *s {
		words := strings.Split(device.ClassName, " ")
		if len(words) > 0 && words[0] == "VGA" {
			ans = append(ans, device)
		}
	}

	return &ans
}
