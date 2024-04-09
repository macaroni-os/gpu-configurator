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

func newEglCommand(config *specs.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "egl [options] eglloader.json",
		Short: "Enable/Disable EGL JSON configurations.",
		PreRun: func(cmd *cobra.Command, args []string) {
			enableJsonLoader, _ := cmd.Flags().GetBool("enable-json-loader")
			disableJsonLoader, _ := cmd.Flags().GetBool("disable-json-loader")
			purge, _ := cmd.Flags().GetBool("purge")

			if enableJsonLoader && disableJsonLoader {
				fmt.Println(
					"Using both --disable-json-loader and --enable-json-loader not admitted.")
				os.Exit(1)
			}

			if purge && !disableJsonLoader {
				fmt.Println(
					"--purge flag to use with --disable-json-loader.")
				os.Exit(1)
			}

			if len(args) == 0 {
				fmt.Println("Missing json filename")
				os.Exit(1)
			}

		},
		Run: func(cmd *cobra.Command, args []string) {
			enableJsonLoader, _ := cmd.Flags().GetBool("enable-json-loader")
			disableJsonLoader, _ := cmd.Flags().GetBool("disable-json-loader")
			purge, _ := cmd.Flags().GetBool("purge")

			jsonLoader := args[0]

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

			eglfiles, jsonfile := analyzer.GetSystem().GetEglLoader(jsonLoader)
			if jsonfile == nil {
				if purge {
					// POST: ignore error if the file is not present
					return
				}
				fmt.Println("No json loader file with name", jsonLoader, "found.")
				os.Exit(1)
			}

			if enableJsonLoader {
				if !jsonfile.Disabled {
					fmt.Println("Json loader file", jsonLoader, "already enabled.")
					return
				}

				fileabs := filepath.Join(eglfiles.Path, jsonfile.Name)
				fileabsDisabled := fileabs + ".disabled"
				err := os.Rename(fileabsDisabled, fileabs)
				if err != nil {
					fmt.Println("Error on rename file:", err.Error())
					os.Exit(1)
				}

			} else if disableJsonLoader {
				fileabs := filepath.Join(eglfiles.Path, jsonfile.Name)

				if purge {
					if jsonfile.Disabled {
						fileabs = fileabs + ".disabled"
					}

					err := os.Remove(fileabs)
					if err != nil {
						fmt.Println("Error on remove file:", err.Error())
						os.Exit(1)
					}
				} else {

					if jsonfile.Disabled {
						fmt.Println("Json loader file", jsonLoader, "already disabled.")
						return
					}

					fileabsDisabled := fileabs + ".disabled"
					err := os.Rename(fileabs, fileabsDisabled)
					if err != nil {
						fmt.Println("Error on rename file:", err.Error())
						os.Exit(1)
					}
				}
			}
		},
	}

	var flags = cmd.Flags()
	flags.Bool("enable-json-loader", false, "Enable EGL JSON loader.")
	flags.Bool("disable-json-loader", false, "Disable EGL JSON loader.")
	flags.Bool("purge", false, "To use with --disable-json-loader to remove the JSON file.")

	return cmd
}
