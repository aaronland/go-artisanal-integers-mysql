package main

import (
	"github.com/aaronland/go-artisanal-integers-mysql/engine"
	"github.com/aaronland/go-artisanal-integers/application"
	"log"
	"os"
)

func main() {

	flags := application.NewServerApplicationFlags()

	var dsn string
	flags.StringVar(&dsn, "dsn", "{USER}:{PSWD}/@{DATABASE}", "The data source name (dsn) for connecting to MySQL.")

	application.ParseFlags(flags)

	eng, err := engine.NewMySQLEngine(dsn)

	if err != nil {
		log.Fatal(err)
	}

	app, err := application.NewServerApplication(eng)

	if err != nil {
		log.Fatal(err)
	}

	err = app.Run(flags)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
