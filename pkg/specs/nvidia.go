/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package specs

func NewNVIDIASetup() *NVIDIASetup {
	return &NVIDIASetup{
		Drivers:       []*NVIDIADriver{},
		VersionActive: "",
	}
}

func (n *NVIDIASetup) HasVersion(v string) bool {
	ans := false
	for idx := range n.Drivers {
		if n.Drivers[idx].Version == v {
			ans = true
			break
		}
	}
	return ans
}

func (n *NVIDIASetup) ElaborateVersion() {
	for idx := range n.Drivers {
		if n.Drivers[idx].WithKernelModules {
			n.VersionActive = n.Drivers[idx].Version
			break
		}
	}
}
