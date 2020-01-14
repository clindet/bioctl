package cmd

import (
	"fmt"
	"os"

	clog "github.com/openbiox/bioctl/log"
	"github.com/spf13/cobra"
)

// RootClisT is the bioctl global flags
type RootClisT struct {
	// version of bioctl
	Version string
	// print debug inforamtion
	Quiet string
	// help flag
	HelpFlags bool
}

var rootClis = RootClisT{
	Version:   "v0.0.1",
	Quiet:     "true",
	HelpFlags: true,
}
var wd string
var log = clog.Logger

var rootCmd = &cobra.Command{
	Use:   "bioctl",
	Short: "Just for fun of bioinformatics.",
	Long:  `Cross-platform command line tools to process sequence, alignment, count, and associated files. More see here https://github.com/openbiox/bioctl.`,
	Run: func(cmd *cobra.Command, args []string) {
		clog.SetQuietLog(log, rootClis.Quiet)
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
	rootCmd.Version = "0.1.0"
	rootCmd.PersistentFlags().StringVarP(&rootClis.Quiet, "quite", "q", "true", "Wheather to print debug information [true or false]")
	rootCmd.AddCommand(FileCmd)
	rootCmd.AddCommand(FmtCmd)
	rootCmd.AddCommand(ParCmd)
}
