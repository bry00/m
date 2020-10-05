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

const DefaultBlockSizeMB =  4
const DefaultTotalSizeMB = 64

var (
	fileName string
	title    string
	removeBackspaces bool
	blockSizeLimitMB int
	totalSizeLimitMB int
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Program %s is designated to view and browse flat, text files.\n", prog)
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
		fmt.Fprintf(os.Stderr, "Press h when browsing, to see list of available shortcuts.\n")
		fmt.Fprintf(os.Stderr, "Configuration file: %s\n", config.GetConfigFileName(prog))
		fmt.Fprintf(os.Stderr, "Copyright (C) 2020 Bartek Rybak (licensed under the MIT license).\n")
	}

	flag.StringVar(&title, "t", "", "title to show")
	flag.BoolVar(&removeBackspaces, "b", false, "remove backspaces")
	flag.IntVar(&blockSizeLimitMB, "block", DefaultBlockSizeMB, "single data block size limit (MB)")
	flag.IntVar(&totalSizeLimitMB, "total", DefaultTotalSizeMB, "total data size limit (MB)")

	flag.Parse()
	setupLogger()
}


func main() {

	if len(flag.Args()) > 0 {
		fileName = strings.Join(flag.Args(), " ")
	}

	conf := config.GetConfig(prog)

	if blockSizeLimitMB <= 0 {
		blockSizeLimitMB = conf.DataBuffer.BlockSizeLimitMB
	}
	if totalSizeLimitMB <= 0 {
		totalSizeLimitMB = conf.DataBuffer.TotalSizeLimitMB
	}

	ctl := controller.NewController(fileName, title,
		buffers.NewBufferedDataMB(blockSizeLimitMB, totalSizeLimitMB),
		tv.NewView(), conf, removeBackspaces)
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
