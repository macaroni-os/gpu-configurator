/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package backend

import (
	"fmt"

	bmacaroni "github.com/macaroni-os/gpu-configurator/pkg/backend/macaroni"
)

type SystemBackend interface {
	GetName() string
	GetEglExternalPlatformsDirs() ([]string, error)
	GetVulkanLayersDirs() ([]string, error)
	GetVulkanICDDirs() ([]string, error)

	// GBM stuff
	GetGBMLibDir() string
	GetEnvironmentDir() string

	// NVIDIA gpu functions
	GetNVIDIAEglWaylandLibDir() string
	GetNVIDIAEglGbmLibDir() string
	GetNVIDIADriverPrefixDir() string
}

func NewBackend(btype string) (SystemBackend, error) {
	var ans SystemBackend
	var err error
	switch btype {
	case "macaroni":
		ans, err = bmacaroni.NewMacaroniBackend()
	case "funtoo":
		ans, err = bmacaroni.NewMacaroniBackend()
		if ans != nil {
			(ans.(*bmacaroni.MacaroniBackend)).Name = "funtoo"
		}
	default:
		return ans, fmt.Errorf("%s backend not supported.", btype)
	}
	return ans, err
}
