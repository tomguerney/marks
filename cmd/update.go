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
	"github.com/abitofoldtom/marks/prompter"
	"github.com/abitofoldtom/marks/runner"
	"github.com/abitofoldtom/marks/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update id [tags...]",
	Short: "Update a bookmark",
	RunE:  runUpdate,
}

func runUpdate(cmd *cobra.Command, argv []string) error {
	args, err := combineUpdateArgs(cmd.Flags(), argv)
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
	runner := runner.NewUpdateRunner(args, config, markService, printer, prompter)
	if err := runner.Run(); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringP("url", "u", "", "(can be partial) --url abc.net.au")
	updateCmd.Flags().StringSliceP("tag", "t", []string{}, "--tag news --tag \"current affairs\"")
	updateCmd.Flags().StringP("new-id", "", "", "--new-id \"Public News Service\"")
	updateCmd.Flags().StringP("new-url", "", "", "--new-url https://www.abc.net.au/news")
	updateCmd.Flags().StringSliceP("new-tag", "", []string{}, "--new-tag free")
	updateCmd.Flags().StringSliceP("remove-tag", "", []string{}, "--remove-tag \"current affairs\"")
	updateCmd.Flags().Bool("remove-url", false, "--remove-url")
}

func combineUpdateArgs(flagSet *pflag.FlagSet, argv []string) (*runner.UpdateArgs, error) {

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

	newId, err := flagSet.GetString("new-id")
	if err != nil {
		return nil, err
	}

	newUrl, err := flagSet.GetString("new-url")
	if err != nil {
		return nil, err
	}

	newTags, err := flagSet.GetStringSlice("new-tag")
	if err != nil {
		return nil, err
	}

	removeTags, err := flagSet.GetStringSlice("remove-tag")
	if err != nil {
		return nil, err
	}

	removeUrl, err := flagSet.GetBool("remove-url")
	if err != nil {
		return nil, err
	}

	return runner.NewUpdateArgs(
		id,
		url,
		newId,
		newUrl,
		tags,
		newTags,
		removeTags,
		removeUrl,
	), nil
}
