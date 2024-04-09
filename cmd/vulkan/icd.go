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

func NewVulkanIcdCommand(config *specs.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "icd [options] icd.json",
		Short: "Enable/Disable Vulcan ICD JSON configurations.",
		PreRun: func(cmd *cobra.Command, args []string) {
			enableIcdFile, _ := cmd.Flags().GetBool("enable-icd-file")
			disableIcdFile, _ := cmd.Flags().GetBool("disable-icd-file")
			purge, _ := cmd.Flags().GetBool("purge")

			if enableIcdFile && disableIcdFile {
				fmt.Println(
					"Using both --disable-icd-file and --enable-icd-file not admitted.")
				os.Exit(1)
			}

			if purge && !disableIcdFile {
				fmt.Println(
					"--purge flag to use with --disable-icd-file.")
				os.Exit(1)
			}

			if len(args) == 0 {
				fmt.Println("Missing json filename")
				os.Exit(1)
			}

		},
		Run: func(cmd *cobra.Command, args []string) {
			enableIcdFile, _ := cmd.Flags().GetBool("enable-icd-file")
			disableIcdFile, _ := cmd.Flags().GetBool("disable-icd-file")
			purge, _ := cmd.Flags().GetBool("purge")

			icdfile := args[0]

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

			icdfiles, jsonfile := analyzer.GetSystem().GetVulkanIcdFile(icdfile)
			if jsonfile == nil {
				if purge {
					// POST: ignore error if the file is not present
					return
				}
				fmt.Println("No json icd file with name", icdfile, "found.")
				os.Exit(1)
			}

			if enableIcdFile {
				if !jsonfile.Disabled {
					fmt.Println("Json icd file", icdfile, "already enabled.")
					return
				}

				fileabs := filepath.Join(icdfiles.Path, jsonfile.Name)
				fileabsDisabled := fileabs + ".disabled"
				err := os.Rename(fileabsDisabled, fileabs)
				if err != nil {
					fmt.Println("Error on rename file:", err.Error())
					os.Exit(1)
				}

			} else if disableIcdFile {
				fileabs := filepath.Join(icdfiles.Path, jsonfile.Name)

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
						fmt.Println("Json loader file", icdfile, "already disabled.")
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
	flags.Bool("enable-icd-file", false, "Enable ICD JSON file.")
	flags.Bool("disable-icd-file", false, "Disable ICD JSON file.")
	flags.Bool("purge", false, "To use with --disable-icd-file to remove the ICD file.")

	return cmd
}
