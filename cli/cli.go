package cli

import (
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
				path := "temp"
				host := "172.16.105.160:8081"

				if len(os.Args) > 2 {
					topic = os.Args[2]
				} else {
					log.Fatalln("please specified which topic to generate")
				}

				if len(host) == 0 {
					log.Fatalln("no schema host")
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
