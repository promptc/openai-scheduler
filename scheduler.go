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
	var clis []*Client
	for _, token := range tokens {
		clis = append(clis, NewClient(token))
	}
	return NewSchedulerFromClients(clis)
}

func NewSchedulerFromClients(clients []*Client) *Scheduler {
	if len(clients) == 0 {
		panic("no gpt token")
	}
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
	s.startedDaemon = true
	go s.daemon()
	s.Statistic()
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
	for {
		select {
		case <-s.stopCh:
			return
		case gpt := <-s.outOfServiceCh:
			time.Sleep(time.Minute)
			gpt.Status = OK
			s.workableCh <- gpt
			fmt.Println("[GPT SCHEDULER]", gpt.Identity, "-> OK")
			s.Statistic()
		}
	}
}

func (s *Scheduler) Statistic() {
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
	s.Statistic()
	return s.GetClient()
}
