package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/aardzhanov/awesomeProject3/ciscoterm"
	"github.com/aardzhanov/awesomeProject3/ciscoworker"
	"golang.org/x/crypto/ssh"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	inJobs := []ciscoworker.CiscoJobs{
		{
			CiscoDevice: ciscoterm.CiscoDevice{
				Hostname:     "172.31.142.14:22",
				Username:     "user",
				Password:     "passwd123",
				Enable:       "enablepasswd",
				KeyExchanges: []string{ssh.InsecureKeyExchangeDH1SHA1},
				Timeout:      1,
			},
			Commands: []string{"show run access-list icmp", "show cpu"},
		},
		{
			CiscoDevice: ciscoterm.CiscoDevice{
				Hostname:     "172.31.142.14:22",
				Username:     "user",
				Password:     "passwd123",
				Enable:       "enablepasswd",
				KeyExchanges: []string{ssh.InsecureKeyExchangeDH1SHA1},
				Timeout:      1,
			},
			Commands: []string{"show clock"},
		},
		{
			CiscoDevice: ciscoterm.CiscoDevice{
				Hostname:     "172.31.142.14:22",
				Username:     "user",
				Password:     "passwd123",
				Enable:       "invalid",
				KeyExchanges: []string{ssh.InsecureKeyExchangeDH1SHA1},
				Timeout:      1,
			},
			Commands: []string{"show clock"},
		},
		{
			CiscoDevice: ciscoterm.CiscoDevice{
				Hostname:     "172.31.142.14:22",
				Username:     "user",
				Password:     "passwd123",
				Enable:       "enablepasswd",
				KeyExchanges: []string{ssh.InsecureKeyExchangeDH1SHA1},
				Timeout:      1,
			},
			Commands: []string{"show aaa local user"},
		},
	}

	myFoo := func(result ciscoworker.CiscoResult) {
		fmt.Println(">>> " + result.Host)
		if result.Error != nil {
			fmt.Println(">>> " + result.Error.Error())
		}
		fmt.Println(">>> " + result.Command)
		for _, val := range result.Result {
			fmt.Println(val)
		}
	}

	worker := ciscoworker.NewCiscoWorker(3)
	worker.Start()
	worker.ResultCallback(ctx, myFoo)
	for _, job := range inJobs {
		worker.Execute(job)
	}

	<-ctx.Done()
}
