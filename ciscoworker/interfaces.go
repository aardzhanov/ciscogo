package ciscoworker

import "context"

type CiscoWorker interface {
	StartWithCallback(ctx context.Context, fn func(ctx context.Context, result CiscoResult))
	Execute(job CiscoJobs)
}

func NewCiscoWorker(maxParallel int) CiscoWorker {
	return &ciscoWorker{
		maxParallel: maxParallel,
		jobs:        make(chan CiscoJobs, maxParallel),
	}
}
