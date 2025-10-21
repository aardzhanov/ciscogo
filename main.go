package main

import (
	"fmt"

	"github.com/aardzhanov/awesomeProject3/ciscoterm"
	"github.com/aardzhanov/awesomeProject3/ciscoworker"
	"golang.org/x/crypto/ssh"
)

func main() {

	inJobs := []ciscoworker.CiscoJobs{
		{
			CiscoDevice: ciscoterm.CiscoDevice{
				Hostname:     "172.31.142.14:22",
				Username:     "user",
				Password:     "passwd123invalid",
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

	worker := ciscoworker.NewCiscoWorker(3)
	worker.Start()
	for _, job := range inJobs {
		worker.Execute(job)
	}

	for {
		select {
		case res := <-worker.Output():
			fmt.Println(">>> " + res.Host)
			if res.Error != nil {
				fmt.Println(">>> " + res.Error.Error())
			}
			fmt.Println(">>> " + res.Command)
			for _, val := range res.Result {
				fmt.Println(val)
			}

		}
	}
}
