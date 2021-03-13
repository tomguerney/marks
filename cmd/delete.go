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
	"github.com/tomguerney/marks/arg"
	"github.com/tomguerney/marks/colorizer"
	"github.com/tomguerney/marks/config"
	"github.com/tomguerney/marks/io"
	"github.com/tomguerney/marks/printer"
	"github.com/tomguerney/marks/prompter"
	"github.com/tomguerney/marks/runner"
	"github.com/tomguerney/marks/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete id [tags...]",
	Short: "Delete a bookmark",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runDelete,
}

func runDelete(cmd *cobra.Command, argv []string) error {
	args, err := combineDeleteArgs(cmd.Flags(), argv)
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
	prompter := prompter.NewPrompter()
	runner := runner.NewDeleteRunner(args, config, markService, printer, prompter)
	if err := runner.Run(); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringP("url", "u", "", "(can be partial) --url abc.net.au")
	deleteCmd.Flags().StringSliceP("tag", "t", []string{}, "--tag news --tag \"current affairs\"")
}

func combineDeleteArgs(flagSet *pflag.FlagSet, argv []string) (*runner.DeleteArgs, error) {

	parser := arg.NewParser(argv)

	id, err := parser.Pop()
	if err != nil {
		return nil, err
	}

	flagTags := parser.Remaining()

	url, err := flagSet.GetString("url")
	if err != nil {
		return nil, err
	}

	tags, err := flagSet.GetStringSlice("tag")
	if err != nil {
		return nil, err
	}

	tags = append(tags, flagTags...)

	return runner.NewDeleteArgs(id, url, tags), nil
}
