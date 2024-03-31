/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/macaroni-os/gpu-configurator/pkg/analyzer/pci"
	"github.com/macaroni-os/gpu-configurator/pkg/specs"

	"github.com/spf13/cobra"
)

func newLsPciCommand(config *specs.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "lspci",
		Short: "Like lspci but in YAML/JSON output.",
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

			lspci, err := pci.GetDevices()
			if err != nil {
				fmt.Println("Error on read pci data:", err.Error())
				os.Exit(1)
			}

			if output == "terminal" {
				fmt.Println("Not yet implemented")
			} else {
				var data []byte

				switch output {
				case "json":
					data, err = lspci.Json()
				default:
					data, err = lspci.Yaml()
				}

				if err != nil {
					fmt.Println("Error on convert data", output, err.Error())
					os.Exit(1)
				}

				fmt.Println(string(data))
			}

		},
	}

	var flags = cmd.Flags()
	flags.StringP("output", "o", "yaml",
		"Modify output format (terminal,yaml,json).")

	return cmd
}
