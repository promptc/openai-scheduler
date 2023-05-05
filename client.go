package openai_scheduler

import (
	"github.com/sashabaranov/go-openai"
	"strings"
)

type Status int

const (
	OK Status = iota
	Banned
	OutOfService
	OutOfQuota
)

func (s Status) String() string {
	switch s {
	case OK:
		return "OK"
	case Banned:
		return "BANNED"
	case OutOfService:
		return "OOS"
	case OutOfQuota:
		return "OOQ"
	default:
		return "Unknown"
	}
}

type Client struct {
	Cli      *openai.Client
	Status   Status
	Identity string
}

func New(token string) (g *Client) {
	return &Client{
		Cli:      openai.NewClient(token),
		Status:   OK,
		Identity: token,
	}
}

func NewWithIdentity(identity, token string) (g *Client) {
	return &Client{
		Cli:      openai.NewClient(token),
		Status:   OK,
		Identity: identity,
	}
}

func NewWithConfig(identity string, config openai.ClientConfig) (g *Client) {
	return &Client{
		Cli:      openai.NewClientWithConfig(config),
		Status:   OK,
		Identity: identity,
	}
}

func (g *Client) IsOk() bool {
	return g.Status == OK
}

func (g *Client) IsBanned() bool {
	return g.Status == Banned
}

func (g *Client) IsOutOfService() bool {
	return g.Status == OutOfService
}

func (g *Client) StatusAdjust(e error) {
	if e == nil {
		return
	}
	err := e.Error()
	if strings.Contains(err, "status code: 200") {
		return
	}
	if strings.Contains(err, "status code: 429") {
		if strings.Contains(err, "Your access was terminated due to violation of our policies") {
			g.Status = Banned
		} else {
			g.Status = OutOfService
		}
	} else if strings.Contains(err, "status code: 403") {
		g.Status = Banned
	} else if strings.Contains(err, "status code: 401") {
		g.Status = Banned
	} else if strings.Contains(err, "status code: 503") {
		g.Status = OutOfService
	} else {
		g.Status = OutOfService
	}
}
