package synch

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const (
	SIMULATION_DIR = "./simulation/"
	LOGS_DIR       = "./log/"
)

type iteration struct {
	id      string
	synch   *Synch
	actions []*action
}

func newIteration(synch *Synch) *iteration {
	return &iteration{
		id:    getNewIterationID(synch),
		synch: synch,
	}
}

func getNewIterationID(synch *Synch) string {
	return synch.cfg.Name + "-" + strconv.FormatInt(time.Now().UnixNano(), 10)
}

func (i *iteration) addAction(act *action) {
	i.actions = append(i.actions, act)
}

func (i *iteration) flush() {
	actionsToJSON, err := i.actionsToJSON()
	if err != nil {
		panic(err)
	}
	actionsToJSONString := strings.Join(actionsToJSON, "\n")
	if i.synch.IsSimulation() {
		err := ioutil.WriteFile(SIMULATION_DIR+i.id, []byte(actionsToJSONString), 0644)
		if err != nil {
			panic(err)
		}
	} else {

	}
}

func (i *iteration) actionsToJSON() ([]string, error) {
	actionsToJSON := make([]string, 0)
	for _, act := range i.actions {
		actionJSON, err := json.Marshal(&act)
		if err != nil {
			return nil, err
			// return false, &SynchReportError{SynchName: r.synch.GetConfig().Name, ErrMsg: err.Error()}
		}
		actionsToJSON = append(actionsToJSON, string(actionJSON))
	}
	return actionsToJSON, nil
}
