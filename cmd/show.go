/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/macaroni-os/gpu-configurator/pkg/analyzer"
	"github.com/macaroni-os/gpu-configurator/pkg/analyzer/pci"
	"github.com/macaroni-os/gpu-configurator/pkg/logger"
	"github.com/macaroni-os/gpu-configurator/pkg/specs"

	"github.com/spf13/cobra"
)

func printSummary(s *specs.System) error {

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "N/A"
	}

	lspci, err := pci.GetDevices()
	if err != nil {
		return err
	}

	vgaDevices := lspci.GetVGADevices()

	fmt.Println(fmt.Sprintf(
		`Copyright (c) 2024 - Macaroni OS - gpu-configurator - %s`,
		specs.GPUCONF_VERSION,
	))
	fmt.Println("---------------------------------------------------------------------")
	fmt.Println(fmt.Sprintf("Hostname:\t\t\t\t\t%s", hostname))
	fmt.Println(fmt.Sprintf("GPUs:\t\t\t\t\t\t%d", len(*vgaDevices)))
	for _, gpu := range *vgaDevices {
		fmt.Println("\t-", fmt.Sprintf("%s [%s]", gpu.Name, gpu.Id))
		if gpu.KernelDriverInUse != "" {
			fmt.Println("\t\tkernel driver in use:", gpu.KernelDriverInUse)
		}
	}
	fmt.Println("")

	if len(s.EglExtPlatformDirs) == 0 {
		fmt.Println("EGL External Platforms Configs Directories:\tNo directories available.")
	} else {
		fmt.Println("EGL External Platforms Configs Directories:")
		for idx := range s.EglExtPlatformDirs {
			fmt.Println("\t-", s.EglExtPlatformDirs[idx].Path)
			for file := range s.EglExtPlatformDirs[idx].Files {
				fmt.Println("\t\t*", file)
			}
		}
	}

	fmt.Println("")

	if len(s.VulkanLayersDirs) == 0 {
		fmt.Println("Vulkan Layers Configs Directories:\tNo directories available.")
	} else {
		fmt.Println("Vulkan Layers Configs Directories:")
		for idx := range s.VulkanLayersDirs {
			fmt.Println("\t-", s.VulkanLayersDirs[idx].Path)
			for file, f := range s.VulkanLayersDirs[idx].Files {
				if f.Disabled {
					fmt.Println("\t\t*", file, "(disabled)")
				} else {
					fmt.Println("\t\t*", file)
				}
			}
		}
	}

	fmt.Println("")

	if len(s.VulkanICDDirs) == 0 {
		fmt.Println("Vulkan ICD Configs Directories:\tNo directories available.")
	} else {
		fmt.Println("Vulkan ICD Configs Directories:")
		for idx := range s.VulkanICDDirs {
			fmt.Println("\t-", s.VulkanICDDirs[idx].Path)
			for file, f := range s.VulkanICDDirs[idx].Files {
				if f.Disabled {
					fmt.Println("\t\t*", file, "(disabled)")
				} else {
					fmt.Println("\t\t*", file)
				}
			}
		}
	}

	fmt.Println("")

	if len(s.GbmLibraries) == 0 {
		fmt.Println("GBM Backend Libraries:\tNo libraries available.")
	} else {
		fmt.Println("GBM Backend Librarires:")
		for idx := range s.GbmLibraries {
			if s.GbmLibraries[idx].Disabled {
				fmt.Println("\t-", s.GbmLibraries[idx].Name, "(disabled)")
			} else {
				fmt.Println("\t-", s.GbmLibraries[idx].Name)
			}
		}
	}

	fmt.Println("")

	if len(s.Nvidia.Drivers) == 0 {
		fmt.Println("NVIDIA Drivers:\tNo drivers available.")
	} else {
		fmt.Println("NVIDIA Drivers:")
		fmt.Println("\tActive version:", s.Nvidia.VersionActive)
		fmt.Println("\tAvailable:")
		for idx := range s.Nvidia.Drivers {
			if s.Nvidia.Drivers[idx].WithKernelModules {
				fmt.Println("\t\t-", s.Nvidia.Drivers[idx].Version, "(with kernel module)")
			} else {
				fmt.Println("\t\t-", s.Nvidia.Drivers[idx].Version)
			}
		}
		if len(s.Nvidia.KModuleAvailable) > 0 {
			fmt.Println("NVIDIA Kernel Modules Available:")
			for idx := range s.Nvidia.KModuleAvailable {
				fmt.Println(fmt.Sprintf("\t* %s - %s",
					s.Nvidia.KModuleAvailable[idx].GetFieldVersion(),
					s.Nvidia.KModuleAvailable[idx].KernelVersion,
				))
			}
		}

		if len(s.Nvidia.KModuleActive) > 0 {
			fmt.Println("NVIDIA Kernel Modules Active:")
			for idx := range s.Nvidia.KModuleActive {
				fmt.Println(fmt.Sprintf("\t* %s - %s",
					s.Nvidia.KModuleActive[idx].GetFieldVersion(),
					s.Nvidia.KModuleActive[idx].KernelVersion,
				))
			}
		}

		if len(s.Nvidia.KOpenModuleAvailable) > 0 {
			fmt.Println("NVIDIA Open Kernel Modules Available:")
			for idx := range s.Nvidia.KOpenModuleAvailable {
				fmt.Println(fmt.Sprintf("\t* %s - %s",
					s.Nvidia.KOpenModuleAvailable[idx].GetFieldVersion(),
					s.Nvidia.KOpenModuleAvailable[idx].KernelVersion,
				))
			}
		}

		if len(s.Nvidia.KOpenModuleActive) > 0 {
			fmt.Println("NVIDIA Open Kernel Modules Active:")
			for idx := range s.Nvidia.KOpenModuleActive {
				fmt.Println(fmt.Sprintf("\t* %s - %s",
					s.Nvidia.KOpenModuleActive[idx].GetFieldVersion(),
					s.Nvidia.KOpenModuleActive[idx].KernelVersion,
				))
			}
		}
	}

	return nil
}

func newShowCommand(config *specs.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "show",
		Short: "Show system configuration.",
		Args:  cobra.NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			output, _ := cmd.Flags().GetString("output")
			switch output {
			case "", "terminal", "json", "yaml":
			default:
				fmt.Println(fmt.Sprintf("Invalid value %s for output.",
					output,
				))
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			output, _ := cmd.Flags().GetString("output")
			log := logger.GetDefaultLogger()

			analyzer, err := analyzer.NewAnalyzer(
				config.GetGeneral().GetBackendType(),
			)
			if err != nil {
				log.Fatal("error on initialize analyzer:", err.Error())
			}

			err = analyzer.Read()
			if err != nil {
				log.Fatal("error on analyze system:", err.Error())
			}

			if output == "terminal" {
				err := printSummary(analyzer.GetSystem())
				if err != nil {
					log.Fatal("Error", err.Error())
				}

			} else {
				var data []byte

				switch output {
				case "json":
					data, err = analyzer.GetSystem().Json()
				default:
					data, err = analyzer.GetSystem().Yaml()
				}

				if err != nil {
					log.Fatal("error on convert system on", output, err.Error())
				}

				fmt.Println(string(data))
			}

		},
	}

	var flags = cmd.Flags()
	flags.StringP("output", "o", "terminal",
		"Modify output format (terminal,yaml,json).")

	return cmd
}
