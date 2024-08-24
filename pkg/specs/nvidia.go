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

func (n *NVIDIASetup) SetVersion(v string) { n.VersionActive = v }

func (n *NVIDIASetup) HasVersion(v string) bool {
	ans := false
	if d := n.GetDriver(v); d != nil {
		ans = true
	}
	return ans
}

func (n *NVIDIASetup) GetDriver(v string) *NVIDIADriver {
	for idx := range n.Drivers {
		if n.Drivers[idx].Version == v {
			return n.Drivers[idx]
		}
	}
	return nil
}

func (n *NVIDIASetup) GetKernelModulesAvailable(nv string, open bool) *[]*KernelModule {
	ans := []*KernelModule{}

	if open {
		for _, km := range n.KOpenModuleAvailable {
			if km.GetFieldVersion() == nv {
				ans = append(ans, km)
			}
		}
	} else {
		for _, km := range n.KModuleAvailable {
			if km.GetFieldVersion() == nv {
				ans = append(ans, km)
			}
		}
	}

	return &ans
}

func (n *NVIDIASetup) GetKernelModulesActive(nv, kv string) *KernelModule {
	// Search on proprietary modules
	for _, km := range n.KModuleActive {
		if km.KernelVersion == kv && km.GetFieldVersion() == nv {
			return km
		}
	}

	// Search on open modules
	for _, km := range n.KOpenModuleActive {
		if km.KernelVersion == kv && km.GetFieldVersion() == nv {
			return km
		}
	}

	return nil
}
