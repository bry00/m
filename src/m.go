package main

import (
	"flag"
	"fmt"
	"github.com/bry00/m/buffers"
	"github.com/bry00/m/config"
	"github.com/bry00/m/controller"
	"github.com/bry00/m/view/tv"
	"log"
	"os"
	"path"
	"strings"
)

var prog string = getProg()

var (
	fileName string
	title    string
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Program %s is designated to view text files.\n", prog)
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\t%s <options> [file]\n", prog)
		fmt.Fprintf(os.Stderr, "where <options> are:\n")

		fmt.Fprintf(os.Stderr, "\t-h\thelp, shows this text\n")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stderr, "\t-%v\t%v\n", f.Name, f.Usage)
			if f.DefValue != "" {
				fmt.Fprintf(os.Stderr, "\t\tdefault: %v\n", f.DefValue)
			}
		})
		fmt.Fprintf(os.Stderr, "Configuration file: %s\n", config.GetConfigFileName(prog))
	}

	flag.StringVar(&title, "t", "", "title to show")

	flag.Parse()
	setupLogger()
}



func main() {

	if len(flag.Args()) > 0 {
		fileName = strings.Join(flag.Args(), " ")
	}

	conf := config.GetConfig(prog)

	ctl := controller.NewController(fileName, title, buffers.NewBufferedDataDefault(), tv.NewView(), conf)
	defer ctl.OnExit()
	ctl.Run()
}

func setupLogger() {
	log.SetPrefix(fmt.Sprintf("%s: ", prog))
	log.SetFlags(0)
}

func getProg() string {
	base := path.Base(os.Args[0])
	if i := strings.LastIndex(base, "."); i < 0 {
		return base
	} else {
		return base[0: i]
	}
}
