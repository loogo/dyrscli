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
			log.Println("task")

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
