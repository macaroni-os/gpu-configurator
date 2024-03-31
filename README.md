<p align="center">
  <img src="https://github.com/macaroni-os/macaroni-site/blob/master/site/static/images/logo.png">
</p>

# GPU Configurator

The `gpu-configurator` tool is born to organize and help users configure
their GPU cards to run correctly on Xorg or Xwayland.

It's mainly usable on Macaroni OS and/or Funtoo/Gentoo environments
but the code is been organized to add support for additional OS.

It's under heavy development and new features will be added to
cover in a better way multiple use cases, but it's first target
is configured NVIDIA cards.

## Commands

### `lspci`

This command runs the system `lspci` command and parse the output
in order to have it in JSON or YAML format.

```bash
$> gpu-configurator lspci --help
Like lspci but in YAML/JSON output.

Usage:
   lspci [flags]

Flags:
  -h, --help            help for lspci
  -o, --output string   Modify output format (terminal,yaml,json). (default "yaml")

Global Flags:
  -c, --config string   Gpu Configurator configfile
  -d, --debug           Enable debug output.
```

### `show`

The `show` command permits to analyze the status of the current system and
print a summary in textual mode or show the system detail in JSON or YAML format.

```bash
$> gpu-configurator show --help
Show system configuration.

Usage:
   show [flags]

Flags:
  -h, --help            help for show
  -o, --output string   Modify output format (terminal,yaml,json). (default "terminal")

Global Flags:
  -c, --config string   Gpu Configurator configfile
  -d, --debug           Enable debug output.
```

An example of the output:

```bash
$> gpu-configurator show
Copyright (c) 2024 - Macaroni OS - gpu-configurator
---------------------------------------------------------------------
Hostname:					nevyl
GPUs:						2
	- NVIDIA Corporation TU106M [GeForce RTX 2060 Mobile] [10de:1f15]
		kernel driver in use: nvidia
	- Advanced Micro Devices, Inc. [AMD/ATI] Picasso [1002:15d8]
		kernel driver in use: amdgpu

EGL External Platforms Configs Directories:
	- /usr/share/egl/egl_external_platform.d
		* 10_nvidia_wayland.json
		* 15_nvidia_gbm.json

Vulkan Layers Configs Directories:
	- /usr/share/vulkan/explicit_layer.d
		* VkLayer_khronos_validation.json
	- /usr/share/vulkan/implicit_layer.d
		* VkLayer_MESA_device_select.json
		* nvidia_layers.json

Vulkan ICD Configs Directories:
	- /usr/share/vulkan/icd.d
		* broadcom_icd.x86_64.json
		* intel_icd.x86_64.json
		* nvidia_icd.json
		* radeon_icd.x86_64.json
	- /etc/vulkan/icd.d
		* nvidia_icd.json
```
