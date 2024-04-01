/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/macaroni-os/gpu-configurator/pkg/analyzer"
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

			if enableDriver && disableDriver {
				fmt.Println(
					"Using both --enable-driver and --disable-driver not admitted.")
				os.Exit(1)
			}

			if purge && !disableDriver {
				fmt.Println(
					"--purge flag to use with --disable-driver.")
				os.Exit(1)
			}

		},
		Run: func(cmd *cobra.Command, args []string) {
			enableDriver, _ := cmd.Flags().GetBool("enable-driver")
			disableDriver, _ := cmd.Flags().GetBool("disable-driver")
			purge, _ := cmd.Flags().GetBool("purge")

			analyzer, err := analyzer.NewAnalyzer(
				config.GetGeneral().GetBackendType(),
			)
			if err != nil {
				fmt.Println("ERROR", err.Error())
				os.Exit(1)
			}

			err = analyzer.Read()
			if err != nil {
				fmt.Println("Error on analyze system", err.Error())
				os.Exit(1)
			}

			libName := "nvidia-drm_gbm.so"

			nvidiagbmlib := analyzer.GetSystem().GetGBMLibrary(
				libName,
			)
			if analyzer.GetSystem().Nvidia == nil ||
				analyzer.GetSystem().Nvidia.VersionActive == "" {
				fmt.Println("No NVIDIA version active available.")
				fmt.Println("Check if the package with NVIDIA drivers is installed.")
				os.Exit(1)
			}

			if enableDriver {
				if nvidiagbmlib == nil {
					// POST: The library link is not present.

					nvidiaDriver := analyzer.GetSystem().Nvidia.GetDriver(
						analyzer.GetSystem().Nvidia.VersionActive,
					)
					if nvidiaDriver == nil {
						fmt.Println("Unexpected error on retrieve nvidia driver data.")
						os.Exit(1)
					}

					linkedFile := filepath.Join(
						nvidiaDriver.Path, "lib64", libName,
					)
					linkFile := filepath.Join(
						analyzer.GetBackend().GetGBMLibDir(),
						libName,
					)

					err := os.Symlink(linkedFile, linkFile)
					if err != nil {
						fmt.Println(fmt.Sprintf(
							"error on create symlink on %s: %s",
							linkFile, err.Error()))
						os.Exit(1)
					}

				} else {

					if !nvidiagbmlib.Disabled {
						fmt.Println("Library", libName, "already active.")
						fmt.Println("Nothing to do.")
						return
					}

					libpath := filepath.Join(
						analyzer.GetBackend().GetGBMLibDir(),
						nvidiagbmlib.Name,
					)
					libpathDisabled := libpath + ".disabled"

					err := os.Rename(libpathDisabled, libpath)
					if err != nil {
						fmt.Println("Error on rename link:", err.Error())
						os.Exit(1)
					}

				}

				fmt.Println("Operation done.")

			} else if disableDriver {

				libpath := filepath.Join(
					analyzer.GetBackend().GetGBMLibDir(),
					nvidiagbmlib.Name,
				)
				libpathDisabled := libpath + ".disabled"

				if purge && nvidiagbmlib != nil {
					if nvidiagbmlib.Disabled {
						err = os.Remove(libpathDisabled)
					} else {
						err = os.Remove(libpath)
					}

					if err != nil {
						fmt.Println("Error on remove link:", err.Error())
						os.Exit(1)
					}

				} else if nvidiagbmlib == nil || nvidiagbmlib.Disabled {
					fmt.Println("Library", libName, "not present or already disable.")
					fmt.Println("Nothing to do.")
					return
				} else {

					err := os.Rename(libpath, libpathDisabled)
					if err != nil {
						fmt.Println("Error on rename link:", err.Error())
						os.Exit(1)
					}
				}

				fmt.Println("Operation done.")
			}

		},
	}

	var flags = cmd.Flags()
	flags.Bool("enable-driver", false, "Enable NVIDIA GBM library.")
	flags.Bool("disable-driver", false, "Disable NVIDIA GBM library.")
	flags.Bool("purge", false, "To use with --disable-driver to remove the link library.")

	return cmd
}
