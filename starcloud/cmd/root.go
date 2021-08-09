/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmd

import (
	"fmt"
	"os"

	"github.com/gzuidhof/starcloud/starcloud/config"
	"github.com/spf13/cobra"
)

var cfgFile string

// Version (injected by goreleaser)
var Version = "<unknown>"
var Date = "date unknown"
var Commit = ""
var Target = ""

var rootCmd = &cobra.Command{
	Use:   "starcloud",
	Short: "Starcloud is the CDN (origin) service for Starboard assets",
	Long:  `Starcloud is the CDN (origin) service for Starboard assets`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() { config.InitConfig(cfgFile) })
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.starcloud.yaml)")

	rootCmd.Version = Version + " " + Target + " (" + Date + ") " + Commit
}
