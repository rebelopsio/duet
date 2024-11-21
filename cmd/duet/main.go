package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/rebelopsio/duet/internal/core/state"
)

var (
	cfgFile string
	store   *state.Store
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "duet",
	Short: "Duet - Infrastructure and Configuration in Harmony",
	Long: `Duet is a tool that orchestrates both infrastructure provisioning
and configuration management using Lua as its configuration language.
Complete documentation is available at https://github.com/rebeleopsio/duet`,
}

var applyCmd = &cobra.Command{
	Use:   "apply [file]",
	Short: "Apply infrastructure and configuration changes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleApply(args[0])
	},
}

var planCmd = &cobra.Command{
	Use:   "plan [file]",
	Short: "Show planned changes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handlePlan(args[0])
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.duet.yaml)")

	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(planCmd)

	// Initialize state store
	var err error
	store, err = state.NewStore("duet.db")
	if err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".duet")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func handleApply(filename string) error {
	// Implementation will go here
	return nil
}

func handlePlan(filename string) error {
	// Implementation will go here
	return nil
}
