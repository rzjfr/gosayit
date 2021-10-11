package cmd

import (
	"github.com/spf13/cobra"
)

var Quiet bool
var Verbose bool
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "sayit",
	Short:   "Simple tool to check given vocabulary against Oxford Learner's Dictionary",
	Long:    "You can get pronunciation and other information about the given word",
	Version: "0.0.1",
	Example: "sayit check hello",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&Quiet, "audio-only", "o", false, "Only play the pronunciation")
}
