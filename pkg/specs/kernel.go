/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package specs

import "gopkg.in/yaml.v3"

func (km *KernelModule) GetFieldVersion() string {
	ans, _ := km.Fields["version"]
	return ans
}

func (km *KernelModule) GetFieldLicense() string {
	ans, _ := km.Fields["license"]
	return ans
}

func (km *KernelModule) IsOpen() bool {
	l := km.GetFieldLicense()
	if l == "" || l == "NVIDIA" {
		return false
	}
	return true
}

func (km *KernelModule) Yaml() ([]byte, error) {
	return yaml.Marshal(km)
}
