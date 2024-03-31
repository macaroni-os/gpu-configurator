/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package specs

func NewVulkanLayersFiles(dir string) *VulkanLayersFiles {
	return &VulkanLayersFiles{
		Path:  dir,
		Files: make(map[string]map[string]interface{}, 0),
	}
}
