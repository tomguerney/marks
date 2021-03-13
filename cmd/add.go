/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/abitofoldtom/marks/arg"
	"github.com/abitofoldtom/marks/colorizer"
	"github.com/abitofoldtom/marks/config"
	"github.com/abitofoldtom/marks/io"
	"github.com/abitofoldtom/marks/printer"
	"github.com/abitofoldtom/marks/runner"
	"github.com/abitofoldtom/marks/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add id",
	Short: "Add a bookmark",
	Args:  cobra.ExactArgs(1),
	RunE:  runAdd,
}

func runAdd(cmd *cobra.Command, argv []string) error {
	args, err := combineAddArgs(cmd.Flags(), argv)
	if err != nil {
		return err
	}
	config, err := config.NewLoader(viper.GetViper()).Load()
	if err != nil {
		return err
	}
	markService := yaml.NewMarkService(config, io.NewReaderWriter())
	colorizer := colorizer.NewColorizer()
	printer := printer.NewPrinter(config, colorizer)
	runner := runner.NewAddRunner(args, config, markService, printer)
	if err := runner.Run(); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("url", "u", "", "--url https://www.abc.net.au/news/")
	addCmd.MarkFlagRequired("url")
	addCmd.Flags().StringSliceP("tag", "t", []string{}, "--tag news --tag \"current affairs\"")
}

func combineAddArgs(flagSet *pflag.FlagSet, argv []string) (*runner.AddArgs, error) {

	parser := arg.NewParser(argv)

	id, err := parser.Pop()
	if err != nil {
		return nil, err
	}

	url, err := flagSet.GetString("url")
	if err != nil {
		return nil, err
	}

	tags, err := flagSet.GetStringSlice("tag")
	if err != nil {
		return nil, err
	}

	return runner.NewAddArgs(id, url, tags), nil
}
