package main

import (
	"flag"
	"github.com/g3co/go-swagger-merger/helpers"
)

func main() {
	merger := helpers.NewMerger()

	var outputFileName string

	flag.StringVar(&outputFileName, "o", "swag.yaml", "")
	flag.Parse()

	for _, f := range flag.Args() {
		err := merger.AddFile(f)
		if err != nil {
			panic(err)
		}
	}

	err := merger.Save(outputFileName)
	if err != nil {
		panic(err)
	}
}



