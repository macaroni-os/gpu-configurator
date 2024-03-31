/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package cmd

import (
	"fmt"
	"os"

	specs "github.com/macaroni-os/gpu-configurator/pkg/specs"

	"github.com/spf13/cobra"
)

func newConfigCommand(config *specs.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "config",
		Short: "Show configuration params",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			jsonOutput, _ := cmd.Flags().GetBool("json")

			var err error
			var data []byte

			if jsonOutput {
				data, err = config.Json()
			} else {
				data, err = config.Yaml()
			}
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fmt.Println(string(data))

		},
	}

	var flags = cmd.Flags()
	flags.Bool("json", false, "JSON output")

	return cmd
}
