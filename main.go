package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/NeowayLabs/neosearch"
	"github.com/NeowayLabs/neosearch/index"
	"github.com/jteeuwen/go-pkg-optarg"
	"github.com/peterh/liner"
)

// SampleData representes data.csv
type SampleData struct {
	ID          int    `json:"id"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address,omitempty"`
}

func main() {
	var fileOpt, dataDirOpt string
	var helpOpt bool

	optarg.Add("f", "file", "Read NeoSearch JSON database from file", "")
	optarg.Add("d", "data-dir", "Data directory", "")
	optarg.Add("h", "help", "Display this help", "")

	for opt := range optarg.Parse() {
		switch opt.ShortName {
		case "f":
			fileOpt = opt.String()
		case "d":
			dataDirOpt = opt.String()
		case "m":
			homeOpt = opt.String()
		case "h":
			helpOpt = true
		}
	}

	if helpOpt {
		optarg.Usage()
		os.Exit(0)
	}

	if homeOpt == "" {
		if homeEnv := os.Getenv("HOME"); homeEnv != "" {
			homeOpt = homeEnv
		}
	}

	if dataDirOpt == "" {
		dataDirOpt, _ = os.Getwd()
	}
	var cmdline string
	var index *index.Index
	var err error

	file := "./data.json"
	indexName := "sample"

	cfg := neosearch.NewConfig()

	cfg.Option(neosearch.DataDir("/data/"))

	neo := neosearch.New(cfg)

	fmt.Println("args", os.Args)

	if len(os.Args) == 2 {
		index, err = neo.CreateIndex(indexName)
	} else {
		fmt.Println("Opening index...")
		index, err = neo.OpenIndex(indexName)
	}

	if err != nil {
		panic(err)
	}

	defer neo.Close()

	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var mapData []SampleData

	err = json.Unmarshal(jsonBytes, &mapData)

	if err != nil {
		panic(err)
	}

	index.Batch()
	var count int
	totalResults := len(mapData)
	batchSize := 100000

	for idx := range mapData {
		sampleData := mapData[idx]

		dv, err := json.Marshal(&sampleData)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = index.Add(uint64(sampleData.ID), dv)
		if err != nil {
			panic(err)
		}

		if count == batchSize {
			count = 0

			index.FlushBatch()
			if idx != (totalResults - 1) {
				index.Batch()
			}
		} else {
			count = count + 1
		}

		//		percent := (idx * 100) / totalResults

		//		if int(percent)%2 == 0 {
		//			remainder := math.Remainder(float64(idx*100), float64(totalResults))

		//			if int(remainder) == 0 {
		//			fmt.Printf("Percent: %d\tCount: %d.\n", percent, count)
		//			}
		//		}
	}

	line := liner.NewLiner()
	defer line.Close()

	// command-line here
	for {
		if cmdline, err = line.Prompt("neosearch>"); err != nil {
			if err.Error() == "EOF" {
				break
			}

			continue
		}

		line.AppendHistory(cmdline)

		if strings.ToLower(cmdline) == "quit" ||
			strings.ToLower(cmdline) == "quit;" {
			break
		}

		docs, err := index.MatchPrefix([]byte("company_name"), []byte(cmdline))

		if err != nil {
			panic(err)
		}

		fmt.Println("Found total: ", len(docs))
		fmt.Println("Results: ")

		for _, doc := range docs {
			fmt.Println(doc)
		}
	}

}