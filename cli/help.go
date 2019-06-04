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
Generators:
    dyrscli g [topic]
Example:
    dyrscli task sink | source | all
    This will show all the sinks | sources tasks or both sink and source.
    Please, see the README in the newly created application
    Good luck!
`)
}
