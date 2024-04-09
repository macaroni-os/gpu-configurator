/*
Copyright © 2024 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	specs "github.com/macaroni-os/gpu-configurator/pkg/specs"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cliName = `Copyright (c) 2024 - Macaroni OS - Daniele Rondina

gpu-configurator - A GPU configurator helper for Xwayland and/or Xorg`
)

var (
	BuildTime   string
	BuildCommit string
)

func initConfig(config *specs.Config) {
	// Set env variable
	config.Viper.SetEnvPrefix(specs.GPUCONF_ENV_PREFIX)
	config.Viper.BindEnv("config")
	config.Viper.SetDefault("config", "")

	config.Viper.AutomaticEnv()

	// Create EnvKey Replacer for handle complex structure
	replacer := strings.NewReplacer(".", "__")
	config.Viper.SetEnvKeyReplacer(replacer)

	config.Viper.SetTypeByDefaultValue(true)

}

func initCommand(rootCmd *cobra.Command, config *specs.Config) {
	var pflags = rootCmd.PersistentFlags()

	pflags.StringP("config", "c", "", "Gpu Configurator configfile")
	pflags.BoolP("debug", "d", config.Viper.GetBool("general.debug"),
		"Enable debug output.")

	config.Viper.BindPFlag("config", pflags.Lookup("config"))
	config.Viper.BindPFlag("general.debug", pflags.Lookup("debug"))

	rootCmd.AddCommand(
		newConfigCommand(config),
		newShowCommand(config),
		newLsPciCommand(config),
		newNvidiaCommand(config),
		newEglCommand(config),
		newVulkanCommand(config),
	)
}

func Execute() {
	// Create Main Instance Config object
	var config *specs.Config = specs.NewConfig(nil)

	initConfig(config)

	var rootCmd = &cobra.Command{
		Short:        cliName,
		Version:      fmt.Sprintf("%s-g%s %s", specs.GPUCONF_VERSION, BuildCommit, BuildTime),
		Args:         cobra.OnlyValidArgs,
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			var v *viper.Viper = config.Viper

			v.SetConfigType("yml")
			if v.Get("config") != "" {
				v.SetConfigFile(v.Get("config").(string))
			}

			// Parse configuration file
			err = config.Unmarshal()
			if err != nil {
				panic(err)
			}
		},
	}

	initCommand(rootCmd, config)

	// Start command execution
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
