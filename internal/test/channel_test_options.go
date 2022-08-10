package test

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/channels"
)

type ChannelTestOptions struct {
	TestOptions[channels.Channel]
}

func NewChannelTestOptions() *ChannelTestOptions {
	return &ChannelTestOptions{
		TestOptions: *NewTestOptions[channels.Channel]("channel"),
	}
}
