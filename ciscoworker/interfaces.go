package ciscoworker

import "context"

type CiscoWorker interface {
	Start()
	Execute(job CiscoJobs)
	ResultCallback(ctx context.Context, fn func(result CiscoResult))
}

func NewCiscoWorker(maxParallel int) CiscoWorker {
	return &ciscoWorker{
		maxParallel: maxParallel,
		jobs:        make(chan CiscoJobs, maxParallel),
		results:     make(chan CiscoResult, maxParallel),
	}
}
