package main

import (
	"flag"
	"fmt"
	"gogen3/command/genapplication"
	"gogen3/command/gencontroller"
	"gogen3/command/gencrud"
	"gogen3/command/genentity"
	"gogen3/command/generror"
	"gogen3/command/gengateway"
	"gogen3/command/genrepository"
	"gogen3/command/genservice"
	"gogen3/command/gentest"
	"gogen3/command/genusecase"
	"gogen3/command/genvalueobject"
	"gogen3/command/genvaluestring"
	"gogen3/command/genwebapp"
)

func main() {

	commandMap := map[string]func(...string) error{
		"usecase":     genusecase.Run,
		"entity":      genentity.Run,
		"valueobject": genvalueobject.Run,
		"valuestring": genvaluestring.Run,
		"repository":  genrepository.Run,
		"service":     genservice.Run,
		"gateway":     gengateway.Run,
		"controller":  gencontroller.Run,
		"error":       generror.Run,
		"test":        gentest.Run,
		"application": genapplication.Run,
		"crud":        gencrud.Run,
		"webapp":      genwebapp.Run,
	}

	flag.Parse()
	cmd := flag.Arg(0)

	if cmd == "" {
		fmt.Printf("Try one of this command to learn how to use it\n" +
			"  gogen usecase\n" +
			"  gogen test\n" +
			"  gogen repository\n" +
			"  gogen service\n" +
			"  gogen gateway\n" +
			"  gogen controller\n" +
			"  gogen registry\n" +
			"  gogen entity\n" +
			"  gogen genvaluestring\n" +
			"  gogen valueobject\n" +
			"  gogen enum\n" +
			"  gogen error\n" +
			"  gogen crud\n" +
			"  gogen webapp\n")
		return
	}

	var values = make([]string, 0)
	if flag.NArg() > 1 {
		values = flag.Args()[1:]
	}

	f, exists := commandMap[cmd]
	if !exists {
		fmt.Printf("Command %s is not recognized\n", cmd)
		return
	}
	err := f(values...)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

}
