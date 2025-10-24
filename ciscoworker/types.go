package ciscoworker

import "github.com/aardzhanov/awesomeProject3/ciscoterm"

type CiscoJobs struct {
	ciscoterm.CiscoDevice
	Commands []string
}

type CiscoResult struct {
	Host    string
	Command string
	Result  []string
	Error   error
}

type ciscoWorker struct {
	maxParallel int
	jobs        chan CiscoJobs
	results     chan CiscoResult
}
