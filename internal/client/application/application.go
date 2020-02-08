package application

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli"
)

type Application struct {
	CLI      *cli.App
	client   http.Client
	simulate string
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

func (a *Application) requestSynch(synchType string, synchName string, simulation bool) {
	paramMap := make(map[string]string)
	paramMap["type"] = synchType
	paramMap["synch"] = synchName
	paramMap["simulation"] = strconv.FormatBool(simulation)
	a.makeGETRequest("http://localhost:8000/init", paramMap)
}

func (a *Application) SetCLI() {
	a.CLI = cli.NewApp()
	a.CLI.UseShortOptionHandling = true

	a.CLI.Name = "Unifier cli"
	a.CLI.Usage = "Database synchronization app."
	author := &cli.Author{Name: "Krzysztof Karpowicz", Email: "christoph.karpowicz@gmail.com"}
	a.CLI.Authors = append(a.CLI.Authors, author)
	a.CLI.Version = "1.0.0"

	a.CLI.Commands = []*cli.Command{
		{
			Name:  "synch",
			Usage: "Start synchronization.",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "simulate",
					Aliases: []string{"s"},
					Usage:   "Simulate a synchronization and show what changes would be made.",
				},
				&cli.StringFlag{
					Name:    "type",
					Aliases: []string{"t"},
					Usage:   "Specify the type of synchronization.",
				},
			},
			Action: func(c *cli.Context) error {
				var synchType string
				synchTypeFlag := c.String("t")
				simulateFlag := c.Bool("simulate")

				switch true {
				case synchTypeFlag == "" || synchTypeFlag == "oo" || synchTypeFlag == "one-off":
					synchType = "one-off"
				case synchTypeFlag == "ng" || synchTypeFlag == "ongoing":
					synchType = "ongoing"
				default:
					log.Fatalln("ERROR: unknown synchronization type: " + synchTypeFlag + ".")
				}

				a.requestSynch(synchType, c.Args().Get(0), simulateFlag)

				return nil
			},
		},
	}
}
