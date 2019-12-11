// Copyright © 2019 VMware, INC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/edgexfoundry-holding/edgex-cli/cmd/db"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/device"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/deviceservice"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/event"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/interval"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/notification"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/profile"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/reading"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/status"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/subscription"
	"github.com/edgexfoundry-holding/edgex-cli/cmd/version"
	"github.com/edgexfoundry-holding/edgex-cli/config"
	"github.com/edgexfoundry-holding/edgex-cli/pkg/pager"
)

// NewCommand returns rootCmd which represents the base command when called without any subcommands
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// set flags
			noPager, err := cmd.Flags().GetBool("no-pager")

			verbose, _ := cmd.Flags().GetBool("verbose")
			if err != nil {
				fmt.Println("couldn't get verbose flag")
			}
			viper.Set("verbose", verbose)
			if verbose {
				noPager = true
			}

			url, _ := cmd.Flags().GetBool("url")
			if err != nil {
				fmt.Println("couldn't get url flag")
			}
			viper.Set("url", url)

			if err != nil {
				fmt.Println("couldn't get no-pager flag")
			}
			viper.Set("writer", os.Stdout)
			if !noPager {
				w, err := pager.NewWriter()
				if err == nil {
					viper.Set("writer", w)
					viper.Set("writerShouldClose", true) // This flag prevents us from calling close on stdout
				}
			}

		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			shouldClose := viper.GetBool("writerShouldClose")
			if shouldClose {
				pw := viper.Get("writer").(io.Closer)
				if pw != os.Stdout {
					err := pw.Close()
					if err != nil {
						_ = fmt.Errorf(err.Error())
					}
				}
			}
		},
		Use:   "edgex",
		Short: "EdgeX command line interface",
		Long: " _____    _           __  __  _____                     _            \n" +
			"| ____|__| | __ _  ___\\ \\/ / |  ___|__  _   _ _ __   __| |_ __ _   _ \n" +
			"|  _| / _` |/ _` |/ _ \\\\  /  | |_ / _ \\| | | | '_ \\ / _` | '__| | | |\n" +
			"| |__| (_| | (_| |  __//  \\  |  _| (_) | |_| | | | | (_| | |  | |_| |\n" +
			"|_____\\__,_|\\__, |\\___/_/\\_\\ |_|  \\___/ \\__,_|_| |_|\\__,_|_|   \\__, |\n" +
			"            |___/                                              |___/ \n" +
			`

https://www.edgexfoundry.org/
	`,
	}

	// Add all subcommands below:
	cmd.AddCommand(device.NewCommand())
	cmd.AddCommand(deviceservice.NewCommand())
	cmd.AddCommand(profile.NewCommand())
	cmd.AddCommand(event.NewCommand())
	cmd.AddCommand(reading.NewCommand())
	cmd.AddCommand(status.NewCommand())
	cmd.AddCommand(db.NewCommand())
	// --- Support Services Commands ---
	cmd.AddCommand(notification.NewCommand())
	cmd.AddCommand(subscription.NewCommand())
	cmd.AddCommand(interval.NewCommand())
	cmd.AddCommand(version.NewCommand())

	// global flags
	Verbose := false
	URL := false
	NoPager := false

	// get flags values
	cmd.PersistentFlags().BoolVarP(&URL, "url", "u", false, "Print URL(s) used by the entered command.")
	cmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Print entire HTTP response.")
	cmd.PersistentFlags().BoolVarP(&NoPager, "no-pager", "", false, "Do not pipe output into a pager.")

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// set default config
	configDirPath := os.Getenv("HOME") + "/.edgex-cli/"
	configFilePath := "config.yaml"
	env := config.NewViperEnv()
	err := config.SetConfig(env, configDirPath, configFilePath)

	if err != nil {
		log.Fatal(err)
	}

	if err := NewCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
