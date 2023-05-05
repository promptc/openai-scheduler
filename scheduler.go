package openai_scheduler

import (
	"fmt"
	"time"
)

type Scheduler struct {
	workableCh     chan *Client
	brokenCh       chan *Client
	outOfServiceCh chan *Client
}

func Init(tokens []string) *Scheduler {
	sche := &Scheduler{}
	var gpts []*Client
	for _, token := range tokens {
		gpts = append(gpts, New(token))
	}
	if len(gpts) == 0 {
		panic("no gpt3 token")
	}
	sche.workableCh = make(chan *Client, len(gpts))
	sche.brokenCh = make(chan *Client, len(gpts))
	sche.outOfServiceCh = make(chan *Client, len(gpts))
	for _, gpt := range gpts {
		sche.workableCh <- gpt
	}
	sche.statistic()
	go sche.daemon()
	return sche
}

func (s *Scheduler) daemon() {
	gpt := <-s.outOfServiceCh
	time.Sleep(time.Minute)
	gpt.Status = OK
	s.workableCh <- gpt
	fmt.Println("[GPT SCHEDULER]", gpt.Identity, "-> OK")
	s.statistic()
}

func (s *Scheduler) statistic() {
	fmt.Println("[GPT STATISTICS]",
		"Valid:", len(s.workableCh),
		"Broken:", len(s.brokenCh),
		"OOS:", len(s.outOfServiceCh))
	if len(s.workableCh) == 0 {
		fmt.Println("[GPT SCHEDULER] NO GPT Available!")
	}
}

func (s *Scheduler) GetClient() *Client {
	gpt := <-s.workableCh
	if gpt.IsOk() {
		s.workableCh <- gpt
		return gpt
	}
	if gpt.IsBanned() {
		fmt.Println("[GPT SCHEDULER]", gpt.Identity, "-> BANNED")
		s.brokenCh <- gpt
	} else if gpt.IsOutOfService() {
		fmt.Println("[GPT SCHEDULER]", gpt.Identity, "-> OOS")
		s.outOfServiceCh <- gpt
	} else {
		fmt.Println("[GPT SCHEDULER]", gpt.Identity, "-> OOS")
		s.outOfServiceCh <- gpt
	}
	s.statistic()
	return s.GetClient()
}
