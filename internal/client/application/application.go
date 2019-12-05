package application

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli"
)

type Application struct {
	CLI    *cli.App
	client http.Client
	Lang   string
}

func (a *Application) Init() {
	timeout := time.Duration(5 * time.Second)
	a.client = http.Client{Timeout: timeout}
}

func (a *Application) makeGETRequest(url string, params map[string]string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	for param, val := range params {
		q.Add(param, val)
	}
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.RawQuery)

	res, err := a.client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(resBody)
}

func (a *Application) requestSynch(synchType string, synchName string) {
	paramMap := make(map[string]string)
	paramMap["type"] = synchType
	paramMap["synch"] = synchName
	a.makeGETRequest("http://localhost:8000/init", paramMap)
}

func (a *Application) SetCLI() {
	a.CLI = cli.NewApp()
	a.CLI.Name = "Unifier CLI"
	a.CLI.Usage = "Database synchronization app."
	a.CLI.Author = "Krzysztof Karpowicz"
	a.CLI.Version = "1.0.0"

	a.CLI.Commands = []cli.Command{
		{
			Name:    "one-off",
			Aliases: []string{"oo"},
			Usage:   "One off synchronization.",
			Action: func(c *cli.Context) {
				a.requestSynch("one-off", c.Args()[0])
			},
		},
		{
			Name:    "ongoing",
			Aliases: []string{"ng"},
			Usage:   "Start ongoing synchronization.",
			Action: func(c *cli.Context) {
				a.requestSynch("ongoing", c.Args()[0])
			},
		},
	}

	a.CLI.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "lang",
			Value:       "english",
			Usage:       "language for the greeting",
			Destination: &a.Lang,
		},
	}
}
