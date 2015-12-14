package client

import (
	"fmt"
	"io"

	Cli "github.com/docker/docker/cli"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/runcom/docker-novolume-plugin/Godeps/_workspace/src/github.com/docker/docker/api/types"
	flag "github.com/runcom/docker-novolume-plugin/Godeps/_workspace/src/github.com/docker/docker/pkg/mflag"
)

var validDrivers = map[string]bool{
	"json-file": true,
	"journald":  true,
}

// CmdLogs fetches the logs of a given container.
//
// docker logs [OPTIONS] CONTAINER
func (cli *DockerCli) CmdLogs(args ...string) error {
	cmd := Cli.Subcmd("logs", []string{"CONTAINER"}, Cli.DockerCommands["logs"].Description, true)
	follow := cmd.Bool([]string{"f", "-follow"}, false, "Follow log output")
	since := cmd.String([]string{"-since"}, "", "Show logs since timestamp")
	times := cmd.Bool([]string{"t", "-timestamps"}, false, "Show timestamps")
	tail := cmd.String([]string{"-tail"}, "all", "Number of lines to show from the end of the logs")
	cmd.Require(flag.Exact, 1)

	cmd.ParseFlags(args, true)

	name := cmd.Arg(0)

	c, err := cli.client.ContainerInspect(name)
	if err != nil {
		return err
	}

	if !validDrivers[c.HostConfig.LogConfig.Type] {
		return fmt.Errorf("\"logs\" command is supported only for \"json-file\" and \"journald\" logging drivers (got: %s)", c.HostConfig.LogConfig.Type)
	}

	options := types.ContainerLogsOptions{
		ContainerID: name,
		ShowStdout:  true,
		ShowStderr:  true,
		Since:       *since,
		Timestamps:  *times,
		Follow:      *follow,
		Tail:        *tail,
	}
	responseBody, err := cli.client.ContainerLogs(options)
	if err != nil {
		return err
	}
	defer responseBody.Close()

	if c.Config.Tty {
		_, err = io.Copy(cli.out, responseBody)
	} else {
		_, err = stdcopy.StdCopy(cli.out, cli.err, responseBody)
	}
	return err
}
