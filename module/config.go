package module

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmdConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "config setting",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initConfig()
		},
	}
	cmd.AddCommand(
		NewCmdConfigGetSites(),
		NewCmdConfigSetSite(),
	)
	return cmd
}

func NewCmdConfigSetSite() *cobra.Command {
	var endpoint string
	var user string
	var password string
	var orgname string

	cmd := &cobra.Command{
		Use:   "set-site ${SITE_NAME}",
		Short: "add vcd site configuration",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			site, err := config.GetSite(name)
			if err != nil {
				site = Site{}
				site.Name = name
				site.Endpoint = endpoint
				site.User = user
				site.SetPassword(password)
				site.OrgName = orgname
				config.Sites = append(config.Sites, site)
				if len(config.Sites) == 1 {
					config.CurrentSite = site.Name
				}
			}

			saveConfig()
		},
	}
	cmd.Flags().StringVarP(&endpoint, "endpoint", "e", "", "endpoint for the new site (https://{vcdmanager})")
	cmd.Flags().StringVarP(&user, "user", "u", "", "user for the new site")
	cmd.Flags().StringVarP(&password, "password", "p", "", "password for the new site user")
	cmd.Flags().StringVarP(&orgname, "orgname", "o", "", "default org name")
	cmd.MarkFlagRequired("endpoint")
	cmd.MarkFlagRequired("user")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("orgname")

	return cmd
}

func NewCmdConfigGetSites() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-sites",
		Short: "show vcd site configurations",
		Run: func(cmd *cobra.Command, args []string) {
			header := []string{"Current", "Name", "Endpoint", "User", "Org"}

			var data [][]string
			for _, s := range config.Sites {
				current := ""
				if s.Name == config.CurrentSite {
					current = "*"
				}
				data = append(data, []string{current, s.Name, s.Endpoint, s.User, s.OrgName})
			}

			PrityPrint(header, data)
		},
	}
	return cmd
}

func initConfig() {
	if configFilePath == "" {
		configFilePath = defaultConfigFilePath()
	}
	viper.SetConfigFile(configFilePath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		config = Config{}
		return
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
}

func saveConfig() {
	if configFilePath == "" {
		configFilePath = defaultConfigFilePath()
	}

	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(configFilePath, file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func defaultConfigFilePath() string {
	homedir := os.Getenv("HOME")
	if homedir == "" && runtime.GOOS == "windows" {
		homedir = os.Getenv("USERPROFILE")
	}
	return homedir + "/.config/vcdctl.json"
}

type Config struct {
	CurrentSite string `json:"current-site" mapstructure:"current-site"`
	Sites       []Site `json:"sites"`
}

type Site struct {
	Name       string `json:"name"`
	Endpoint   string `json:"endpoint"`
	User       string `json:"user"`
	Password   string `json:"password"`
	OrgName    string `json:"orgname"`
	ApiVersion string `json:"apiversion"`
}

func (c *Config) GetCurrentSite() (Site, error) {
	for _, s := range c.Sites {
		if s.Name == c.CurrentSite {
			return s, nil
		}
	}
	return Site{}, fmt.Errorf("site '%s' not found", c.CurrentSite)
}

func (c *Config) GetSite(name string) (Site, error) {
	for _, s := range c.Sites {
		if s.Name == name {
			return s, nil
		}
	}
	return Site{}, fmt.Errorf("site '%s' not found", name)
}

func (s *Site) GetCredential() string {
	passwordText, err := base64.StdEncoding.DecodeString(s.Password)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString([]byte(s.User + ":" + string(passwordText)))
}

func (s *Site) SetPassword(password string) {
	s.Password = base64.StdEncoding.EncodeToString([]byte(password))
}
