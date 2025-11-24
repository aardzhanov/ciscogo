package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/aardzhanov/ciscogo/ciscoterm"
	"github.com/aardzhanov/ciscogo/ciscoworker"
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

	myFoo := func(ctx context.Context, result ciscoworker.CiscoResult) {
		fmt.Println(">>> " + result.Host)
		if result.Error != nil {
			fmt.Println(">>> " + result.Error.Error())
		}

		for k, v := range result.Result {
			fmt.Println(">>> " + k)
			if v.Error != nil {
				fmt.Println(v.Error.Error())
			}
			for _, val := range v.Result {
				fmt.Println(val)
			}
		}
	}

	worker := ciscoworker.NewCiscoWorker(3)
	worker.StartWithCallback(ctx, myFoo)
	for _, job := range inJobs {
		worker.Execute(job)
	}
	<-ctx.Done()
}
