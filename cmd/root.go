package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var globalFlagConfigFile string
var globalFlagCore string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "corecli",
	Short:        "A CLI for the datarhei Core",
	Long:         `A CLI for the datarhei Core.`,
	SilenceUsage: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func SetArgs(a []string) {
	rootCmd.SetArgs(a)
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&globalFlagConfigFile, "config", "", "config file (default is $HOME/.corecli.json)")
	rootCmd.PersistentFlags().StringVar(&globalFlagCore, "core", "", "core to connect to instead of the selected one")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	viper.SetDefault("editor", "")
	viper.SetDefault("cores.selected", "")
	viper.SetDefault("cores.list", map[string]string{})

	if globalFlagConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(globalFlagConfigFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".corecli" (without extension).
		//viper.SetConfigFile(filepath.Join(home, ".corecli.json"))
		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".corecli")
		viper.SafeWriteConfig()
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			cobra.CheckErr(err)
		}
	}
}
