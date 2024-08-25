/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package nvidia

import (
	"os"

	"github.com/macaroni-os/gpu-configurator/pkg/analyzer"
	"github.com/macaroni-os/gpu-configurator/pkg/logger"
	"github.com/macaroni-os/gpu-configurator/pkg/specs"

	"github.com/spf13/cobra"
)

func NewGbmLibCommand(config *specs.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "gbmlib",
		Short: "GBM Backend Library configuration.",
		PreRun: func(cmd *cobra.Command, args []string) {
			enableDriver, _ := cmd.Flags().GetBool("enable-driver")
			disableDriver, _ := cmd.Flags().GetBool("disable-driver")
			purge, _ := cmd.Flags().GetBool("purge")
			log := logger.GetDefaultLogger()

			if enableDriver && disableDriver {
				log.InfoC(
					"Using both --enable-driver and --disable-driver not admitted.")
				os.Exit(1)
			}

			if purge && !disableDriver {
				log.InfoC(
					"--purge flag to use with --disable-driver.")
				os.Exit(1)
			}

		},
		Run: func(cmd *cobra.Command, args []string) {
			enableDriver, _ := cmd.Flags().GetBool("enable-driver")
			disableDriver, _ := cmd.Flags().GetBool("disable-driver")
			purge, _ := cmd.Flags().GetBool("purge")
			log := logger.GetDefaultLogger()

			analyzer, err := analyzer.NewAnalyzer(
				config.GetGeneral().GetBackendType(),
			)
			if err != nil {
				log.Fatal(err.Error())
			}

			err = analyzer.Read()
			if err != nil {
				log.Fatal("error on analyze system", err.Error())
			}

			if analyzer.GetSystem().Nvidia == nil ||
				analyzer.GetSystem().Nvidia.VersionActive == "" {
				log.Warning("No NVIDIA version active available.")
				log.Warning("Check if the package with NVIDIA drivers is installed.")
				os.Exit(1)
			}

			if purge {

				err = analyzer.GetBackend().PurgeGBMLinks(
					analyzer.GetSystem(),
					[]string{"nvidia-drm_gbm.so", "nvidia_gbm.so"},
				)

			} else {

				enabled := enableDriver
				if disableDriver {
					enabled = false
				}
				err = analyzer.GetBackend().ConfigureGBMLinks(
					analyzer.GetSystem(),
					analyzer.GetSystem().Nvidia.VersionActive,
					enabled,
				)

			}

			if err != nil {
				log.Fatal(err.Error())
			}

			log.InfoC("Operation done.")
		},
	}

	var flags = cmd.Flags()
	flags.Bool("enable-driver", false, "Enable NVIDIA GBM library.")
	flags.Bool("disable-driver", false, "Disable NVIDIA GBM library.")
	flags.Bool("purge", false, "To use with --disable-driver to remove the link library.")

	return cmd
}
