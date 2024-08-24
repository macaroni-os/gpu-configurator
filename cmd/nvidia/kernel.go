/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package nvidia

import (
	"fmt"
	"os"

	"github.com/macaroni-os/gpu-configurator/pkg/analyzer"
	"github.com/macaroni-os/gpu-configurator/pkg/logger"
	"github.com/macaroni-os/gpu-configurator/pkg/specs"

	"github.com/spf13/cobra"
)

func NewKernlCommand(config *specs.Config) *cobra.Command {

	var cmd = &cobra.Command{
		Use:     "kernel [nvidia-version|*] [[kernel-version]]",
		Short:   "Activate or disable a specific kernel version of NVIDIA driver.",
		Aliases: []string{"k"},
		PreRun: func(cmd *cobra.Command, args []string) {
			log := logger.GetDefaultLogger()
			purge, _ := cmd.Flags().GetBool("purge")
			if len(args) == 0 {
				log.InfoC("Missing nvidia driver version argument.")
				os.Exit(1)
			}
			targetVersion := args[0]
			kernelVersion := ""
			if len(args) > 1 {
				kernelVersion = args[1]
			}

			if !purge && targetVersion == "*" {
				log.InfoC("nvidia-version * could be used only with --purge")
				os.Exit(1)
			}
			if !purge && kernelVersion == "" {
				log.InfoC("kernel-version mandatory without --purge")
				os.Exit(1)
			}

			if targetVersion == "" {
				log.InfoC("Invalid nvidia-version field")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.GetDefaultLogger()
			purge, _ := cmd.Flags().GetBool("purge")
			proprietary, _ := cmd.Flags().GetBool("proprietary")

			analyzer, err := analyzer.NewAnalyzer(
				config.GetGeneral().GetBackendType(),
			)

			err = analyzer.Read()
			if err != nil {
				log.Fatal("Error on analyze system", err.Error())
			}

			kernelVersion := ""
			targetVersion := args[0]

			if len(args) > 1 {
				kernelVersion = args[1]
			}

			if purge {
				// POST: Purge

				err := analyzer.GetBackend().PurgeNVIDIAKernelDriverActive(
					analyzer.GetSystem().GetNvidia(),
					targetVersion, kernelVersion)
				if err != nil {
					log.Fatal(err.Error())
				}

			} else {

				// Check if the kernel driver is present.
				if !analyzer.GetSystem().GetNvidia().HasVersion(targetVersion) {
					log.Fatal(
						fmt.Sprintf("Driver with version %s not found.",
							targetVersion))
				}

				err := analyzer.GetBackend().ActiveNVIDIAKernelDriver(
					analyzer.GetSystem().GetNvidia(),
					targetVersion, kernelVersion, !proprietary,
				)
				if err != nil {
					log.Fatal(err.Error())
				}

			}

			log.InfoC("All done.")
		},
	}

	var flags = cmd.Flags()
	flags.Bool("purge", false, "remove the selected kernels from the /lib/modules/<kernel>/video directory.")
	flags.Bool("proprietary", true, "Use NVIDIA proprietary driver instead of Open NVIDIA driver.")

	return cmd
}
