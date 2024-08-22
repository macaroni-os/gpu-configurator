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
	"github.com/macaroni-os/gpu-configurator/pkg/logger"
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

			log := logger.GetDefaultLogger()

			jsonLoader := args[0]

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

			eglfiles, jsonfile := analyzer.GetSystem().GetEglLoader(jsonLoader)
			if jsonfile == nil {
				if purge {
					// POST: ignore error if the file is not present
					return
				}
				log.Fatal("No json loader file with name", jsonLoader, "found.")
			}

			if enableJsonLoader {
				if !jsonfile.Disabled {
					log.InfoC("Json loader file", jsonLoader, "already enabled.")
					return
				}

				fileabs := filepath.Join(eglfiles.Path, jsonfile.Name)
				fileabsDisabled := fileabs + ".disabled"
				err := os.Rename(fileabsDisabled, fileabs)
				if err != nil {
					log.Fatal("error on rename file:", err.Error())
				}

			} else if disableJsonLoader {
				fileabs := filepath.Join(eglfiles.Path, jsonfile.Name)

				if purge {
					if jsonfile.Disabled {
						fileabs = fileabs + ".disabled"
					}

					err := os.Remove(fileabs)
					if err != nil {
						log.Fatal("error on remove file:", err.Error())
					}
				} else {

					if jsonfile.Disabled {
						log.InfoC("Json loader file", jsonLoader, "already disabled.")
						return
					}

					fileabsDisabled := fileabs + ".disabled"
					err := os.Rename(fileabs, fileabsDisabled)
					if err != nil {
						log.Fatal("error on rename file:", err.Error())
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
