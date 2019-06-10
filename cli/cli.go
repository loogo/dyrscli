package cli

import (
	"flag"
	"log"
	"os"
)

func Route() {
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

			if len(os.Args) > 2 {
				f.StringVar(&host, "host", "", "connector host with port")
				f.StringVar(&name, "name", "", "connector name")
				f.StringVar(&mode, "type", "source", "connector type")
				f.StringVar(&target, "target", "task", "which task needs to operat (task|connector)")

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
					t.restart()

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
		}
	}
}
