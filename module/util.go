package module

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func PrityPrint(header []string, value [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
	w.Write([]byte(strings.Join(header, "\t") + "\n"))
	for _, v := range value {
		w.Write([]byte(strings.Join(v, "\t") + "\n"))
	}
	w.Flush()
}

func LastOne(str string, spliter string) string {
	tmp := strings.Split(str, spliter)
	if len(tmp) > 0 {
		return tmp[len(tmp)-1]
	}
	return ""
}

func Fatal(v ...any) {
	fmt.Printf("error: %v\n", v...)
	os.Exit(1)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
