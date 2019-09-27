package main

import (
	"fmt"
	"log"

	"github.com/thecsw/mira"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	confFlag = kingpin.Flag("conf", "path to conf file").Short(rune('c')).Default(".conf").ExistingFile()
)

func main() {
	kingpin.Parse()
	cred := mira.ReadCredsFromFile(*confFlag)
	fmt.Println(cred.ClientId)

	reddit, err := mira.Authenticate(&cred)
	if err != nil {
		log.Fatal(err)
	}
	_ = reddit

}
