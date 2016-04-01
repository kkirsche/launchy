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
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "View a process statuses list",
	Long: `View a list of process statuses which can be filtered using a regular
expression. If the process is active, it will also include the PID and label for
the process.`,
	Run: func(cmd *cobra.Command, args []string) {
		launchCmd := exec.Command("launchctl", "list")
		stdout, err := launchCmd.StdoutPipe()
		if err != nil {
			fmt.Printf("Could not gain stdout access to communicate with launchctl with error: %s.\n", err.Error())
			os.Exit(1)
		}

		if err = launchCmd.Start(); err != nil {
			fmt.Printf("Could not start 'launchctl list' with error: %s.\n", err.Error())
			os.Exit(1)
		}

		globalizedArgs = args
		userRegexpString := viper.GetString("regexp")
		if userRegexpString == "" && len(globalizedArgs) > 0 {
			userRegexpString = strings.Join(globalizedArgs, `\s`)
		}

		// read command's stdout line by line
		in := bufio.NewScanner(stdout)
		if in.Scan() {
			fmt.Println(in.Text()) // Headers we want to print separately if they are there
		}

		for in.Scan() {
			// Check if we need to worry about what the titles of them are
			if userRegexpString != "" {
				userRegexp, err := regexp.Compile(userRegexpString)
				if err != nil {
					fmt.Printf("Could not compile regular expression '%s' with error %s.\n", userRegexpString, err.Error())
					os.Exit(1)
				}
				matchIndex := userRegexp.FindStringIndex(in.Text())
				if matchIndex != nil {
					fmt.Println(in.Text())
				}
			} else {
				fmt.Println(in.Text())
			}
		}

		launchCmd.Wait()
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
