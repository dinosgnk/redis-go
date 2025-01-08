package main

type Command struct {
	name   []byte
	args   [][]byte
	client *Client
}

func NewCommand(cmdArgs [][]byte, client *Client) *Command {
	return &Command{
		name:   cmdArgs[0],
		args:   cmdArgs,
		client: client,
	}
}
