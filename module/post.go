package module

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func NewCmdPost() *cobra.Command {
	var fileName string
	var header string
	cmd := &cobra.Command{
		Use:   "post",
		Short: "exec post api",
		Args:  cobra.ExactArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initClient()
		},
		Run: func(cmd *cobra.Command, args []string) {
			api := args[0]

			var data []byte
			if fileName != "" {
				data = ReadRequestData(fileName)
			}

			var header_map map[string]string
			if header != "" {
				header_map = map[string]string{}
				entry := strings.Split(header, ": ")
				header_map[entry[0]] = entry[1]
			}
			res := client.Request("POST", api, header_map, data)
			fmt.Println(string(res.Body))
		},
	}
	cmd.Flags().StringVarP(&fileName, "filename", "f", "", "file name for send data(xml)")
	cmd.Flags().StringVarP(&header, "header", "", "", "additional header (cf. \"Content-Type: application/vnd.vmware.vcloud.vm+xml\"")
	return cmd
}

func ReadRequestData(fileName string) []byte {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
