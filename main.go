package main

import (
	"fmt"
	"os"
	"time"

	"github.com/inloop/go-transport-queue/model"
	"github.com/inloop/go-transport-queue/transports"
	"github.com/urfave/cli"
)

func main() {

	// transports.NewFCMTransport("AAAAg4AwZxI:APA91bFrKi9nnIisLh-33VLGuNFqgk3_V2YaGAWOLCaywwadhOiBdSb3PkKFSFnIru00Ge67RtFfwwrmLZdhdg-ktEXcHZMK_MdLgywoWBY8KCRuoeeYqX6M1HiofzFDrJB9GrID--3K")

	app := cli.NewApp()
	app.Name = "go-transport-queue"
	app.Usage = "..."
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "transport,t",
			EnvVar: "TRANSPORT",
			Usage:  "[log|smtp|fcm]",
		},
		cli.StringFlag{
			Name:   "api-port,p",
			EnvVar: "API_PORT",
			Value:  "80",
		},
		cli.IntFlag{
			Name:   "batch-size,b",
			EnvVar: "BATCH_SIZE",
			Value:  100,
		},
		cli.DurationFlag{
			Name:   "interval,i",
			EnvVar: "INTERVAL",
			Value:  time.Second / 10,
		},
		cli.StringFlag{
			Name:   "data-path",
			EnvVar: "DATA_PATH",
			Value:  "/data",
		},
		cli.StringFlag{
			Name:   "smtp-url",
			EnvVar: "SMTP_URL",
		},
		cli.StringFlag{
			Name:   "smtp-sender",
			EnvVar: "SMTP_SENDER",
		},
		cli.StringFlag{
			Name:   "fcm-api-key",
			EnvVar: "FCM_API_KEY",
		},
	}

	app.Action = func(c *cli.Context) error {

		t, err := getTransport(c.String("transport"), c)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		q := NewQueue(c.String("data-path"), t)
		// b := NewMessageBuffer(q, c.Int("batch-size"))
		d := distributor{queue: q, transport: t}
		port := c.String("api-port")

		// b.Start(c.Duration("interval"))
		d.start(c.Duration("interval"), c.Int("batch-size"))

		srv := createHTTPServer(q, t)
		fmt.Println("transport:", c.String("transport"))
		fmt.Println("batch size:", c.Int("batch-size"))
		fmt.Println("interval:", c.Duration("interval"))
		fmt.Println("queue length:", q.Length())
		fmt.Println("listening on:", port)
		if err := srv.Run(":" + port); err != nil {
			return cli.NewExitError(err, 1)
		}

		return nil
	}

	app.Run(os.Args)
}

func getTransport(transport string, c *cli.Context) (model.Transport, error) {
	var t model.Transport

	switch transport {
	case "log":
		t = transports.NewLogTransport()
	case "smtp":
		url := c.String("smtp-url")
		sender := c.String("smtp-sender")

		if url == "" {
			return t, MissingVariableError{VariableName: "SMTP_URL"}
		}
		if sender == "" {
			return t, MissingVariableError{VariableName: "SMTP_SENDER"}
		}
		t = transports.NewSMTPTransport(url, sender)
	case "fcm":
		apiKey := c.String("fcm-api-key")

		if apiKey == "" {
			return t, MissingVariableError{VariableName: "FCM_API_KEY"}
		}
		t = transports.NewFCMTransport(apiKey)
	default:
		return t, fmt.Errorf("unkown transport type '%s' (--transport attribute)", transport)
	}
	return t, nil
}

// MissingVariableError ...
type MissingVariableError struct {
	VariableName string
}

func (e MissingVariableError) Error() string {
	return fmt.Sprintf("required variable %s not specified", e.VariableName)
}
