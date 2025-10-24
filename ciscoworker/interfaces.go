package ciscoworker

type CiscoWorker interface {
	Start()
	Execute(job CiscoJobs)
	Output() chan CiscoResult
}

func NewCiscoWorker(maxParallel int) CiscoWorker {
	return &ciscoWorker{
		maxParallel: maxParallel,
		jobs:        make(chan CiscoJobs, maxParallel),
		results:     make(chan CiscoResult, maxParallel),
	}
}
