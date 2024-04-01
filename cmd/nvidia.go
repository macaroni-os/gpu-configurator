/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package cmd

import (
	. "github.com/macaroni-os/gpu-configurator/cmd/nvidia"
	"github.com/macaroni-os/gpu-configurator/pkg/specs"

	"github.com/spf13/cobra"
)

func newNvidiaCommand(config *specs.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "nvidia",
		Short: "NVIDIA setup commands.",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(
		NewGbmLibCommand(config),
	)

	return cmd
}
