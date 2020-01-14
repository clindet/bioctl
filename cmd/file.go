package cmd

import (
	gfile "github.com/openbiox/bioctl/file"
	clog "github.com/openbiox/bioctl/log"
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
	Use:   "fn",
	Short: "Conduct basic file operations.",
	Long:  `Conduct basic file operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		clog.SetQuietLog(log, rootClis.Quiet)
		fileCmdOptions(cmd, args)
	},
}

func fileCmdOptions(cmd *cobra.Command, args []string) {
	err := []error{}
	if FileClis.CountLines {
		_, _, err = gfile.LineCounterByNameSlice(args)
		rootClis.HelpFlags = false

	} else if FileClis.CountBytes {
		_, _, err = gfile.BytesCounterByNameSlice(args)
		rootClis.HelpFlags = false
	}
	if err != nil {
		for _, v := range err {
			log.Warnln(v)
		}
	}
	// || FileClis.CountWords  || FileClis.CountChars ||
	if rootClis.HelpFlags {
		cmd.Help()
	}
}

func init() {
	FileCmd.Flags().BoolVarP(&FileClis.CountLines, "count-lines", "l", false, "count the lines")
	FileCmd.Flags().BoolVarP(&FileClis.CountBytes, "count-bytes", "c", false, "count the bytes")
	FileCmd.Flags().StringVarP(&FileClis.Format, "format", "", "plain", "set the output format [plain, json, table]")
}
