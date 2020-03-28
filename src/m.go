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
	"path/filepath"
)



func main() {
	setupLogger()
	flag.Parse()

	var fileName string

	if len(flag.Args()) > 0 {
		fileName = flag.Arg(0)
	}

	conf := config.NewDefaultConfig()

	ctl := controller.NewController(fileName, buffers.NewBufferedDataDefault(), tv.NewView(), conf)
	defer ctl.OnExit()
	ctl.Run()

	//fileName := flag.Arg(0)
	//file, err := os.Open(fileName)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//defer file.Close()
	//data := buffers.NewBufferedDataDefault()
	//
	////scanner := bufio.NewScanner(os.Stdin)
	//scanner := bufio.NewScanner(file)
	//for scanner.Scan() {
	//	data.AddLine(scanner.Text())
	//}
	//
	//if err := scanner.Err(); err != nil {
	//	fmt.Fprintln(os.Stderr, err.Error())
	//}
	//
	////fmt.Fprintln(os.Stderr, data.Len())
	////
	////i := data.NewLineIndexer()
	////for k := 0; k < 5; k++ {
	////	i.IndexBegin()
	////	for j:=0 ; i.IndexOK() && j < 5; i.IndexIncrement()  {
	////		line, _ := i.GetLine()
	////		fmt.Println(line)
	////		j++
	////	}
	////	i.IndexSet(data.Len() - 5, false)
	////	for ; i.IndexOK() ; i.IndexIncrement()  {
	////		line, _ := i.GetLine()
	////		fmt.Println(line)
	////	}
	////
	////}
	////
	////
	////
	//data.Close()

}

func setupLogger() {
	log.SetPrefix(fmt.Sprintf("%s: ", filepath.Base(os.Args[0])))
	log.SetFlags(0)
}
