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

func NewConfigureCommand(config *specs.Config) *cobra.Command {

	var cmd = &cobra.Command{
		Use:     "configure [version]",
		Short:   "Configure a specific version of NVIDIA driver.",
		Aliases: []string{"c", "conf", "set"},
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println(
					"Missing nvidia driver version argument.")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			force, _ := cmd.Flags().GetBool("force")
			purge, _ := cmd.Flags().GetBool("purge")
			withVideoGroup, _ := cmd.Flags().GetBool("with-video-group")

			log := logger.GetDefaultLogger()

			analyzer, err := analyzer.NewAnalyzer(
				config.GetGeneral().GetBackendType(),
			)

			err = analyzer.Read()
			if err != nil {
				log.Fatal("Error on analyze system", err.Error())
			}

			targetVersion := args[0]

			if purge {
				if analyzer.GetSystem().GetNvidia().VersionActive == targetVersion || force {
					err = analyzer.GetBackend().PurgeNVIDIADriver(
						analyzer.GetSystem().GetNvidia(),
					)
					if err != nil {
						log.Fatal(fmt.Sprintf(
							"Error on purge active version %s: %s",
							targetVersion, err.Error()))
					}
					log.InfoC(fmt.Sprintf("Nvidia generated files purged."))
				}

			} else {
				log.InfoC(fmt.Sprintf("Setting version %s...", targetVersion))

				if analyzer.GetSystem().GetNvidia().VersionActive != targetVersion || force {

					if analyzer.GetSystem().GetNvidia().VersionActive == targetVersion {
						err = analyzer.GetBackend().PurgeNVIDIADriver(
							analyzer.GetSystem().GetNvidia(),
						)
						if err != nil {
							log.Fatal(fmt.Sprintf(
								"Error on purge active version %s: %s",
								targetVersion, err.Error()))
						}
					}

					err = analyzer.GetBackend().SetNVIDIAVersion(
						analyzer.GetSystem().GetNvidia(),
						targetVersion)

					if err != nil {
						log.Fatal(fmt.Sprintf(
							"Error on set version %s: %s",
							targetVersion,
							err.Error()))
					}

					// Setting modprobe
					err = analyzer.GetBackend().SetNVIDIAModprobeFiles(
						analyzer.GetSystem().GetNvidia(),
						withVideoGroup, force,
					)

					if err != nil {
						log.Fatal(fmt.Sprintf(
							"Error on setup modprobe files: %s",
							targetVersion,
							err.Error()))
					}
				}
			}
		},
	}

	var flags = cmd.Flags()
	flags.BoolP("force", "f", false, "Forcing set of the selected version.")
	flags.Bool("purge", false, "Remove generated file of the the selected version.")
	flags.Bool("with-video-group", true,
		"Set the kernel NVIDIA driver with the video group id")

	return cmd
}
