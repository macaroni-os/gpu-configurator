/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package macaroni

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/macaroni-os/gpu-configurator/pkg/specs"
	"github.com/macaroni-os/macaronictl/pkg/utils"
)

var (
	udevNvidiaUdevsh = `#!/bin/sh
# Autogenered file by gpu-configurator.

restorecon_nvidia()
{
	if [ -x /sbin/restorecon ]; then
		ebegin "Restoring SELinux contexts for /dev/nvidia*"
		restorecon -F /dev/nvidia* >/dev/null 2>&1
		eend $?
	fi

	return 0
}

if [ $# -ne 1 ]; then
	echo "Invalid args" >&2
	exit 1
fi

case $1 in
	add|ADD)
		#hopefully this prevents infinite loops like bug #454740
		if lsmod | grep -iq nvidia; then
			/usr/bin/nvidia-smi > /dev/null
			restorecon_nvidia
		fi
		;;
	remove|REMOVE)
		rm -f /dev/nvidia*
		;;
esac

exit 0
`
)

func (b *MacaroniBackend) createUdevScript() error {
	targetDir := "/lib/udev"
	rulesDir := "/lib/udev/rules.d"
	targetFile := "nvidia-udev.sh"

	if !utils.Exists(targetDir) {
		err := os.MkdirAll(targetDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	udevScript := filepath.Join(targetDir, targetFile)
	err := os.WriteFile(udevScript, []byte(udevNvidiaUdevsh),
		0755)

	if err != nil {
		return fmt.Errorf("Error on write %s: %s",
			targetDir, targetFile)
	}

	// Create /lib/udev/rules.d/99-nvidia.rules file
	udevRule := filepath.Join(rulesDir, "99-nvidia.rules")

	if !utils.Exists(rulesDir) {
		err := os.MkdirAll(rulesDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(udevRule, []byte(
		`ACTION=="add", DEVPATH=="/module/nvidia", SUBSYSTEM=="module", RUN+="nvidia-udev.sh $env{ACTION}"
# Previously the ACTION was "add|remove" but one user on bug #376527 had a
# problem until he recompiled udev-171-r5, which is one of the versions I
# tested with and it was fine. I'm breaking the rules out just to be safe
# so someone else doesn't have an issue
ACTION=="remove", DEVPATH=="/module/nvidia", SUBSYSTEM=="module", RUN+="nvidia-udev.sh $env{ACTION}"`),
		0644,
	)
	if err != nil {
		return fmt.Errorf("error on write file %s: %s",
			udevRule, err.Error())
	}

	return nil
}

func (b *MacaroniBackend) SetNVIDIAModprobeFiles(setup *specs.NVIDIASetup, withVideoGroup, force bool) error {
	modprobeDir := "/etc/modprobe.d"

	rmmodfile := "nvidia-rmmod.conf"
	modfile := "nvidia-conf"

	if !utils.Exists(modprobeDir) {
		err := os.MkdirAll(modprobeDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	modfileabs := filepath.Join(modprobeDir, modfile)
	rmmodfileabs := filepath.Join(modprobeDir, rmmodfile)

	// Creating the file /etc/modprobe.d/nvidia.conf if not present.
	if !utils.Exists(modfileabs) || force {
		var content string

		// In macaroni we use static gid for video group.
		// We can retrieve it with in the near future:
		// $> entities list groups --filter video --json
		vgroup := "27"

		if withVideoGroup {

			content = fmt.Sprintf(`# Nvidia drivers support
alias char-major-195 nvidia
alias /dev/nvidiactl char-major-195

# This configures the device nodes to be part of the video group.
options nvidia NVreg_DeviceFileGID=%s NVreg_DeviceFileMode=432 NVreg_DeviceFileUID=0 NVreg_ModifyDeviceFiles=1
`, vgroup)

		} else {

			content = `# Nvidia drivers support
alias char-major-195 nvidia
alias /dev/nvidiactl char-major-195`

		}
		err := os.WriteFile(modfileabs, []byte(content), 0644)
		if err != nil {
			return fmt.Errorf("error on write file %s: %s",
				modfileabs, err.Error())
		}
	}

	if !utils.Exists(rmmodfileabs) || force {

		content := `# Nvidia UVM support
remove nvidia modprobe -r --ignore-remove nvidia-drm nvidia-modeset nvidia-uvm nvidia
`

		err := os.WriteFile(rmmodfileabs, []byte(content), 0644)
		if err != nil {
			return fmt.Errorf("error on write file %s: %s",
				rmmodfileabs, err.Error())
		}
	}

	return nil
}
