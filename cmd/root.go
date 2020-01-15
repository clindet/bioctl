package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootClisT is the bioctl global flags
type RootClisT struct {
	// version of bioctl
	Version   string
	Quiet     bool
	SaveLog   bool
	TaskID    string
	LogDir    string
	Clean     bool
	HelpFlags bool
}

var rootClis = RootClisT{
	Version:   "v0.1.0",
	Quiet:     false,
	HelpFlags: true,
}

var rootCmd = &cobra.Command{
	Use:   "bioctl",
	Short: "A simple command line tool to facilitate the data analysis",
	Long:  `A simple command line tool to facilitate the data analysis. More see here https://github.com/openbiox/bioctl.`,
	Run: func(cmd *cobra.Command, args []string) {
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
	rootCmd.AddCommand(FileCmd)
	rootCmd.AddCommand(FmtCmd)
	rootCmd.AddCommand(ParCmd)
	setGlobalFlag(rootCmd)
}
