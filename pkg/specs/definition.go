/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package specs

type System struct {
	EglExtPlatformDirs []*EglExternalPlatformFiles `json:"egl_external_platforms_dirs,omitempty" yaml:"egl_external_platforms_dirs,omitempty"`
	VulkanLayersDirs   []*VulkanLayersFiles        `json:"vulkan_layers_dirs,omitempty" yaml:"vulkan_layers_dirs,omitempty"`
	VulkanICDDirs      []*EglExternalPlatformFiles `json:"vulkan_icd_dirs,omitempty" yaml:"vulkan_icd_dirs,omitempty"`
}

type VulkanLayersFiles struct {
	Path  string                            `json:"path" yaml:"path"`
	Files map[string]map[string]interface{} `json:"files,omitempty" yaml:"files,omitempty"`
}

type EglExternalPlatformFiles struct {
	Path  string              `json:"path" yaml:"path"`
	Files map[string]*ICDJson `json:"files,omitempty" yaml:"files,omitempty"`
}

type ICDJson struct {
	FileFormatVersion string      `json:"file_format_version" yaml:"file_format_version"`
	ICD               ICDJsonData `json:"ICD" yaml:"ICD"`
}

type ICDJsonData struct {
	LibraryPath string `json:"library_path" yaml:"library_path"`
	ApiVersion  string `json:"api_version,omitempty" yaml:"api_version,omitempty"`
}
