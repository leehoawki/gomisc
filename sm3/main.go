package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/tjfoc/gmsm/sm3"
	"io/ioutil"
	"os"
)

func main() {
	parser := argparse.NewParser("smgm", "SMGM Command line")
	f := parser.String("f", "file", &argparse.Options{Required: true, Help: "input file"})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(-1)
	}

	data, err := ioutil.ReadFile(*f)
	if err != nil {
		panic(err)
	}
	h := sm3.New()
	h.Write(data)
	sum := h.Sum(nil)
	fmt.Printf("%x\n", sum)
}
