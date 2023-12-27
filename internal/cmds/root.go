package cmds

import "github.com/jimmykodes/gommand"

var rootCmd = &gommand.Command{
	Name: "root",
}

func Cmd() *gommand.Command {
	rootCmd.SubCommand(buildCmd)

	return rootCmd
}
