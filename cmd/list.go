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
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List services which can be controlled using launchy or launchctl",
	Long: `List services which can be controlled using launchy or launchtl. This
command is the equivalent command ls followed by:
	- ~/Library/LaunchAgents
	- /Library/LaunchDaemons (root only)
	- ~/Library/LaunchAgents
	- ~/Library/LaunchDaemons (root only)`,
	Run: func(cmd *cobra.Command, args []string) {
		user, err := user.Current()
		if err != nil {
			fmt.Println(err.Error())
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

		globalizedArgs = args
		for _, path := range servicesPaths {
			verbosePrintf("# Walking path %s", path)
			filepath.Walk(path, walkPathAction)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// // and all subcommands, e.g.:
	// listCmd.PersistentFlags().StringP("regexp", "r", "", "A regex to match filenames against")
	// viper.BindPFlag("regexp", listCmd.PersistentFlags().Lookup("regexp"))
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func walkPathAction(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	userRegexpString := viper.GetString("regexp")
	if userRegexpString == "" && len(globalizedArgs) > 0 {
		userRegexpString = strings.Join(globalizedArgs, `\s`)
	}

	// Make sure we only grab plist files
	if hasPlistExtension(path) {
		nameLength := len(info.Name())
		// Check if we need to worry about what the titles of them are
		if hasRegexpProvided(userRegexpString) {
			userRegexp, err := regexp.Compile(userRegexpString)
			if err != nil {
				return fmt.Errorf("Could not compile regular expression: `%s` due to an error: %s\n", userRegexpString, err.Error())
			}
			matchIndex := userRegexp.FindStringIndex(info.Name())
			if matchIndex != nil {
				fmt.Println(info.Name()[:nameLength-6])
			}
		} else {
			fmt.Println(info.Name()[:nameLength-6])
		}
	}
	return nil
}

func hasPlistExtension(path string) bool {
	if filepath.Ext(path) == ".plist" {
		return true
	}
	return false
}

func hasRegexpProvided(provided string) bool {
	if provided != "" {
		return true
	}
	return false
}
