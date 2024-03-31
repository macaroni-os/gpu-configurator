/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package backend

type MacaroniBackend struct {
	Name string
}

func NewMacaroniBackend() (*MacaroniBackend, error) {
	return &MacaroniBackend{
		Name: "macaroni",
	}, nil
}

func (b *MacaroniBackend) GetName() string {
	return b.Name
}

func (b *MacaroniBackend) GetEglExternalPlatformsDirs() ([]string, error) {
	return []string{
		"/usr/share/egl/egl_external_platform.d",
	}, nil
}

func (b *MacaroniBackend) GetVulkanLayersDirs() ([]string, error) {
	return []string{
		"/usr/share/vulkan/explicit_layer.d",
		"/usr/share/vulkan/implicit_layer.d",
	}, nil
}

func (b *MacaroniBackend) GetVulkanICDDirs() ([]string, error) {
	return []string{
		"/usr/share/vulkan/icd.d",
		"/etc/vulkan/icd.d",
	}, nil
}

func (b *MacaroniBackend) GetEnvironmentDir() string { return "/etc/env.d" }

func (b *MacaroniBackend) GetGBMLibDir() string { return "/usr/lib64/gbm" }

func (b *MacaroniBackend) GetNVIDIAEglWaylandLibDir() string { return "/usr/lib64" }
func (b *MacaroniBackend) GetNVIDIAEglGbmLibDir() string     { return "/usr/lib64" }
func (b *MacaroniBackend) GetNVIDIADriverPrefixDir() string  { return "/opt/nvidia" }
