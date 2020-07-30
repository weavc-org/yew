package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mogolade/yuu/internal"
)

var (
	output = flag.String("output", ".", "plugin output directory; default .")
	dir    = flag.String("dir", "", "comma-separated list of directories to build; must be set")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of builder:\n")
	fmt.Fprintf(os.Stderr, "\tbuilder -output=[output directory] -dir=[...directory]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("builder: ")

	flag.Usage = Usage
	flag.Parse()

	var dirs []string
	if len(*dir) > 0 {
		dirs = strings.Split(*dir, ",")
	} else {
		panic("no directories to build")
	}

	e := internal.BuildPlugins(*output, dirs)
	if e != nil {
		panic(e)
	}
}
