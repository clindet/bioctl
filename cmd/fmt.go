package cmd

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"

	cvrt "github.com/openbiox/ligo/convert"
	"github.com/openbiox/ligo/flag"
	gfmt "github.com/openbiox/ligo/fmt"
	"github.com/spf13/cobra"
)

// FmtClis is the parameters to run FmtCmd
var FmtClis = gfmt.ClisT{}

// FmtCmd is the command line cobra object for basic re-format
var FmtCmd = &cobra.Command{
	Use:   "fmt [input1 input2]",
	Short: "A set of format (fmt) command.",
	Long:  `A set of format (fmt) command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmtCmdRunOptions(cmd, args)
	},
}

func fmtCmdRunOptions(cmd *cobra.Command, args []string) {
	cleanArgs := []string{}
	hasStdin := false
	if cleanArgs, hasStdin = flag.CheckStdInFlag(cmd); hasStdin {
		reader := bufio.NewReader(os.Stdin)
		result, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal(err)
		} else if len(result) > 0 {
			FmtClis.Stdin = &result
		}
	} else {
		FmtClis.Stdin = nil
	}

	if len(cleanArgs) >= 1 || hasStdin {
		initCmd(cmd, args)
		logEnv.Infof("env (fmt): %v", cvrt.Struct2Map(FmtClis))
		logBash.Infof("%s %s", cmd.CommandPath(), strings.Join(args, " "))
		FmtClis.Files = &cleanArgs
		runFlag := false
		if FmtClis.PrettyJSON {
			gfmt.PrettyJSON(&FmtClis, FmtClis.Thread)
			runFlag = true
		} else if FmtClis.JSONToSlice {
			gfmt.JSON2Slice(&FmtClis, FmtClis.Thread)
			runFlag = true
		}
		if !runFlag && hasStdin {
			io.Copy(os.Stdout, bytes.NewBuffer(*FmtClis.Stdin))
		}
		rootClis.HelpFlags = false
	}
	if rootClis.HelpFlags {
		cmd.Help()
	}
}

func init() {
	FmtCmd.Flags().IntVarP(&FmtClis.Thread, "thread", "t", 1, "thread to process.")
	FmtCmd.Flags().BoolVarP(&FmtClis.JSONToSlice, "json-to-slice", "", false, "convert key-value JSON  to []key-value and easy to export to readable table.")
	FmtCmd.Flags().BoolVarP(&FmtClis.PrettyJSON, "json-pretty", "", false, "pretty json files.")
	FmtCmd.Flags().IntVarP(&FmtClis.Indent, "indent", "", 4, "control the indent of output json files.")
	FmtCmd.Flags().BoolVarP(&FmtClis.SortKeys, "sort-keys", "", false, "control wheather to sort JSON key.")
	FmtCmd.Example = `  bget api ncbi -q "Galectins control MTOR and AMPK in response to lysosomal damage to induce autophagy OR MTOR-independent autophagy induced by interrupted endoplasmic reticulum-mitochondrial Ca2+ communication: a dead end in cancer cells. OR The PARK10 gene USP24 is a negative regulator of autophagy and ULK1 protein stability OR Coordinate regulation of autophagy and the ubiquitin proteasome system by MTOR." | bget api ncbi --xml2json pubmed - | sed 's;}{;,;g' | bioctl fmt --json-to-slice --indent 4 -| json2csv -o final.csv`
	JSON := make(map[int]map[string]interface{})
	FmtClis.JSON = &JSON
	Table := make(map[int][]interface{})
	FmtClis.Table = &Table
}
