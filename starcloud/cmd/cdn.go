/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package cmd

import (
	"github.com/gzuidhof/starcloud/starcloud/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cdnCmd = &cobra.Command{
	Use:   "cdn",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		server.StartCDNApp()
	},
	PreRun: bindCDNCmdToViper,
}

func init() {
	rootCmd.AddCommand(cdnCmd)
	cdnCmd.Flags().IntP("port", "p", 8080, "Port to listen on")
	cdnCmd.Flags().String("cache_folder", "starcloud_cache", "Folder path to use for caching files into")
}

func bindCDNCmdToViper(cmd *cobra.Command, args []string) {
	viper.BindPFlag("port", cmd.Flags().Lookup("port"))
	viper.BindPFlag("cache_folder", cmd.Flags().Lookup("cache_folder"))
}
