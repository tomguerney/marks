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
	"github.com/abitofoldtom/marks/opener"
	"github.com/abitofoldtom/marks/printer"
	"github.com/abitofoldtom/marks/prompter"
	"github.com/abitofoldtom/marks/runner"
	"github.com/abitofoldtom/marks/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open id [tags...]",
	Short: "Open a url in a browser",
	RunE:  runOpen,
}

func runOpen(cmd *cobra.Command, argv []string) error {
	args, err := combineOpenArgs(cmd.Flags(), argv)
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
	opener := opener.NewOpener(config)
	runner := runner.NewOpenRunner(args, config, markService, printer, prompter, opener)
	if err := runner.Run(); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(openCmd)
	openCmd.Flags().StringP("url", "u", "", "(can be partial) --url abc.net.au")
	openCmd.Flags().StringSliceP("tag", "t", []string{}, "--tag news --tag \"current affairs\"")
	openCmd.PersistentFlags().StringP("browser", "b", "", "--browser firefox")
	viper.BindPFlag("browser", openCmd.PersistentFlags().Lookup("browser"))
}

func combineOpenArgs(flagSet *pflag.FlagSet, argv []string) (*runner.OpenArgs, error) {

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

	return runner.NewOpenArgs(id, url, tags), nil
}
