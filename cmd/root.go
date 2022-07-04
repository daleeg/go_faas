package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cliName        = "go_faas"
	cliDescription = "A simple command line client for go_faas"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   cliName,
	Short: cliDescription,
	Long:  cliDescription,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("go faas success ")
	},
}

// Start adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Start() {
	rootCmd.SetUsageFunc(usageFunc)
	// Make help just show the usage
	rootCmd.SetHelpTemplate(`{{.UsageString}}`)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./test.yaml", "config file (default is ./test.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("test")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())

		// fmt.Println("mysql password: ", viper.Get("mysql.password"))
		// fmt.Println("mysql port: ", viper.Get("mysql.port"))
	}
}

func usageFunc(c *cobra.Command) error {
	return nil
}