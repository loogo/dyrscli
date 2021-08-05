package cli

import (
	"fmt"
)

func Help() {
	fmt.Println(`
Usage:
    dyrscli COMMAND [FLAGS]
Commands:
    task        Show status of connectors and tasks 
    g           Generating schema avrc file
    help        Output this message again
    recreate    Delete and Recreate an specified Debezium connector.
                --server Debezium server address, sample: http://127.0.0.1:8083, default value: http://172.16.10.246:8083
                --connect_name Source or Sink name, default value: sample-sink
Generators:
    dyrscli g [topic]
Example:
    dyrscli task sink | source | all
    This will show all the sinks | sources tasks or both sink and source.
    Please, see the README in the newly created application
    Good luck!
`)
}
