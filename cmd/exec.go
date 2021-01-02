package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec [CMD]",
	Args:  cobra.MaximumNArgs(1),
	Short: "exec command",
	Long:  "exec command",
	Run: func(cmd *cobra.Command, args []string) {
		listObj(globalFlags.language)
		runUnixCmd(args, globalFlags.language)
	},
}

func listObj(language string) {
	if language == "ja" {
		fmt.Println("##### ファイル #####")
	} else if language == "en" {
		fmt.Println("##### FILES #####")
	}

	cmdstr := "find $PWD | xargs ls -Fd | cut -d '/' -f 1-"
	out, err := exec.Command("sh", "-c", cmdstr).Output()
	if err != nil {
		log.Fatal(err)
	}
	slice := strings.Split(string(out), "\n")
	for i, fl := range slice {
		if fl != "" {
			fmt.Println(i, ":", fl)
		}
	}
	fmt.Printf("\n")
}

func runUnixCmd(args []string, language string) {
	c := strings.Join(args, " ")
	if language == "ja" {
		fmt.Println("##### 実行結果 #####")
		fmt.Println("実行 コマンド : ", c)
		fmt.Printf("----------------\n")
	} else if language == "en" {
		fmt.Println("##### COMMAND RESULT #####")
		fmt.Println("EXECUTE COMMAND : ", c)
		fmt.Printf("----------------\n")
	}

	out, err := exec.Command("sh", "-c", c).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}

func init() {
	rootCmd.AddCommand(execCmd)
}
