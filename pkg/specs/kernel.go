/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package specs

func (km *KernelModule) GetFieldVersion() string {
	ans, _ := km.Fields["version"]
	return ans
}
