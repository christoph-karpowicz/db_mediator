package application

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

// Application is the main struct of the app.
type Application struct {
	CLI      *cli.App
	client   http.Client
	simulate string
}

// Init initializes the client side app.
func (a *Application) Init() {
	timeout := time.Duration(5 * time.Second)
	a.client = http.Client{Timeout: timeout}

	a.setCLI()

	err := a.CLI.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// setCLI configures the app's command line interface.
func (a *Application) setCLI() {
	a.CLI = cli.NewApp()
	a.CLI.UseShortOptionHandling = true

	a.CLI.Name = "Unifier CLI"
	a.CLI.Usage = "Database synchronization app."
	author := &cli.Author{Name: "Krzysztof Karpowicz", Email: "christoph.karpowicz@gmail.com"}
	a.CLI.Authors = append(a.CLI.Authors, author)
	a.CLI.Version = "1.0.0"

	a.CLI.Commands = []*cli.Command{
		{
			Name:  "run",
			Usage: "Start the specified synchronization.",
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

				a.runSynch(synchType, c.Args().Get(0), simulateFlag)

				return nil
			},
		},
		{
			Name:  "stop",
			Usage: "Stop the specified synchronization.",
			// Flags: []cli.Flag{
			// 	&cli.BoolFlag{
			// 		Name:    "all",
			// 		Aliases: []string{"a"},
			// 		Usage:   "Stop all synchronizations.",
			// 	},
			// },
			Action: func(c *cli.Context) error {
				// allFlag := c.Bool("all")

				a.stopSynch(c.Args().Get(0))

				return nil
			},
		},
		{
			Name:  "watch",
			Usage: "Watch tables from the specified databases.",
			Action: func(c *cli.Context) error {
				a.runWatch(c.Args().Get(0))

				return nil
			},
		},
	}
}

// runSynch prepares the parameters for a synchronization run request and invokes a GET function.
func (a *Application) runSynch(synchType string, synchName string, simulation bool) {
	paramMap := make(map[string]string)
	paramMap["type"] = synchType
	paramMap["run"] = synchName
	paramMap["simulation"] = strconv.FormatBool(simulation)

	response := a.makeGETRequest("http://localhost:8000/runSynch", paramMap)

	printRunResponse(response, synchType)
}

// stopSynch prepares the parameters for a synchronization stop request and invokes a GET function.
func (a *Application) stopSynch(synchName string) {
	paramMap := make(map[string]string)
	paramMap["stop"] = synchName
	// paramMap["all"] = strconv.FormatBool(all)

	response := a.makeGETRequest("http://localhost:8000/stopSynch", paramMap)

	printStopResponse(response)
}

// runWatch prepares the parameters for a watcher run request and invokes a GET function.
func (a *Application) runWatch(watcherName string) {
	paramMap := make(map[string]string)
	paramMap["run"] = watcherName

	response := a.makeGETRequest("http://localhost:8000/runWatch", paramMap)

	printStopResponse(response)
}

func (a *Application) makeGETRequest(url string, params map[string]string) map[string]interface{} {
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

	res, err := a.client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return parseResponse(resBody)
}
