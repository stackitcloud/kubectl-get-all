/*
Copyright 2019 Cornelius Weig

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
	"io"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/util/homedir"

	"github.com/corneliusweig/ketall/cmd/internal"
	"github.com/corneliusweig/ketall/pkg/ketall"
	"github.com/corneliusweig/ketall/pkg/ketall/constants"
	"github.com/corneliusweig/ketall/pkg/ketall/options"
)

var (
	ketallOptions = options.NewCmdOptions()
	v             string
)

const (
	ketallLongDescription = `
Like 'kubectl get all', but get _really_ all resources

Ketall retrieves all resources which allow to be fetched. This complements the
usual "kubectl get all" command, which excludes all cluster-level and some
namespaced resources.

More on https://github.com/corneliusweig/ketall/blob/v1.2.1/doc/USAGE.md#usage
`
	ketallExamples = `
  Get all resources, excluding events
   $ ketall

  Get all resources, including events
   $ ketall --exclude=

  Get all resources created in the last minute
   $ ketall --since 1m

  Get all resources in the default namespace
   $ ketall --namespace=default

  Get all cluster level resources
   $ ketall --only-scope=cluster

  Some options can also be configured in the config file './ketall.yaml' or '~/.kube/ketall.yaml'
`
)

var rootCmd = &cobra.Command{
	Use:     internal.CommandName,
	Short:   "Like `kubectl get all`, but get _really_ all resources",
	Long:    internal.HelpTextMapName(ketallLongDescription),
	Args:    cobra.NoArgs,
	Example: internal.HelpTextMapName(ketallExamples),
	Run: func(cmd *cobra.Command, args []string) {
		ketall.KetAll(ketallOptions)
	},
}

func Execute() error {
	rootCmd.SetOutput(ketallOptions.Streams.Out)
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&ketallOptions.CfgFile, "config", "", "config file (default is $HOME/.kube/ketall.yaml)")
	rootCmd.PersistentFlags().StringVarP(&v, "verbosity", "v", constants.DefaultLogLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")

	rootCmd.Flags().BoolVar(&ketallOptions.UseCache, constants.FlagUseCache, false, "use cached list of server resources")
	rootCmd.Flags().StringVar(&ketallOptions.Scope, constants.FlagScope, "", "only resources with scope cluster|namespace")
	rootCmd.Flags().StringVar(&ketallOptions.Since, constants.FlagSince, "", "only resources younger than given age")
	rootCmd.Flags().StringVarP(&ketallOptions.Selector, constants.FlagSelector, "l", "", "selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	rootCmd.Flags().StringVar(&ketallOptions.FieldSelector, constants.FlagFieldSelector, "", "Selector (field query) to filter on, supports '=', '==', and '!='.(e.g. --field-selector key1=value1,key2=value2). The common field queries for all types are metadata.name and metadata.namespace.")
	rootCmd.Flags().StringSliceVar(&ketallOptions.Exclusions, constants.FlagExclude, []string{"events"}, "filter by resource name (plural form or short name)")

	ketallOptions.GenericCliFlags.AddFlags(rootCmd.Flags())
	ketallOptions.PrintFlags.AddFlags(rootCmd)

	err := viper.BindPFlags(rootCmd.Flags())
	if err != nil {
		logrus.Errorf("Cannot bind flags: %s", err)
	}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := SetUpLogs(ketallOptions.Streams.ErrOut, v); err != nil {
			return err
		}
		return nil
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if ketallOptions.CfgFile != "" {
		viper.SetConfigFile(ketallOptions.CfgFile)
	} else {
		// Search for "ketall.yaml" in "." and "~/.kube/"
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Join(homedir.HomeDir(), ".kube"))
		viper.SetConfigName("ketall")
	}

	// read in environment variables that match
	viper.SetEnvPrefix("ketall")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logrus.Debug("Using config file:", viper.ConfigFileUsed())
	}
}

func SetUpLogs(out io.Writer, level string) error {
	logrus.SetOutput(out)
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return errors.Wrap(err, "parsing log level")
	}
	logrus.SetLevel(lvl)
	logrus.Debugf("Set log-level to %s", level)
	return nil
}
