package openai_scheduler

import (
	"fmt"
	"time"
)

type Scheduler struct {
	workableCh     chan *Client
	brokenCh       chan *Client
	outOfServiceCh chan *Client
	startedDaemon  bool
	stopCh         chan bool
}

func NewScheduler(tokens []string) *Scheduler {
	sche := &Scheduler{
		startedDaemon: false,
		stopCh:        make(chan bool),
	}
	var gpts []*Client
	for _, token := range tokens {
		gpts = append(gpts, NewClient(token))
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
	//sche.statistic()
	return sche
}

func NewSchedulerFromClients(clients []*Client) *Scheduler {
	sche := &Scheduler{
		startedDaemon: false,
	}
	sche.workableCh = make(chan *Client, len(clients))
	sche.brokenCh = make(chan *Client, len(clients))
	sche.outOfServiceCh = make(chan *Client, len(clients))
	for _, gpt := range clients {
		sche.workableCh <- gpt
	}
	//sche.statistic()
	return sche
}

func (s *Scheduler) StartDaemon() {
	if s.startedDaemon {
		return
	}
	go s.daemon()
}

func (s *Scheduler) Dispose() {
	if s.startedDaemon {
		s.stopCh <- true
	}
	close(s.workableCh)
	close(s.brokenCh)
	close(s.outOfServiceCh)
}

func (s *Scheduler) DaemonStarted() bool {
	return s.startedDaemon
}

func (s *Scheduler) GetWorkable() int {
	return len(s.workableCh)
}

func (s *Scheduler) daemon() {
	select {
	case <-s.stopCh:
		return
	default:
		gpt := <-s.outOfServiceCh
		time.Sleep(time.Minute)
		gpt.Status = OK
		s.workableCh <- gpt
		fmt.Println("[GPT SCHEDULER]", gpt.Identity, "-> OK")
		s.statistic()

	}

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
