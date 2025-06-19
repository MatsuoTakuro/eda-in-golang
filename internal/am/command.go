package am

import (
	"eda-in-golang/internal/ddd"
)

const (
	CommandHdrPrefix       = "COMMAND_"
	CommandNameHdr         = CommandHdrPrefix + "NAME"
	CommandReplyChannelHdr = CommandHdrPrefix + "REPLY_CHANNEL"
)

type Command interface {
	ddd.Command
	Destination() string
}

type command struct {
	ddd.Command
	destination string
}

var _ Command = (*command)(nil)

func NewCommand(name, destination string, payload ddd.CommandPayload, options ...ddd.CommandOption) command {
	return command{
		Command:     ddd.NewCommand(name, payload, options...),
		destination: destination,
	}
}

func (c command) Destination() string {
	return c.destination
}
