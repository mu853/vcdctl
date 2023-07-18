package module

import (
	"encoding/json"
	"log"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	client VcdClient
	config Config
)

func GetCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vcdctl",
		Short: "vCD command-line client",
		Long:  "vCD command-line client",
	}
	cmd.AddCommand(
		NewCmdGet(),
		NewCmdPost(),
	)

	// Get config file
	homedir := os.Getenv("HOME")
	if homedir == "" && runtime.GOOS == "windows" {
		homedir = os.Getenv("USERPROFILE")
	}
	var configfile string
	cmd.PersistentFlags().StringVarP(&configfile, "config", "c", homedir+"/.config/vcdctl.json", "path to vcdctl config file")

	// Get site
	file, err := os.ReadFile(configfile)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(file, &config)
	site, err := config.GetCurrentSite()
	if err != nil {
		log.Fatal(err)
	}

	// Get vcd client
	client = *NewVcdClient(site)

	return cmd
}
