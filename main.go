package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
	"github.com/NeowayLabs/neosearch/lib/neosearch/index"
	"github.com/jteeuwen/go-pkg-optarg"
)

func main() {
	var (
		fileOpt, dataDirOpt, databaseName string
		helpOpt, newIndex, debugOpt       bool
		err                               error
		index                             *index.Index
		batchSize                         int
	)

	optarg.Header("General options")
	optarg.Add("f", "file", "Read NeoSearch JSON database from file. (Required)", "")
	optarg.Add("c", "create", "Create new index database", false)
	optarg.Add("b", "batch-size", "Batch size", 1000)
	optarg.Add("n", "name", "Name of index database", "")
	optarg.Add("d", "data-dir", "Data directory", "")
	optarg.Add("t", "trace-debug", "Enable trace for debug", false)
	optarg.Add("h", "help", "Display this help", false)

	for opt := range optarg.Parse() {
		switch opt.ShortName {
		case "f":
			fileOpt = opt.String()
		case "b":
			batchSize = opt.Int()
		case "d":
			dataDirOpt = opt.String()
		case "n":
			databaseName = opt.String()
		case "c":
			newIndex = true
		case "t":
			debugOpt = true
		case "h":
			helpOpt = true
		}
	}

	if helpOpt {
		optarg.Usage()
		os.Exit(0)
	}

	if dataDirOpt == "" {
		dataDirOpt, _ = os.Getwd()
	}

	if fileOpt == "" {
		optarg.Usage()
		os.Exit(1)
	}

	cfg := neosearch.NewConfig()

	cfg.Option(neosearch.DataDir(dataDirOpt))
	cfg.Option(neosearch.Debug(debugOpt))

	neo := neosearch.New(cfg)

	if newIndex {
		log.Printf("Creating index %s\n", databaseName)
		index, err = neo.CreateIndex(databaseName)
	} else {
		log.Printf("Opening index %s ...\n", databaseName)
		index, err = neo.OpenIndex(databaseName)
	}

	if err != nil {
		log.Fatalf("Failed to open database '%s': %v", err)
		return
	}

	defer neo.Close()

	jsonBytes, err := ioutil.ReadFile(fileOpt)

	if err != nil {
		panic(err)
	}

	var data []map[string]interface{}

	err = json.Unmarshal(jsonBytes, &data)

	if err != nil {
		panic(err)
	}

	startTime := time.Now()

	index.Batch()
	var count int
	totalResults := len(data)

	for idx := range data {
		dataEntry := data[idx]

		if dataEntry["_id"] == nil {
			dataEntry["_id"] = idx
		}

		entryJSON, err := json.Marshal(&dataEntry)
		if err != nil {
			log.Println(err)
			return
		}

		err = index.Add(uint64(dataEntry["_id"].(int)), entryJSON)
		if err != nil {
			panic(err)
		}

		if count == batchSize {
			count = 0

			fmt.Println("Flushing batch: ", idx, " from ", totalResults)
			index.FlushBatch()
			if idx != (totalResults - 1) {
				index.Batch()
			}
		} else {
			count = count + 1
		}
	}

	index.FlushBatch()

	elapsed := time.Since(startTime)

	log.Printf("Database indexed in %v\n", elapsed)
}
