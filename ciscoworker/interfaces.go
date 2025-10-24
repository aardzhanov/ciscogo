package ciscoworker

type CiscoWorker interface {
	Start()
	Execute(job CiscoJobs)
	Output() chan CiscoResult
}
