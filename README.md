# OpenAI Scheduler

This is a simple scheduler for OpenAI Client in Go.

It contains auto detected OpenAI client status and handle the wrong status.

## Usage

```go
package main

import (
    "github.com/promptc/openai-scheduler"
)

var scheduler *openai_scheduler.Scheduler

func main() {
	tokens := []string{"token1", "token2"}
	scheduler = openai_scheduler.NewScheduler(tokens)
	scheduler.StartDaemon()
	// Some codes
}

func feedPrompt() {
	gpt := scheduler.GetClient()
	// Do codes just like do on *openai.Client
}
```