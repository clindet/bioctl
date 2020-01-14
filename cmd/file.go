package cmd

import (
	gfile "github.com/openbiox/ganker/file"
	"github.com/spf13/cobra"
)

// FileClisT is the type to run FileCmd
type FileClisT struct {
	// CountLines run the countLine()
	CountLines bool
	CountChars bool
	CountBytes bool
	CountWords bool
	Format     string
}

// FileClis is the parameters to run FileCmd
var FileClis = FileClisT{}

// FileCmd is the command line cobra object for basic file operations
var FileCmd = &cobra.Command{
	Use:   "f",
	Short: "Basic file operations.",
	Long:  `Basic file operations..`,
	Run: func(cmd *cobra.Command, args []string) {
		setQuietLog(log, rootClis.Quite)
		fileCmdOptions(cmd, args)
	},
}

func fileCmdOptions(cmd *cobra.Command, args []string) {
	_, _, err := gfile.CountLineNameSlice(args)
	if err != nil {
		for _, v := range err {
			log.Warnln(v)
		}
		rootClis.HelpFlags = false
	}
	if rootClis.HelpFlags {
		cmd.Help()
	}
}

func init() {
	FileCmd.Flags().BoolVarP(&FileClis.CountLines, "count-lines", "l", false, "Count the lines")
	FileCmd.Flags().BoolVarP(&FileClis.CountChars, "count-chars", "", false, "Count the chars")
	FileCmd.Flags().BoolVarP(&FileClis.CountBytes, "count-bytes", "", false, "Count the bytes")
	FileCmd.Flags().BoolVarP(&FileClis.CountWords, "count-words", "", false, "Count the words")
	FileCmd.Flags().StringVarP(&FileClis.Format, "format", "", "plain", "Set the output format [plain, json, table]")
}
