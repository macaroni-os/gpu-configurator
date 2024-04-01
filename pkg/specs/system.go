/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package specs

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

func NewSystem() *System {
	return &System{
		EglExtPlatformDirs: []*EglExternalPlatformFiles{},
		VulkanLayersDirs:   []*VulkanLayersFiles{},
		VulkanICDDirs:      []*EglExternalPlatformFiles{},
		GbmLibraries:       []*Library{},
	}
}

func (s *System) Yaml() ([]byte, error) {
	return yaml.Marshal(s)
}

func (s *System) Json() ([]byte, error) {
	return json.Marshal(s)
}

func (s *System) GetGBMLibrary(lib string) *Library {
	var ans *Library = nil

	if len(s.GbmLibraries) > 0 {
		for idx := range s.GbmLibraries {
			if s.GbmLibraries[idx].Name == lib {
				ans = s.GbmLibraries[idx]
				break
			}
		}
	}

	return ans
}
