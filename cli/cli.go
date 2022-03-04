package cli

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var topicMap = map[string]string{}

// kafkacat -b localhost:29092 -C -t my_connect_offsets -f '%p %k\n' > /opt/dev/report/offse

func Route() {
	// load source partition
	file, _ := os.Open("offset")

	if file == nil {
		_, _ = os.Create("offset")
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		txts := strings.Split(line, " ")
		var vals []string
		json.Unmarshal([]byte(txts[1]), &vals)
		topicMap[vals[0]], _ = getStringInBetween(txts[0], "Partition(", ")")
		topicMap[vals[0]+"_server"], _ = getStringInBetween(txts[1], ",", "]")
	}
	b, _ := json.MarshalIndent(topicMap, "", "  ")
	fmt.Println(string(b))
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help":
			Help()

		case "task":
			f := flag.NewFlagSet("generate", flag.ExitOnError)
			var host string
			var name string
			var mode string
			var target string
			var duration int

			if len(os.Args) > 2 {
				f.StringVar(&host, "host", "127.0.0.1:8083", "connector host with port")
				f.StringVar(&name, "name", "", "connector name")
				f.StringVar(&mode, "type", "source", "connector type")
				f.StringVar(&target, "target", "task", "which task needs to operat (task|connector|all)")
				f.IntVar(&duration, "dur", 60, "how long time is the duration, seconds")

				f.Parse(os.Args[3:])
				if len(host) == 0 {
					log.Fatalln("no connector api host specified")
				}
				t := task{
					host:     host,
					name:     name,
					taskType: mode,
					target:   target,
				}

				switch os.Args[2] {
				case "ls", "list":
					t.listTask()

				case "start", "restart":
					if target == "all" {
						for {
							taskList := t.getNotRunningTasks()
							for _, value := range taskList {
								if value.status == "RUNNING" && value.connectorStatus == "RUNNING" {
									continue
								}
								target := "task"
								if value.connectorStatus != "RUNNING" {
									target = "connector"
								}
								log.Printf("Connector: %s`s  status is invalid", value.connector)
								log.Println(value.trace)
								//Caused by: org.apache.kafka.connect.errors.DataException: Failed to serialize Avro data from topic qtesvc_sdv_product_small_fee_cfg_org_rel :
								errTopic, found := getAVROTopic(value.trace)
								if found && errTopic != "" {
									fmt.Println("found avro topic:" + errTopic)
									keyScript := fmt.Sprintf("curl -X DELETE http://localhost:8081/subjects/%s-key", errTopic)
									fmt.Println(keyScript)

									valueScript := fmt.Sprintf("curl -X DELETE http://localhost:8081/subjects/%s-value", errTopic)
									fmt.Println(valueScript)

									exec.Command("bash", "-c", keyScript).Run()
									exec.Command("bash", "-c", valueScript).Run()
								}
								if strings.Contains(value.trace, "Could not find first log file name in binary log index file") ||
									strings.Contains(value.trace, "Failed to deserialize data of EventHeaderV4") {
									if sourcePartition, ok := topicMap[value.connector]; ok {
										execScript := "echo '[\"" + value.connector + "\"," + topicMap[value.connector+"_server"] + "]|' | kafkacat -P -Z -b localhost:29092 -t my_connect_offsets -K \\| -p " + sourcePartition
										fmt.Println(execScript)
										stdout, err := exec.Command("bash", "-c", execScript).Output()
										if err != nil {
											log.Println(err)
										}
										fmt.Println(string(stdout))
									}
								}

								tt := task{
									host:     host,
									name:     value.connector,
									taskType: mode,
									target:   target,
								}
								tt.restart()
							}
							log.Println("connector status heart rate...")
							time.Sleep(time.Second * time.Duration(duration))
						}
					} else {
						t.restart()
					}

				default:
					log.Fatalln("unrecognized input")
				}
			} else {
				log.Fatalln("please input subcommand to task")
			}
		case "generate":
		case "g":
			{
				var topic string
				var path string
				var host string

				if len(os.Args) > 2 {
					topic = os.Args[2]
				} else {
					log.Fatalln("please specified which topic to generate")
				}
				f := flag.NewFlagSet("generate", flag.ExitOnError)

				f.StringVar(&path, "path", "temp", "the path where to save")
				f.StringVar(&host, "host", "", "schema register host with port")

				f.Parse(os.Args[3:])

				if len(host) == 0 {
					log.Fatalln("no schema register host input")
				}
				scaffold := &scaffold{
					Topic: topic,
					Path:  path,
					Host:  host,
				}

				scaffold.Scaffold()
			}
		case "recreate":
			var server string
			var connect_name string

			f := flag.NewFlagSet("recreate", flag.ExitOnError)
			f.StringVar(&server, "server", "http://172.16.10.246:8083", "Debezium server address, sample: http://127.0.0.1:8083")
			f.StringVar(&connect_name, "connect_name", "sample-sink", "Source or Sink name")
			f.Parse(os.Args[2:])

			Recreate(server, connect_name)
		}
	}
}

func getAVROTopic(trace string) (result string, found bool) {
	errTopic, found := getStringInBetweenTwoString(trace, "Failed to serialize Avro data from topic ", " :")
	if found {
		return errTopic, found
	} else {
		return getStringInBetweenTwoString(trace, "Failed to access Avro data from topic ", " :")
	}
}

func getStringInBetweenTwoString(str string, startS string, endS string) (result string, found bool) {
	return getStringInBetween(str, startS, endS)
}

func getStringInBetween(str string, startS string, endS string) (result string, found bool) {
	s := strings.Index(str, startS)
	if s == -1 {
		return result, false
	}
	newS := str[s+len(startS):]
	e := strings.Index(newS, endS)
	if e == -1 {
		return result, false
	}
	result = newS[:e]
	return result, true
}
