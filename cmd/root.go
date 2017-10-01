// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"regexp"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gh",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		// Inside work tree?
		inside, err := exec.Command("git",
																"rev-parse",
																"--is-inside-work-tree").Output()
		if err != nil || strings.TrimRight(string(inside), "\n") != "true" {
			fmt.Println(err)
			fmt.Println(string(inside))
			fmt.Println("Not a git repository.")
			fmt.Println("Please move to inside work tree.")
			os.Exit(0)
		}

		// Get remote names
		names, err := exec.Command("git",
															 "remote").Output()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		names_str := strings.TrimRight(string(names), "\n")
		name_arr := strings.Split(names_str, "\n")

		// if length of name_arr == 1
		length_names := len(name_arr)
		remote_name := ""
		if length_names != 1 {

			// TODO: ひとまずは、upstream をいれる..
			// あとで、標準入力から選択できるようにしたい
			// そのデータをどう保存していいかわからない
			if len(args) == 1 {
				remote_name = args[0]
			} else {
				remote_name = "upstream"
			}
		} else {
			remote_name = name_arr[0]
		}

		url_byte, err := exec.Command("git",
																	"config",
																	"--get",
																	"remote." + remote_name + ".url").Output()
		if err != nil {
			fmt.Println("remote name '" + remote_name + "' does not exist.")
			fmt.Println("Usage:")
			fmt.Println("gh {remote_name}")
			// fmt.Println(err)
			os.Exit(0)
		}

		// Replace [git@github.com:] to [https://github.com/]
		url := string(url_byte)
		r := regexp.MustCompile(`git@github.com:`)
		if r.MatchString(url) {
			url = r.ReplaceAllString(url, "https://github.com/")
		}

		// Open url
		fmt.Println(url)
		exec.Command("open", url).Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Printf("hello error\n")

		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gh.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gh" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gh")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
