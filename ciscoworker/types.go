package ciscoworker

import "github.com/aardzhanov/awesomeProject3/ciscoterm"

type CiscoJobs struct {
	ciscoterm.CiscoDevice
	Commands []string
}

type commandResult struct {
	Result []string
	Error  error
}
type CiscoResult struct {
	Host   string
	Result map[string]commandResult
	Error  error
}
type ciscoWorker struct {
	maxParallel int
	jobs        chan CiscoJobs
}
