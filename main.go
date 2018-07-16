package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

const (
	exitCode = 1
)

//Attributes is input command line attributes
type Attributes struct {
	dDir     string
	nDir     string
	oDir     string
	settings string
}

func main() {
	var data Attributes

	flag.StringVar(&data.dDir, "dDir", "", "Directory with default files.")
	flag.StringVar(&data.nDir, "nDir", "", "Directory with files for compare with default files.")
	flag.StringVar(&data.oDir, "oDir", "./", "Directory for save compare results.")
	flag.StringVar(&data.settings, "config", "./config.json", "Path to configuration file.")
	flag.Parse()

	if data.dDir == "" || data.nDir == "" {
		flag.PrintDefaults()
		os.Exit(exitCode)
	}

	settings, err := GetSettings(data.settings)

	if err != nil {
		fmt.Print("Configuration can't be read\n\n")
		flag.PrintDefaults()
		os.Exit(exitCode)
	}

	defFiles, err := getDefaultFiles(data.dDir)

	if err != nil {
		fmt.Printf("Files from path %s can not be read\n\n", data.dDir)
		flag.PrintDefaults()
		os.Exit(exitCode)
	}

	chSettings, err := CheckFileSettings(defFiles, settings)

	if err != nil {
		fmt.Printf("%s\n\n", err)
		flag.PrintDefaults()
		os.Exit(exitCode)
	}

	var tDifs Differences
	var wg sync.WaitGroup

	for _, set := range chSettings {
		wg.Add(1)

		go func(data Attributes, set CompareSetting) {
			difs, err := compare(data, set)

			if err != nil {
				fmt.Printf("%s\n\n", err)
			}

			if len(difs) > 0 {
				tDifs = append(tDifs, difs...)
			}

			wg.Done()
		}(data, set)
	}
	wg.Wait()

	writeResult(tDifs, data.oDir)

	fmt.Println("Compare process complete!")
}

func getDefaultFiles(dir string) ([]string, error) {
	var files []string

	f, err := os.Open(dir)

	if err != nil {
		return files, err
	}

	fileInfo, err := f.Readdir(-1)
	f.Close()

	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}

	return files, nil
}
