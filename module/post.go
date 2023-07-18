package module

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func NewCmdPost() *cobra.Command {
	var fileName string
	cmd := &cobra.Command{
		Use:   "post",
		Short: "exec post api",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			api := args[0]

			var data []byte
			if fileName != "" {
				data = ReadRequestData(fileName)
			}
			res := client.Request("POST", api, nil, data)
			fmt.Println(string(res.Body))
		},
	}
	cmd.Flags().StringVarP(&fileName, "filename", "f", "", "file name for send data(xml)")
	return cmd
}

func ReadRequestData(fileName string) []byte {
	raw_data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	data, err := xml.Marshal(raw_data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
