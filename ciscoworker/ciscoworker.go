package ciscoworker

import (
	"context"
	"errors"
	"time"

	"github.com/aardzhanov/awesomeProject3/ciscoterm"
)

func (cw *ciscoWorker) Start() {
	for i := 0; i < cw.maxParallel; i++ {
		go cw.netDevWorker()
	}
}

func (cw *ciscoWorker) Execute(job CiscoJobs) {
	cw.jobs <- job
}

func (cw *ciscoWorker) Output() chan CiscoResult {
	return cw.results
}

func (cw *ciscoWorker) netDevWorker() {
	for job := range cw.jobs {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(job.Timeout)*time.Second)
		defer cancel()

		func(ctx context.Context) {
			select {
			case <-ctxTimeout.Done():
				cw.results <- CiscoResult{
					Host:  job.CiscoDevice.Hostname,
					Error: errors.New("execution timeout"),
				}
				return
			default:
				term := ciscoterm.NewTerminal()
				err := term.Connect(job.CiscoDevice)
				defer term.Close()
				if err != nil {
					cw.results <- CiscoResult{
						Host:  job.CiscoDevice.Hostname,
						Error: err,
					}
					return
				}
				err = term.EnableTerm(job.CiscoDevice.Enable)
				if err != nil {
					cw.results <- CiscoResult{
						Host:  job.CiscoDevice.Hostname,
						Error: err,
					}
					return
				}
				err = term.DisablePagination()
				if err != nil {
					cw.results <- CiscoResult{
						Host:  job.CiscoDevice.Hostname,
						Error: err,
					}
					return
				}

				for _, command := range job.Commands {
					result, err := term.ExecuteCommand(command)
					if err != nil {
						cw.results <- CiscoResult{
							Host:  job.CiscoDevice.Hostname,
							Error: err,
						}
					}
					cw.results <- CiscoResult{
						Host:    job.CiscoDevice.Hostname,
						Command: command,
						Result:  result,
					}
				}
				return
			}

		}(ctxTimeout)

	}
}
