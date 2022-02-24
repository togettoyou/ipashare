package tools

import (
	"fmt"
	"os/exec"
)

type cmdClient struct {
	OnWatch func(out string)
}

func (c *cmdClient) Write(p []byte) (int, error) {
	if c.OnWatch != nil {
		c.OnWatch(string(p))
	}
	return len(p), nil
}

var CmdClient = &cmdClient{
	OnWatch: func(out string) {
		fmt.Print(out)
	},
}

func NewCmdClient(onWatch func(out string)) *cmdClient {
	return &cmdClient{
		OnWatch: onWatch,
	}
}

func (c *cmdClient) Command(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = c
	cmd.Stderr = c
	return cmd.Run()
}

func (c *cmdClient) Output(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	return cmd.Output()
}
