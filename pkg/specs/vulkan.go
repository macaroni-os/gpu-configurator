/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package specs

import "strings"

func NewVulkanLayersFiles(dir string) *VulkanLayersFiles {
	return &VulkanLayersFiles{
		Path:  dir,
		Files: make(map[string]*VulkanLayersFile, 0),
	}
}

func NewVulkanLayersFile(n string) *VulkanLayersFile {
	ans := &VulkanLayersFile{
		Name:     n,
		Disabled: false,
		Content:  make(map[string]interface{}, 0),
	}

	if strings.HasSuffix(n, ".disabled") {
		ans.Disabled = true
		ans.Name = ans.Name[0 : len(ans.Name)-9]
	}

	return ans
}
