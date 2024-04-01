/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package specs

type System struct {
	EglExtPlatformDirs []*EglExternalPlatformFiles `json:"egl_external_platforms_dirs,omitempty" yaml:"egl_external_platforms_dirs,omitempty"`
	VulkanLayersDirs   []*VulkanLayersFiles        `json:"vulkan_layers_dirs,omitempty" yaml:"vulkan_layers_dirs,omitempty"`
	VulkanICDDirs      []*EglExternalPlatformFiles `json:"vulkan_icd_dirs,omitempty" yaml:"vulkan_icd_dirs,omitempty"`

	GbmLibraries []*Library `json:"gbm_libs,omitempty" yaml:"gbm_libs,omitempty"`

	Nvidia NVIDIASetup `json:"nvidia,omitempty" yaml:"nvidia,omitempty"`
}

type NVIDIASetup struct {
	Drivers       []*NVIDIADriver `json:"drivers,omitempty" yaml:"drivers,omitempty"`
	VersionActive string          `json:"version_active,omitempty" yaml:"version_active,omitempty"`
}

type NVIDIADriver struct {
	Path    string `json:"path" yaml:"path"`
	Version string `json:"version" yaml:"version"`
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

type Library struct {
	Name       string `json:"library" yaml:"library"`
	Disabled   bool   `json:"disabled,omitempty" yaml:"disabled,omitempty"`
	LinkedFile string `json:"linked_libpath,omitempty" yaml:"linked_libpath,omitempty'`
}
