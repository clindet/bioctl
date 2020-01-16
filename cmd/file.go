package cmd

import (
	cvrt "github.com/openbiox/ligo/convert"
	gfile "github.com/openbiox/ligo/file"
	"github.com/spf13/cobra"
)

// FnClisT is the type to run FnCmd
type FnClisT struct {
	// CountLines run the countLine()
	CountLines bool
	CountChars bool
	CountBytes bool
	CountWords bool
	Format     string
}

// FnClis is the parameters to run FnCmd
var FnClis = FnClisT{}
var fileSubCmdName = "fn"

// FnCmd is the command line cobra object for basic file operations
var FnCmd = &cobra.Command{
	Use:   fileSubCmdName,
	Short: "Conduct basic file operations.",
	Long:  `Conduct basic file operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		FnCmdOptions(cmd, args)
	},
}

func FnCmdOptions(cmd *cobra.Command, args []string) {
	err := []error{}
	if FnClis.CountLines || FnClis.CountBytes {
		initCmd(cmd, args)
		if rootClis.Verbose == 2 {
			logEnv.Infof("fnClis: %v", cvrt.Struct2Map(FnClis))
		}
	}
	if FnClis.CountLines {
		_, _, err = gfile.LineCounterByNameSlice(args)
		rootClis.HelpFlags = false

	} else if FnClis.CountBytes {
		_, _, err = gfile.BytesCounterByNameSlice(args)
		rootClis.HelpFlags = false
	}
	if err != nil {
		for _, v := range err {
			log.Warnln(v)
		}
	}
	// || FnClis.CountWords  || FnClis.CountChars ||
	if rootClis.HelpFlags {
		cmd.Help()
	}
}

func init() {
	FnCmd.Flags().BoolVarP(&FnClis.CountLines, "count-lines", "l", false, "count the lines")
	FnCmd.Flags().BoolVarP(&FnClis.CountBytes, "count-bytes", "c", false, "count the bytes")
	FnCmd.Flags().StringVarP(&FnClis.Format, "format", "", "plain", "set the output format [plain, json, table]")
}
