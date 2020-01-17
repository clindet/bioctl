package cmd

import (
	"bufio"
	"io/ioutil"
	"os"

	cvrt "github.com/openbiox/ligo/convert"
	"github.com/openbiox/ligo/flag"
	"github.com/spf13/cobra"
)

// ArchiveClisT is the type to run bioctl Archive
type ArchiveClisT struct {
}

// ArchiveClis is the parameters to run par.Tasks
var ArchiveClis ArchiveClisT

// ArchiveCmd is the command line of bioctl Archive
var ArchiveCmd = &cobra.Command{
	Use:   "arc",
	Short: "Archive related functions (compress or uncompress).",
	Long:  `Archive related functions (compress or uncompress).`,
	Run: func(cmd *cobra.Command, args []string) {
		ArchiveCmdRunOptions(cmd, args)
	},
}

func ArchiveCmdRunOptions(cmd *cobra.Command, args []string) {
	cleanArgs := []string{}
	hasStdin := false
	if cleanArgs, hasStdin = flag.CheckStdInFlag(cmd); hasStdin {
		reader := bufio.NewReader(os.Stdin)
		result, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal(err)
		} else if len(result) > 0 {
		}
	}
	if len(cleanArgs) >= 1 {
		initCmd(cmd, args)
		if rootClis.Verbose == 2 {
			logEnv.Infof("ArchiveClis: %v", cvrt.Struct2Map(ArchiveClis))
		}

		rootClis.HelpFlags = false
	}
	if rootClis.HelpFlags {
		cmd.Help()
	}
}

func init() {
	ArchiveCmd.Example = ``
}
