package main

import "github.com/mu853/vcdctl/module"

func main() {
	rootCmd := module.GetCmdRoot()
	rootCmd.Execute()
}
