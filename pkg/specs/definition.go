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

	Nvidia *NVIDIASetup `json:"nvidia,omitempty" yaml:"nvidia,omitempty"`
}

type NVIDIASetup struct {
	Drivers              []*NVIDIADriver `json:"drivers,omitempty" yaml:"drivers,omitempty"`
	VersionActive        string          `json:"version_active,omitempty" yaml:"version_active,omitempty"`
	KModuleAvailable     []*KernelModule `json:"kernel_modules,omitempty" yaml:"kernel_modules,omitempty"`
	KOpenModuleAvailable []*KernelModule `json:"kernel_open_modules,omitempty" yaml:"kernel_open_modules,omitempty"`
}

type NVIDIADriver struct {
	Path              string `json:"path" yaml:"path"`
	Version           string `json:"version" yaml:"version"`
	WithKernelModules bool   `json:"with_kernel_modules,omitempty" yaml:"with_kernel_modules,omitempty"`
}

type VulkanLayersFiles struct {
	Path  string                       `json:"path" yaml:"path"`
	Files map[string]*VulkanLayersFile `json:"files,omitempty" yaml:"files,omitempty"`
}

type EglExternalPlatformFiles struct {
	Path  string               `json:"path" yaml:"path"`
	Files map[string]*JsonFile `json:"files,omitempty" yaml:"files,omitempty"`
}

type JsonFile struct {
	Name     string   `json:"name" yaml:"name"`
	Disabled bool     `json:"disabled,omitempty" yaml:"disabled,omitempty"`
	File     *ICDJson `json:"file,omitempty" yaml:"file,omitempty"`
}

type ICDJson struct {
	FileFormatVersion string      `json:"file_format_version" yaml:"file_format_version"`
	ICD               ICDJsonData `json:"ICD" yaml:"ICD"`
}

type VulkanLayersFile struct {
	Name     string                 `json:"name" yaml:"name"`
	Disabled bool                   `json:"disabled,omitempty" yaml:"disabled,omitempty"`
	Content  map[string]interface{} `json:"content,omitempty" yaml:"content,omitempty"`
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

type KernelModule struct {
	Path          string            `json:"path,omitempty" yaml:"path,omitempty"`
	KernelVersion string            `json:"kernel_version,omitempty" yaml:"kernel_version,omitempty"`
	Fields        map[string]string `json:"fields,omitempty" yaml:"fields,omitempty"`
	Name          string            `json:"name,omitempty" yaml:"name,omitempty"`
}
