package cmd

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	uuid "github.com/satori/go.uuid"

	cvrt "github.com/openbiox/ligo/convert"
	"github.com/openbiox/ligo/flag"
	"github.com/spf13/cobra"
)

// RandClisT is the type to run bioctl Rand
type RandClisT struct {
	UUID   bool
	Date   bool
	String bool
	Int    bool
	Float  bool
	Length int
	Number int
}

// RandClis is the parameters to run par.Tasks
var RandClis RandClisT

// RandCmd is the command line of bioctl Rand
var RandCmd = &cobra.Command{
	Use:   "rand",
	Short: "Functions related to random events.",
	Long:  `Functions related to random events.`,
	Run: func(cmd *cobra.Command, args []string) {
		RandCmdRunOptions(cmd, args)
	},
}

func RandCmdRunOptions(cmd *cobra.Command, args []string) {
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
	if len(cleanArgs) >= 1 || validFlag() {
		initCmd(cmd, args)
		if rootClis.Verbose == 2 {
			logEnv.Infof("RandClis: %v", cvrt.Struct2Map(RandClis))
		}
		if RandClis.UUID {
			randOut(func() { fmt.Println(uuid.NewV4()) })
		} else if RandClis.Int {
			randOut(func() { randInt(RandClis.Length) })
		} else if RandClis.String {
			randOut(func() { randStr(RandClis.Length) })
		}
		rootClis.HelpFlags = false
	}
	if rootClis.HelpFlags {
		cmd.Help()
	}
}

func validFlag() bool {
	return RandClis.UUID || RandClis.Date || RandClis.Float ||
		RandClis.Int || RandClis.String
}

func randOut(f func()) {
	count := 1
	for {
		f()
		count++
		if count > RandClis.Number {
			break
		}
	}
}

func randInt(len int) {
	var pool = []byte{1, 2, 3, 4, 5, 7, 8, 9}
	var container string
	length := bytes.NewReader(pool).Len()

	for i := 1; i <= len; i++ {
		random, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
		if err != nil {

		}
		container += fmt.Sprintf("%d", pool[random.Int64()])
	}
	fmt.Println(container)
}

func randStr(len int) {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	fmt.Println(container)
}

func init() {
	RandCmd.Flags().BoolVarP(&RandClis.UUID, "uuid", "", false, "generate random uuid")
	RandCmd.Flags().BoolVarP(&RandClis.Date, "date", "", false, "generate random date")
	RandCmd.Flags().BoolVarP(&RandClis.String, "str", "", false, "generate random string")
	RandCmd.Flags().BoolVarP(&RandClis.Int, "int", "", false, "generate random string")
	RandCmd.Flags().BoolVarP(&RandClis.Float, "float", "", false, "generate random string")

	RandCmd.Flags().IntVarP(&RandClis.Length, "len", "l", 15, "length of random string or number")
	RandCmd.Flags().IntVarP(&RandClis.Number, "num", "n", 1, "number of random output")

	RandCmd.Example = `  bioctl rand --uuid -n 10
  bioctl rand --str -l 35 -n 22
  bioctl rand --int -l 23 -n 10`
}
