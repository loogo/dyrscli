package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
)

type task struct {
	name     string
	host     string
	taskType string
	target   string
}

type connectorStatus struct {
	Name      string
	Connector struct {
		State    string
		WorkerID string `json:"worker_id"`
	}
	Tasks []struct {
		State    string
		ID       int
		WorkerID string `json:"worker_id"`
	}
	Type string
}

type taskStatus struct {
	id              int
	status          string
	connector       string
	connectorStatus string
	tskType         string
}

func (tsk *task) listTask() {
	baseURL := fmt.Sprintf("http://%s/connectors/", tsk.host)
	var taskList []taskStatus
	var result []string

	resp, err := http.Get(baseURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&result)

	for _, value := range result {
		if len(tsk.name) > 0 && !strings.Contains(value, tsk.name) {
			continue
		}
		if tsk.taskType != "all" && !strings.HasSuffix(value, "-"+tsk.taskType) {
			continue
		}
		url := fmt.Sprintf("%s%s/status", baseURL, value)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		var connectorstatus connectorStatus
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&connectorstatus)
		for _, tskstatus := range connectorstatus.Tasks {
			taskList = append(taskList, taskStatus{
				connector:       connectorstatus.Name,
				connectorStatus: connectorstatus.Connector.State,
				id:              tskstatus.ID,
				status:          tskstatus.State,
				tskType:         connectorstatus.Type,
			})
		}
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "Connector\tConnectorStatus\tTaskID\tTaskStatus\tType")
	count := 0
	for _, value := range taskList {
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n", value.connector, value.connectorStatus, value.id, value.status, value.tskType)
		count++
		if count%20 == 0 {
			w.Flush()

			fmt.Scanln() // wait for Enter Key
		}
	}
	w.Flush()
}

func (tsk *task) restart() {
	baseURL := fmt.Sprintf("http://%s/connectors", tsk.host)
	if len(tsk.name) == 0 {
		log.Fatalln("please specified connector name")
	}
	url := fmt.Sprintf("%s/%s/status", baseURL, tsk.name)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var connectorstatus connectorStatus
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&connectorstatus)
	if tsk.target == "task" {
		for _, value := range connectorstatus.Tasks {
			url := fmt.Sprintf("%s/%s/tasks/%d/restart", baseURL, tsk.name, value.ID)
			resp, err := http.Post(url, "application/json", nil)
			if err != nil {
				log.Fatalln(err)
			}
			defer resp.Body.Close()
		}
	} else {
		url := fmt.Sprintf("%s/%s/restart", baseURL, tsk.name)
		resp, err := http.Post(url, "application/json", nil)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
	}
}
