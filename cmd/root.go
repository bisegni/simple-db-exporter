// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// Version ...
var Version string

// Buildtime ...
var Buildtime string

// ----------------------------------------------------------------------------
// COBRA COMMAND
// ----------------------------------------------------------------------------

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "simple-db-export",
	Short: "Tools for exporting data from database(postgres)",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		showVersion, _ := cmd.Flags().GetBool("version")
		if showVersion {
			fmt.Printf("Build Time: %s\n", Buildtime)
			fmt.Printf("Commit Hash: %s\n", Version)
			os.Exit(0)
		}
	},
}

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "completion [bash|zsh|fish|powershell]",
	Long: `Generate completion script:
Bash:

$ source <(simple-db-export completion bash)

# To load completions for each session, execute once:
Linux:
	$ simple-db-export completion bash > /etc/bash_completion.d/simple-db-export
MacOS:
	$ simple-db-export completion bash > /usr/local/etc/bash_completion.d/simple-db-export

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ simple-db-export completion zsh > "${fpath[1]}/_simple-db-export"

# You will need to start a new shell for this setup to take effect.

Fish:

$ simple-db-export completion fish | source

# To load completions for each session, execute once:
$ simple-db-export completion fish > ~/.config/fish/completions/simple-db-export.fish
	`,
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}

// RootFlagsType è la struttura di RootFlags
type RootFlagsType struct {
	Verbose string
	DryRun  bool
	WorkDir string
}

// ----------------------------------------------------------------------------
// INIT
// ----------------------------------------------------------------------------

func init() {
	logrus.SetReportCaller(false)
	rootCmd.AddCommand(completionCmd)
	rootCmd.Flags().BoolP("version", "v", false, "Show version")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SetOutput(os.Stderr)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
