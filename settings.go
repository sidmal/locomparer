package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

//CompareSetting is settings for compare excel rows
type CompareSetting struct {
	File          string `json:"file"`
	ColumnDefault string `json:"column_default"`
	ColumnNew     string `json:"column_new"`
}

//Settings is global settings object
type Settings struct {
	Compare []CompareSetting `json:"compare"`
}

//GetSettings is return application settings
func GetSettings(configPath string) (Settings, error) {
	var config Settings

	ioConfig, err := ioutil.ReadFile(configPath)

	if err != nil {
		return config, err
	}

	err = json.Unmarshal(ioConfig, &config)

	if err != nil {
		return config, err
	}

	return config, nil
}

//CheckFileSettings is checking and excluding files if setting not exist for it
func CheckFileSettings(files []string, settings Settings) ([]CompareSetting, error) {
	var checkedSettingsList []CompareSetting

	scm := make(map[string]CompareSetting)

	for _, s := range settings.Compare {
		scm[s.File] = s
	}

	for _, fn := range files {
		if s, ok := scm[fn]; ok {
			checkedSettingsList = append(checkedSettingsList, s)
			continue
		}

		fmt.Printf("Settings for file %s not found. File processing will be skip.", fn)
	}

	if len(checkedSettingsList) <= 0 {
		return checkedSettingsList, errors.New("settings for files not found. processing reject")
	}

	return checkedSettingsList, nil
}
