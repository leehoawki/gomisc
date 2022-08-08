package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/shima-park/agollo"
	"github.com/valyala/fasttemplate"
	"io/ioutil"
	"os"
)

func main() {
	parser := argparse.NewParser("apollo script", "Apollo based template rendering Command line")
	a := parser.String("a", "apollo", &argparse.Options{Required: false, Help: "apollo url"})
	n := parser.String("n", "name", &argparse.Options{Required: true, Help: "target app name"})
	f := parser.String("f", "file", &argparse.Options{Required: true, Help: "target file"})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(-1)
	}
	apollo := os.Getenv("APOLLO_CONFIGSERVICE")
	if *a != "" {
		apollo = *a
	}

	aa, err := agollo.New(apollo, *n, agollo.AutoFetchOnCacheMiss())
	if err != nil {
		panic(err)
	}

	values := aa.GetNameSpace("application.properties")
	template, err := ioutil.ReadFile(*f)
	if err != nil {
		panic(err)
	}
	t := fasttemplate.New(string(template), "{{", "}}")
	s := t.ExecuteString(values)
	fmt.Printf("%s", s)
}
