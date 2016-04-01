// Copyright © 2016 Kevin Kirsche <kev.kirsche@gmail.com>
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
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loadMatches []string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		user, err := user.Current()
		if err != nil {
			stdPrintf(err.Error())
			os.Exit(1)
		}

		servicesPaths := []string{
			"/Library/LaunchAgents",
			user.HomeDir + "/Library/LaunchAgents",
		}

		if user.Username == "root" {
			servicesPaths = append(servicesPaths, "/Library/LaunchDaemons")
			servicesPaths = append(servicesPaths, user.HomeDir+"/Library/LaunchDaemons")
		}

		launchCmd := exec.Command("launchctl", "load")

		if force {
			launchCmd.Args = append(launchCmd.Args, "-F")
		}

		if enable {
			launchCmd.Args = append(launchCmd.Args, "-w")
		}

		globalizedArgs = args
		for _, path := range servicesPaths {
			verbosePrintf("Walking path %s", path)
			filepath.Walk(path, launchDuringWalk)
		}

		if len(loadMatches) > 1 {
			stdPrintf("More than one service matched. Cannot start multiple services. Exiting...")
			os.Exit(1)
		}

		launchCmd.Args = append(launchCmd.Args, loadMatches...)

		err = launchCmd.Run()
		if err != nil {
			stdPrintf("Could not start services: %s\n due to error: %s", strings.Join(loadMatches, "\n"), err.Error())
			os.Exit(1)
		}

		stdPrintf("Service %s started correctly.\n", filepath.Base(strings.Join(loadMatches, "")))
	},
}

func init() {
	RootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	startCmd.Flags().BoolVarP(&force, "force", "f", false, "Forcibly start the service.")
	viper.BindPFlag("force", startCmd.Flags().Lookup("force"))

	startCmd.Flags().BoolVarP(&enable, "enable", "e", false, "If the service is disabled, it will be enabled.")
	viper.BindPFlag("enable", startCmd.Flags().Lookup("enable"))
}

func launchDuringWalk(path string, info os.FileInfo, err error) error {
	userRegexpString := viper.GetString("regexp")
	if userRegexpString == "" && len(globalizedArgs) > 0 {
		userRegexpString = strings.Join(globalizedArgs, `\s`)
	}

	userRegexp, err := regexp.Compile(userRegexpString)
	if err != nil {
		return fmt.Errorf("Could not compile regular expression '%s' with error %s.\n", userRegexpString, err.Error())
	}

	match := userRegexp.FindStringIndex(path)
	if match != nil {
		verbosePrintf("Found matching service: %s.\n", path)
		loadMatches = append(loadMatches, path)
	}

	return nil
}
