package cmd

import (
	"fmt"
	"github.com/convee/go-large-files/pkgs"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "文件导入",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		importFiles()
	},
}

func importFiles() {
	var (
		total int64
	)
	var files []string
	for _, file := range files {
		cnt, err := importFile2db(file)
		if err != nil {
			continue
		}
		total += cnt
	}
	fmt.Println("file imported, total:", total)
}

// importFile2db 文件导入到db
func importFile2db(file string) (int64, error) {
	var f pkgs.FileReader
	cnt, err := f.ReadConcurrentWithSkip(file, true, root.threads, 0, func(i int64, s string) error {
		decodeData, err := decode(i, s)
		if err != nil {
			return err
		}
		// todo DB 操作
		return CheckAndReport("checkAndReport:", file, decodeData)
	})
	if err != nil {
		return cnt, err
	}
	return cnt, nil
}

// decode 解析文件行
func decode(i int64, line string) (s []string, err error) {
	s = strings.Split(line, ",")
	if len(s) != 2 {
		err = errors.New("Illegal data")
		return
	}
	return
}
