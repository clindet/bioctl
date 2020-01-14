package cmd

import (
	"fmt"
	"os"

	glog "github.com/openbiox/ganker/log"
	"github.com/spf13/cobra"
)

// RootClisT is the ganker global flags
type RootClisT struct {
	// version of ganker
	Version string
	// print debug inforamtion
	Quite string
	// help flag
	HelpFlags bool
}

var rootClis = RootClisT{
	Version:   "v0.0.1",
	Quite:     "true",
	HelpFlags: true,
}
var log = glog.Logger

var rootCmd = &cobra.Command{
	Use:   "ganker",
	Short: "Just for fun of bioinformatics.",
	Long:  `Cross-platform command line tools to process sequence, alignment, count, and associated files. More see here https://github.com/openbiox/ganker.`,
	Run: func(cmd *cobra.Command, args []string) {
		setQuietLog(log, rootClis.Quite)
		rootClis.HelpFlags = true
		if rootClis.HelpFlags {
			cmd.Help()
		}
	},
}

// Execute main interface of bget
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if !rootCmd.HasFlags() && !rootCmd.HasSubCommands() {
			rootCmd.Help()
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootClis.Quite, "quite", "q", "true", "Wheather to print debug information [true or false]")
	rootCmd.AddCommand(FileCmd)
}
