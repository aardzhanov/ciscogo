package ciscoworker

import (
	"context"
	"errors"
	"time"

	"github.com/aardzhanov/awesomeProject3/ciscoterm"
)

func (wrk *ciscoWorker) StartWithCallback(ctx context.Context, fn func(ctx context.Context, result CiscoResult)) {
	for i := 0; i < wrk.maxParallel; i++ {
		go wrk.netDevWorker(ctx, fn)
	}
}

func (cw *ciscoWorker) Execute(job CiscoJobs) {
	cw.jobs <- job
}

func (wrk *ciscoWorker) netDevWorker(ctx context.Context, fn func(ctx context.Context, result CiscoResult)) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-wrk.jobs:
			func(ctx context.Context) {
				ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(job.Timeout)*time.Second)
				defer cancel()
				result := CiscoResult{
					Host: job.CiscoDevice.Hostname,
				}

				select {
				case <-ctxTimeout.Done():
					result.Error = errors.New("execution timeout")
					fn(ctx, result)
					return
				default:
					term := ciscoterm.NewTerminal()
					err := term.Connect(job.CiscoDevice)
					defer term.Close()
					if err != nil {
						result.Error = err
						fn(ctx, result)
						return
					}
					err = term.EnableTerm(ctxTimeout, job.CiscoDevice.Enable)
					if err != nil {
						result.Error = err
						fn(ctx, result)
						return
					}
					err = term.DisablePagination(ctxTimeout)
					if err != nil {
						result.Error = err
						fn(ctx, result)
						return
					}

					resultMap := make(map[string]commandResult)
					for _, command := range job.Commands {
						cmdResult, err := term.ExecuteCommand(ctxTimeout, command)
						resultMap[command] = commandResult{
							Result: cmdResult,
							Error:  err,
						}
					}
					result.Result = resultMap
					fn(ctx, result)
					return
				}
			}(ctx)
		}
	}
}
