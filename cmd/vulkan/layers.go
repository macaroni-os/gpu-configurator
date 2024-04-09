/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package vulkan

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/macaroni-os/gpu-configurator/pkg/analyzer"
	"github.com/macaroni-os/gpu-configurator/pkg/specs"

	"github.com/spf13/cobra"
)

func NewVulkanLayerCommand(config *specs.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "layers [options] layers.json",
		Short: "Enable/Disable Vulcan Layers JSON configurations.",
		PreRun: func(cmd *cobra.Command, args []string) {
			enableLayersFile, _ := cmd.Flags().GetBool("enable-layers-file")
			disableLayersFile, _ := cmd.Flags().GetBool("disable-layers-file")
			purge, _ := cmd.Flags().GetBool("purge")

			if enableLayersFile && disableLayersFile {
				fmt.Println(
					"Using both --disable-layers-file and --enable-layers-file not admitted.")
				os.Exit(1)
			}

			if purge && !disableLayersFile {
				fmt.Println(
					"--purge flag to use with --disable-layers-file.")
				os.Exit(1)
			}

			if len(args) == 0 {
				fmt.Println("Missing json filename")
				os.Exit(1)
			}

		},
		Run: func(cmd *cobra.Command, args []string) {
			enableLayersFile, _ := cmd.Flags().GetBool("enable-layers-file")
			disableLayersFile, _ := cmd.Flags().GetBool("disable-layers-file")
			purge, _ := cmd.Flags().GetBool("purge")

			jfile := args[0]

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

			vulkanldir, jsonfile := analyzer.GetSystem().GetVulkanLayerFile(jfile)
			if jsonfile == nil {
				if purge {
					// POST: ignore error if the file is not present
					return
				}
				fmt.Println("No json icd file with name", jfile, "found.")
				os.Exit(1)
			}

			if enableLayersFile {
				if !jsonfile.Disabled {
					fmt.Println("Json layer file", jfile, "already enabled.")
					return
				}

				fileabs := filepath.Join(vulkanldir.Path, jsonfile.Name)
				fileabsDisabled := fileabs + ".disabled"
				err := os.Rename(fileabsDisabled, fileabs)
				if err != nil {
					fmt.Println("Error on rename file:", err.Error())
					os.Exit(1)
				}

			} else if disableLayersFile {
				fileabs := filepath.Join(vulkanldir.Path, jsonfile.Name)

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
						fmt.Println("Json layer file", jfile, "already disabled.")
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
	flags.Bool("enable-layers-file", false, "Enable Vulkan Layers JSON file.")
	flags.Bool("disable-layers-file", false, "Disable Vulkan Layers JSON file.")
	flags.Bool("purge", false, "To use with --disable-layers-file to remove the file.")

	return cmd
}
