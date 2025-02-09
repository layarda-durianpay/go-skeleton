package main

import (
	"fmt"
	"os"

	"github.com/layarda-durianpay/go-skeleton/internal/server"
	"github.com/urfave/cli/v2"
)

func main() {
	if err := NewServiceApp().Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func NewServiceApp() (cliApp *cli.App) {
	cliApp = cli.NewApp()
	cliApp.Name = "DurianPay Disbursement Service" // can update with env
	cliApp.Version = "1.0.0"                       // can update with env

	cliApp.Before = func(context *cli.Context) error {
		return server.Init()
	}

	cliApp.Commands = append(cliApp.Commands,
		startServerCommand(),
		startConsumerCommand(),
	)

	return
}

func startServerCommand() (cmd *cli.Command) {
	cmd = &cli.Command{
		Name:  "start",
		Usage: "start server",
		Action: func(c *cli.Context) error {
			fmt.Println("acction start")
			return server.Start()
		},
	}

	return
}

func startConsumerCommand() (cmd *cli.Command) {
	cmd = &cli.Command{
		Name:  "consumer",
		Usage: "start consumer",
		Action: func(c *cli.Context) error {
			fmt.Println("acction start readers")
			return server.StartReaders()
		},
	}

	return
}
