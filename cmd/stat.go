package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/openbiox/ligo/io"
	"github.com/tealeg/xlsx"

	"github.com/go-gota/gota/dataframe"
	gfile "github.com/openbiox/ligo/file"
	"github.com/openbiox/ligo/stringo"

	"github.com/montanaflynn/stats"
	cvrt "github.com/openbiox/ligo/convert"
	"github.com/openbiox/ligo/flag"
	"github.com/spf13/cobra"
)

// StatClisT is the type to run bioctl Stat
type StatClisT struct {
	Max           bool
	Min           bool
	Median        bool
	Mean          bool
	Mfreq         bool
	Variance      bool
	Sum           bool
	Pearson       bool
	Percentile    float64
	Freq          bool
	CuSum         bool
	GeometricMean bool
	HarmonicMean  bool
	Entropy       bool
	Covariance    bool
	SplitArgs     string
}

// StatDfClisT is the type to run bioctl statdf
type StatDfClisT struct {
	Print      bool
	Header     string
	At         string
	SheetIndex int
}

type StatFnClisT struct {
	// CountLines run the countLine()
	CountLines bool
	CountChars bool
	CountBytes bool
	CountWords bool
}

// StatClis is the parameters to run stat
var StatClis StatClisT

// StatDfClis is the parameters to run statdf
var StatDfClis StatDfClisT

// StatFnClis is the parameters to run StatFnCmd
var StatFnClis = StatFnClisT{}

// StatCmd is the command line of bioctl Stat
var StatCmd = &cobra.Command{
	Use:   "stat",
	Short: "Functions related to stat.",
	Long:  `Functions related to stat.`,
	Run: func(cmd *cobra.Command, args []string) {
		StatCmdRunOptions(cmd, args)
		if rootClis.HelpFlags {
			cmd.Help()
		}
	},
}

// StatDfCmd is the command line of bioctl statdf
var StatDfCmd = &cobra.Command{
	Use:   "statdf",
	Short: "Functions related to stat data.frame.",
	Long:  `Functions related to stat data.frame.`,
	Run: func(cmd *cobra.Command, args []string) {
		StatDfCmdOptions(cmd, args)
		if rootClis.HelpFlags {
			cmd.Help()
		}
	},
}

// StatFnCmd is the command line of bioctl statfn
var StatFnCmd = &cobra.Command{
	Use:   "statfn",
	Short: "Functions related to stat file.",
	Long:  `Functions related to stat file.`,
	Run: func(cmd *cobra.Command, args []string) {
		StatFnCmdOptions(cmd, args)
		if rootClis.HelpFlags {
			cmd.Help()
		}
	},
}

func StatCmdRunOptions(cmd *cobra.Command, args []string) {
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
	argsNew := []string{}
	for _, v := range args {
		if stringo.StrDetect(v, "^ *") {
			v = stringo.StrRemoveAll(v, "^ *")
		}
		if strings.Contains(v, " ") && StatClis.SplitArgs == "true" {
			argsNew = append(argsNew, stringo.StrSplit(v, " ", 10000000)...)
		} else if strings.Contains(v, ",") && StatClis.SplitArgs == "true" {
			argsNew = append(argsNew, stringo.StrSplit(v, ",", 10000000)...)
		} else if stringo.StrDetect(v, "\t") && StatClis.SplitArgs == "true" {
			argsNew = append(argsNew, stringo.StrSplit(v, "\t", 10000000)...)
		} else {
			argsNew = append(argsNew, v)
		}
	}
	args = argsNew
	if StatClis.Freq {
		statFreq(args)
		rootClis.HelpFlags = false
		return
	}
	if len(cleanArgs) >= 1 || enableStat() {
		initCmd(cmd, args)
		if rootClis.Verbose == 2 {
			logEnv.Infof("StatClis: %v", cvrt.Struct2Map(StatClis))
		}
		var data, data2 stats.Float64Data
		if !twoDimData() {
			data = stats.LoadRawData(args)
		} else {
			data = stats.LoadRawData(args[0:(len(args) / 2)])
			data2 = stats.LoadRawData(args[len(args)/2 : len(args)])
			if data.Len() < 100 {
				fmt.Printf("Vector1: %v\n", data)
			} else {
				fmt.Printf("Vector1: [")
				for i := 0; i < 100; i++ {
					fmt.Printf("%v ", data.Get(i))
				}
				fmt.Println("...]")
			}
			if data2.Len() < 100 {
				fmt.Printf("Vector2: %v\n", data2)
			} else {
				fmt.Printf("Vector2: [")
				for i := 0; i < 100; i++ {
					fmt.Printf("%v ", data2.Get(i))
				}
				fmt.Println("...]")
			}
		}
		runStat(&StatClis, &data, &data2)
		rootClis.HelpFlags = false
	}
}

func StatDfCmdOptions(cmd *cobra.Command, args []string) {
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
	if len(cleanArgs) >= 1 || StatDfClis.Print {
		initCmd(cmd, args)
		if rootClis.Verbose == 2 {
			logEnv.Infof("StatDfClis: %v", cvrt.Struct2Map(StatDfClis))
		}
		var df []*dataframe.DataFrame
		if StatDfClis.Print {
			df = readFile(args)
			for k, v := range df {
				fmt.Printf("input[%d]:", k)
				fmt.Println(v)
			}
		}
		data := []string{}
		if StatDfClis.At != "" {
			if len(df) == 0 {
				df = readFile(args)
			}
			index := strings.Split(StatDfClis.At, ",")
			row, colum := []int{-1}, []int{-1}
			if index[0] != "" {
				row = convertIndex(index[0])
			}
			if len(index) > 1 && index[1] != "" {
				colum = convertIndex(index[1])
			}
			for k, v := range df {
				if row[0] != -1 && colum[0] != -1 {
					fmt.Printf("selected[%d]:", k)
					fmt.Println(v.Elem(row[0], colum[0]))
				} else if row[0] != -1 {
					row2 := v.Subset(row)
					if enableStat() {
						data = row2.Records()[1]
					}
					if StatDfClis.Print {
						fmt.Printf("selected[%d]:", k)
						fmt.Println(row2)
					}
				} else if colum[0] != -1 {
					colum2 := v.Select(colum)
					if StatDfClis.Print {
						fmt.Printf("selected[%d]:", k)
						fmt.Println(colum2)
					}
					if enableStat() {
						for i := 0; i < len(colum); i++ {
							for _, v2 := range colum2.Records()[1:] {
								data = append(data, v2[i])
							}
						}
					}
				}
			}
		}
		if len(data) > 0 && enableStat() {
			StatCmdRunOptions(cmd, data)
		}
		rootClis.HelpFlags = false
	}
}

func enableStat() bool {
	return StatClis.Max || StatClis.Min || StatClis.Mean || StatClis.Median ||
		StatClis.Freq || StatClis.Mfreq || StatClis.Pearson || StatClis.Sum || StatClis.Variance || StatClis.Percentile != -1 || StatClis.CuSum || StatClis.GeometricMean || StatClis.HarmonicMean || StatClis.Entropy || StatClis.Covariance
}

func twoDimData() bool {
	return StatClis.Covariance || StatClis.Pearson
}

func StatFnCmdOptions(cmd *cobra.Command, args []string) {
	err := []error{}
	if StatFnClis.CountLines || StatFnClis.CountBytes {
		initCmd(cmd, args)
		if rootClis.Verbose == 2 {
			logEnv.Infof("StatFnClis: %v", cvrt.Struct2Map(StatFnClis))
		}
	}
	if StatFnClis.CountLines {
		_, _, err = gfile.LineCounterByNameSlice(args)
		rootClis.HelpFlags = false

	} else if StatFnClis.CountBytes {
		_, _, err = gfile.BytesCounterByNameSlice(args)
		rootClis.HelpFlags = false
	}
	if err != nil {
		for _, v := range err {
			log.Warnln(v)
		}
	}
}

func statFreq(args []string) {
	count := make(map[string]int)
	for _, v := range args {
		count[v]++
	}
	for k, v := range count {
		fmt.Printf("%d\t%v\n", v, k)
	}
	return
}

func readFile(args []string) (df []*dataframe.DataFrame) {
	for _, v := range args {
		dfTmp := dataframe.DataFrame{}

		if stringo.StrDetect(v, ".csv$") {
			of, _ := io.Open(v)
			dfTmp = dataframe.ReadCSV(of, dataframe.HasHeader(StatDfClis.Header == "true"))
			defer of.Close()
		} else if stringo.StrDetect(v, ".txt$") {
			rows := io.ReadLines(v)
			dfTxt := [][]string{}
			for _, row := range rows {
				dfTxt = append(dfTxt, stringo.StrSplit(row, "\t", 10000000))
			}
			dfTmp = dataframe.LoadRecords(
				dfTxt,
				dataframe.HasHeader(StatDfClis.Header == "true"),
			)
		} else if stringo.StrDetect(v, ".json$") {
			of, _ := io.Open(v)
			dfTmp = dataframe.ReadJSON(of, dataframe.HasHeader(StatDfClis.Header == "true"))
			defer of.Close()
		} else if stringo.StrDetect(v, ".xlsx$") {
			xlFile, _ := xlsx.OpenFile(v)
			sheet := xlFile.Sheets[StatDfClis.SheetIndex]
			dfXlsx := [][]string{}
			for _, row := range sheet.Rows {
				rowTmp := []string{}
				for _, cell := range row.Cells {
					rowTmp = append(rowTmp, cell.String())
				}
				dfXlsx = append(dfXlsx, rowTmp)
			}
			dfTmp = dataframe.LoadRecords(
				dfXlsx,
				dataframe.HasHeader(StatDfClis.Header == "true"),
			)
		}
		df = append(df, &dfTmp)
	}
	return df
}

func runStat(StatClis *StatClisT, data *stats.Float64Data, data2 *stats.Float64Data) {
	var val float64
	var valSlice []float64
	if StatClis.Max {
		val, _ = data.Max()
	} else if StatClis.Min {
		val, _ = data.Min()
	} else if StatClis.Mean {
		val, _ = data.Mean()
	} else if StatClis.Median {
		val, _ = data.Median()
	} else if StatClis.Mfreq {
		valSlice, _ = data.Mode()
	} else if StatClis.Variance {
		val, _ = data.Variance()
	} else if StatClis.Sum {
		val, _ = data.Sum()
	} else if StatClis.Pearson {
		val, _ = data.Pearson(*data2)
	} else if StatClis.Percentile != -1 {
		val, _ = data.Percentile(StatClis.Percentile)
	} else if StatClis.CuSum {
		valSlice, _ = data.CumulativeSum()
	} else if StatClis.GeometricMean {
		val, _ = data.GeometricMean()
	} else if StatClis.HarmonicMean {
		val, _ = data.HarmonicMean()
	} else if StatClis.Entropy {
		val, _ = data.Entropy()
	} else if StatClis.Covariance {
		val, _ = data.Covariance(*data2)
	}

	if !(StatClis.Mfreq || StatClis.CuSum) {
		fmt.Println(val)
	} else {
		fmt.Println(stringo.StrReplaceAll(fmt.Sprintf("%v", valSlice), "[[]|[]]", ""))
	}
}

func init() {
	setStatFlag(StatCmd)
	setStatFlag(StatDfCmd)
	StatDfCmd.Flags().IntVarP(&(StatDfClis.SheetIndex), "sheet", "", 0, "xlsx sheet index for input data.frame.")
	StatDfCmd.Flags().BoolVarP(&(StatDfClis.Print), "print", "", false, "print involved data.frame.")
	StatDfCmd.Flags().StringVarP(&(StatDfClis.Header), "header", "", "true", "set first row as header.")
	StatDfCmd.Flags().StringVarP(&(StatDfClis.At), "at", "", "", "extract values: '1: first row', ',1: first colum', '1,2: [1,2] element'")

	StatFnCmd.Flags().BoolVarP(&StatFnClis.CountLines, "count-lines", "l", false, "count the lines")
	StatFnCmd.Flags().BoolVarP(&StatFnClis.CountBytes, "count-bytes", "c", false, "count the bytes")

	StatCmd.Example = `  bioctl stat --min 1 2 3 4 5 100
  bioctl stat --max ' -1 1 2 3 4 5'
  bioctl stat --mean ' -10 1 3 5 7 9 26 100'
  bioctl stat --median ' -10 1 10 3 4 100 143 123 12 22.2'
  bioctl stat --mfreq '2 10 2 2 10 10 11 12 14 -10 -10'
  bioctl stat --var '2 10 2 2 10 20'
  bioctl stat --sum ' -30 1 1 2 3 2 14234 12 12 1331 23 12 12'
  // vector1: -1-6 vector2: -10-51
  bioctl stat --pearson ' -1 2 3 4 5 6 -10 20 30 40 50 51'
  bioctl stat --percentile 30 '1 2 3 4 5 6 7 8 9 10'
  bioctl stat --freq 'a -1 -1 -2 -2 1 2 3 2 1 1 1 1 1 1' | sort
  bioctl stat --csum '1,2,3,2,1,1,1,1,1,1'
  bioctl stat --mean-geometric '1 2 3 2 1 1 1 1 1'
  bioctl stat --mean-harmonic '1 2 3 2 1 1 1 1 1'
  bioctl stat --entropy '1 2 3 2 1 1 1 1 1'`

	StatDfCmd.Example = `  bioctl statdf --min _examples/test.csv --at ",0"
  bioctl statdf --max _examples/test.csv --at ",0"
  bioctl statdf --mean _examples/test.csv --at ",0"
  bioctl statdf --median _examples/test.csv --at ",0"
  bioctl statdf --mfreq _examples/test.csv --at ",0"
  bioctl statdf --var _examples/test.csv --at ",0"
  bioctl statdf --sum _examples/test.csv --at ",0"
  bioctl statdf --pearson _examples/test.csv --at ",0:1" --print
  bioctl statdf --percentile 30 _examples/test.csv --at ",0"
  bioctl statdf --freq _examples/test.csv --at ",0" | sort
  bioctl statdf --csum _examples/test.csv --at ",0"
  bioctl statdf --mean-geometric _examples/test.csv --at ",0"
  bioctl statdf --mean-harmonic _examples/test.csv --at ",0"
  bioctl statdf --entropy _examples/test.csv --at ",0"`
}

func setStatFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&(StatClis.Max), "max", "", false, "returns max value.")
	cmd.Flags().BoolVarP(&(StatClis.Min), "min", "", false, "returns min value.")
	cmd.Flags().BoolVarP(&(StatClis.Median), "median", "", false, "returns median value.")
	cmd.Flags().BoolVarP(&(StatClis.Mean), "mean", "", false, "returns mean value.")
	cmd.Flags().BoolVarP(&(StatClis.Mfreq), "mfreq", "", false, "returns the most frequent value.")
	cmd.Flags().BoolVarP(&(StatClis.Variance), "var", "", false, "returns the amount of variation in the dataset.")
	cmd.Flags().BoolVarP(&(StatClis.Sum), "sum", "", false, "returns the sum of value.")
	cmd.Flags().BoolVarP(&(StatClis.Pearson), "pearson", "", false, "returns the pearson product-moment correlation coefficient between two group variables.")
	cmd.Flags().Float64VarP(&(StatClis.Percentile), "percentile", "", -1, "returns the relative standing in a slice of floats.")
	cmd.Flags().BoolVarP(&(StatClis.Freq), "freq", "", false, "returns the frequency stat.")
	cmd.Flags().BoolVarP(&(StatClis.CuSum), "csum", "", false, "calculates the cumulative sum of the input slice.")
	cmd.Flags().BoolVarP(&(StatClis.GeometricMean), "mean-geometric", "", false, "calculates the geometric mean.")
	cmd.Flags().BoolVarP(&(StatClis.HarmonicMean), "mean-harmonic", "", false, "calculates the harmonic mean.")
	cmd.Flags().BoolVarP(&(StatClis.Entropy), "entropy", "", false, "calculates the entropy.")
	cmd.Flags().BoolVarP(&(StatClis.Covariance), "covar", "", false, "a measure of how much two sets of data change.")
	cmd.Flags().StringVarP(&(StatClis.SplitArgs), "split-args", "", "true", "split input args by comma, space or tab .")

}
