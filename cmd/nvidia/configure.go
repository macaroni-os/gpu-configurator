/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package nvidia

import (
	"fmt"
	"os"

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

		},
	}

	return cmd
}
