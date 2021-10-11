package cmd

import (
	"fmt"
	"os"

	"github.com/rzjfr/sayit/audio"
	"github.com/rzjfr/sayit/dict"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks the give word",
	Long:  "Pronunciation of the given word. Also checks for the definition.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("Please provide a word to be checked")
		}
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		audio.InitLogger(Verbose)
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := dict.Spell(args[0])
		if err != nil {
			fmt.Println("Cannot say that!")
			fmt.Println(err)
			os.Exit(1)
		}
		err = audio.Play(args[0])
		if err != nil {
			fmt.Println("Cannot say that!")
			if !Verbose {
				fmt.Println("Please rerun with -v or --verbose to see more details")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().BoolP("fav", "f", false, "Add to favourite list")
}
