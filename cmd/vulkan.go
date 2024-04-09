/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package cmd

import (
	. "github.com/macaroni-os/gpu-configurator/cmd/vulkan"
	"github.com/macaroni-os/gpu-configurator/pkg/specs"

	"github.com/spf13/cobra"
)

func newVulkanCommand(config *specs.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "vulkan",
		Short: "Vulkan setup commands.",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(
		NewVulkanIcdCommand(config),
		NewVulkanLayerCommand(config),
	)

	return cmd
}
