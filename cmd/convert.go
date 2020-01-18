package cmd

import (
	"bufio"
	"io/ioutil"
	"os"

	cvrt "github.com/openbiox/ligo/convert"
	"github.com/openbiox/ligo/flag"
	"github.com/openbiox/ligo/parse"
	"github.com/openbiox/ligo/stringo"
	"github.com/spf13/cobra"
)

// ConvertClisT is the type to run bioctl Convert
type ConvertClisT struct {
	NcbiXMLToJSON string
	NcbiXMLPaths  []string
	NcbiKeywords  string
}

// ConvertClis is the parameters to run par.Tasks
var ConvertClis ConvertClisT

// ConvertCmd is the command line of bioctl Convert
var ConvertCmd = &cobra.Command{
	Use:   "cvrt",
	Short: "Convert related functions.",
	Long:  `Convert related functions.`,
	Run: func(cmd *cobra.Command, args []string) {
		ConvertCmdRunOptions(cmd, args)
	},
}

func ConvertCmdRunOptions(cmd *cobra.Command, args []string) {
	cleanArgs := []string{}
	var stdin []byte
	var err error
	hasStdin := false
	if cleanArgs, hasStdin = flag.CheckStdInFlag(cmd); hasStdin {
		reader := bufio.NewReader(os.Stdin)
		stdin, err = ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(cleanArgs) >= 1 || ConvertClis.NcbiXMLToJSON != "" {
		initCmd(cmd, args)
		if rootClis.Verbose == 2 {
			logEnv.Infof("ConvertClis: %v", cvrt.Struct2Map(ConvertClis))
		}
		if ConvertClis.NcbiXMLToJSON == "pubmed" {
			if len(cleanArgs) >= 1 || len(stdin) > 0 {
				ConvertClis.NcbiXMLPaths = append(ConvertClis.NcbiXMLPaths, cleanArgs...)
				keywordsList := stringo.StrSplit(ConvertClis.NcbiKeywords, ", |,", 10000)
				parse.PubmedXML(&ConvertClis.NcbiXMLPaths, &stdin, rootClis.Outfn, &keywordsList, rootClis.Thread)
			}
			rootClis.HelpFlags = false
		} else if ConvertClis.NcbiXMLToJSON == "sra" {
			if len(cleanArgs) >= 1 || len(stdin) > 0 {
				ConvertClis.NcbiXMLPaths = append(ConvertClis.NcbiXMLPaths, cleanArgs...)
				keywordsList := stringo.StrSplit(ConvertClis.NcbiKeywords, ", |,", 10000)
				parse.SraXML(&ConvertClis.NcbiXMLPaths, &stdin, rootClis.Outfn, &keywordsList, rootClis.Thread)
			}
			rootClis.HelpFlags = false
		}
		rootClis.HelpFlags = false
	}
	if rootClis.HelpFlags {
		cmd.Help()
	}
}

func init() {
	ConvertCmd.Flags().StringVarP(&ConvertClis.NcbiXMLToJSON, "xml2json", "", "", "Convert XML files to json [e.g. pubmed, sra].")
	ConvertCmd.Flags().IntVarP(&rootClis.Thread, "thread", "t", 1, "thread to process.")

	ConvertCmd.Example = `  # convert Pubmed XML to clean JSON string
	bget api ncbi -q "Galectins control MTOR and AMPK in response to lysosomal damage to induce autophagy OR MTOR-independent autophagy induced by interrupted endoplasmic reticulum-mitochondrial Ca2+ communication: a dead end in cancer cells. OR The PARK10 gene USP24 is a negative regulator of autophagy and ULK1 protein stability OR Coordinate regulation of autophagy and the ubiquitin proteasome system by MTOR." | bioctl cvrt --xml2json pubmed
	
	# convert SRA XML to clean JSON string
  bget api ncbi -d 'sra' -q PRJNA527715 | bioctl cvrt --xml2json sra -`
}
