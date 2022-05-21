package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Buzz2d0/wapptester"
)

var (
	urlFlag  = flag.String("url", "", "target url")
	exprFlag = flag.String("expr", "", "expression")
)

func main() {
	flag.Parse()

	if *urlFlag == "" || *exprFlag == "" {
		fmt.Println("Usage: wapptester [options]")
		flag.PrintDefaults()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	v, err := wapptester.Match(ctx, *urlFlag, *exprFlag)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("taget: %s\n\t `%s` => %v \n", *urlFlag, *exprFlag, v)
}
