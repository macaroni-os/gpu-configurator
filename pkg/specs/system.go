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

func (s *System) GetEglLoader(jsonfile string) (*EglExternalPlatformFiles, *JsonFile) {
	if len(s.EglExtPlatformDirs) > 0 {
		for idx := range s.EglExtPlatformDirs {
			if json, present := s.EglExtPlatformDirs[idx].Files[jsonfile]; present {
				return s.EglExtPlatformDirs[idx], json
			}
		}
	}

	return nil, nil
}

func (s *System) GetVulkanIcdFile(jsonfile string) (*EglExternalPlatformFiles, *JsonFile) {
	if len(s.VulkanICDDirs) > 0 {
		for idx := range s.VulkanICDDirs {
			if json, present := s.VulkanICDDirs[idx].Files[jsonfile]; present {
				return s.VulkanICDDirs[idx], json
			}
		}
	}

	return nil, nil
}

func (s *System) GetVulkanLayerFile(jsonfile string) (*VulkanLayersFiles, *VulkanLayersFile) {
	if len(s.VulkanLayersDirs) > 0 {
		for idx := range s.VulkanLayersDirs {
			if json, present := s.VulkanLayersDirs[idx].Files[jsonfile]; present {
				return s.VulkanLayersDirs[idx], json
			}
		}
	}

	return nil, nil
}
