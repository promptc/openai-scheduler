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
	Raw      *openai.Client
	Status   Status
	Identity string
}

func New(token string) (g *Client) {
	return &Client{
		Raw:      openai.NewClient(token),
		Status:   OK,
		Identity: token,
	}
}

func NewWithIdentity(identity, token string) (g *Client) {
	return &Client{
		Raw:      openai.NewClient(token),
		Status:   OK,
		Identity: identity,
	}
}

func NewWithConfig(identity string, config openai.ClientConfig) (g *Client) {
	return &Client{
		Raw:      openai.NewClientWithConfig(config),
		Status:   OK,
		Identity: identity,
	}
}

func (c *Client) IsOk() bool {
	return c.Status == OK
}

func (c *Client) IsBanned() bool {
	return c.Status == Banned
}

func (c *Client) IsOutOfService() bool {
	return c.Status == OutOfService
}

func (c *Client) StatusAdjust(e error) {
	if e == nil {
		return
	}
	err := e.Error()
	if strings.Contains(err, "status code: 200") {
		return
	}
	if strings.Contains(err, "status code: 429") {
		if strings.Contains(err, "Your access was terminated due to violation of our policies") {
			c.Status = Banned
		} else {
			c.Status = OutOfService
		}
	} else if strings.Contains(err, "status code: 403") {
		c.Status = Banned
	} else if strings.Contains(err, "status code: 401") {
		c.Status = Banned
	} else if strings.Contains(err, "status code: 503") {
		c.Status = OutOfService
	} else {
		c.Status = OutOfService
	}
}
