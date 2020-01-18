package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/openbiox/ligo/stringo"

	cvrt "github.com/openbiox/ligo/convert"
	"github.com/openbiox/ligo/flag"
	"github.com/spf13/cobra"
)

// RangeClisT is the type to run bioctl Range
type RangeClisT struct {
	Start       float64
	StartChar   string
	EndChar     string
	End         float64
	Step        float64
	RefSequence string
	Sep         string
	Mode        string
}

// RangeClis is the parameters to run par.Tasks
var RangeClis RangeClisT

// RangeCmd is the command line of bioctl Range
var RangeCmd = &cobra.Command{
	Use:   "range [start end step]",
	Short: "Functions to manipulate intervals.",
	Long:  `Functions to manipulate intervals.`,
	Run: func(cmd *cobra.Command, args []string) {
		RangeCmdRunOptions(cmd, args)
	},
}

func RangeCmdRunOptions(cmd *cobra.Command, args []string) {
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
	if len(args) >= 1 && RangeClis.Mode == "num" {
		i64, _ := strconv.ParseFloat(args[0], 64)
		RangeClis.Start = i64
	}
	if len(args) >= 2 && RangeClis.Mode == "num" {
		i64, _ := strconv.ParseFloat(args[1], 64)
		RangeClis.End = i64
	}
	if len(args) >= 3 && RangeClis.Mode == "num" {
		i64, _ := strconv.ParseFloat(args[2], 64)
		RangeClis.Step = i64
	}
	if len(cleanArgs) >= 1 || RangeClis.Step != 0 {
		initCmd(cmd, args)
		if rootClis.Verbose == 2 {
			logEnv.Infof("RangeClis: %v", cvrt.Struct2Map(RangeClis))
		}
		if RangeClis.Start > RangeClis.End {
			log.Warn("start need lower than end")
			return
		}
		if RangeClis.Mode == "num" {
			tmpStr := stringo.StrReplaceAll(fmt.Sprintf("%v", newNumSlice(RangeClis.Start, RangeClis.End, RangeClis.Step)), "[[]|[]]", "")
			tmpStr = stringo.StrReplaceAll(tmpStr, " ", RangeClis.Sep)
			fmt.Println(tmpStr)
		} else {
			fmt.Println(strings.Join(newAlphabetSlice(RangeClis.StartChar, RangeClis.EndChar, int(RangeClis.Step)), RangeClis.Sep))
		}
		rootClis.HelpFlags = false
	}
	if rootClis.HelpFlags {
		cmd.Help()
	}
}

func newNumSlice(start, end, step float64) []float64 {
	if step <= 0 || end < start {
		return []float64{}
	}
	s := make([]float64, 0, 1+(int64(end)-int64(start)/int64(step)))
	for start <= end {
		s = append(s, start)
		start += step
	}
	return s
}

func newAlphabetSlice(start, end string, step int) []string {
	var k int
	var startIdx int
	var endIdx int
	var final []string
	for k < len(RangeClis.RefSequence) {
		if string(RangeClis.RefSequence[k]) == start {
			startIdx = k
		}
		if string(RangeClis.RefSequence[k]) == end {
			endIdx = k
		}
		k++
	}
	for startIdx <= endIdx {
		final = append(final, string(RangeClis.RefSequence[startIdx]))
		startIdx += step
	}
	return final
}

func init() {
	RangeCmd.Flags().Float64VarP(&RangeClis.Start, "start", "", 0, "start of range")
	RangeCmd.Flags().Float64VarP(&RangeClis.End, "end", "", 0, "end of range")
	RangeCmd.Flags().StringVarP(&RangeClis.StartChar, "start-char", "", "a", "start of char range")
	RangeCmd.Flags().StringVarP(&RangeClis.EndChar, "end-char", "", "z", "end of char range")
	RangeCmd.Flags().StringVarP(&RangeClis.RefSequence, "ref-str", "", "abcdefghijklmnopqrstuvwxyz", "reference char sequence")
	RangeCmd.Flags().Float64VarP(&RangeClis.Step, "step", "", 0, "step of range")
	RangeCmd.Flags().StringVarP(&RangeClis.Mode, "mode", "m", "num", "num or alphabet")
	RangeCmd.Flags().StringVarP(&RangeClis.Sep, "sep", "", " ", "seperator of output")

	RangeCmd.Example = `  # range number sequence
  bioctl range 1 100 2.5
  bioctl range --start 2 --end 1000 --step 15
	
  # char mode
  bioctl range --mode char --step 5
  bioctl range --mode char --step 3 --start-char s --end-char z --step 2 --ref-str qrstuvwxy12442z
  bioctl range --mode char --step 3 --start-char s --end-char z --step 2 --ref-str qrstuvwxy12442z --sep ''`
}
