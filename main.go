package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/Jeffail/gabs/v2"
)

const (
	appConfigPath = "/data/bbriot/appconfig/appconfig.json"
)

func main() {
	for {
		fmt.Println("This is the Barbara OS app log sample")
		if fileExists(appConfigPath) {
			appConfigLogString, appConfigLogStringErr := getStringAppConfig("log")
			if appConfigLogStringErr != nil {
				fmt.Println("Appconfig file exists, but failed to get log string ('log' key)")
			} else {
				fmt.Println("This is the appConfig log string: " + appConfigLogString)
			}
		} else {
			fmt.Println("There is no appconfig inside Barbara OS, update it from Barbara Panel")
		}
		time.Sleep(5 * time.Second)
	}
}

func fileExists(filePath string) bool {
	_, statErr := os.Stat(filePath)
	if os.IsNotExist(statErr) {
		return false
	}

	return true
}

func parseJSON(jsonPath string) (*gabs.Container, error) {
	jsonFile, jsonFileErr := os.Open(jsonPath)
	if jsonFileErr != nil {
		return nil, errors.New("parseJSON: " + jsonFileErr.Error())
	}
	defer func() {
		jsonFileErr := jsonFile.Close()
		if jsonFileErr != nil {
			return
		}
	}()
	jsonBytes, jsonBytesErr := ioutil.ReadAll(jsonFile)
	if jsonBytesErr != nil {
		return nil, errors.New("failed to parse json file: " + jsonPath + ". Error: " + jsonBytesErr.Error())
	}

	return gabs.ParseJSON(jsonBytes)
}

func getAppConfig(configKey string) (interface{}, error) {
	appConfigMap, appConfigMapErr := parseJSON(appConfigPath)
	if appConfigMapErr != nil {
		return nil, appConfigMapErr
	}

	for appConfigKey, appConfigValue := range appConfigMap.ChildrenMap() {
		if configKey == appConfigKey {
			return appConfigValue.Data(), nil
		}
	}
	return nil, errors.New("No " + configKey + " app config key found")
}

func getStringAppConfig(configKey string) (string, error) {
	getStringAppConfigValue, getStringAppConfigValueErr := getAppConfig(configKey)
	if getStringAppConfigValueErr != nil {
		return "", getStringAppConfigValueErr
	}
	appConfigStringValue, appConfigStringValueErr := interfaceToString(getStringAppConfigValue)
	if appConfigStringValueErr != nil {
		return "", appConfigStringValueErr
	}
	return appConfigStringValue, nil
}

func interfaceToString(providedInterface interface{}) (string, error) {
	if providedInterface != nil {
		switch providedInterface.(type) {
		case string:
			return fmt.Sprintf("%v", providedInterface), nil
		default:
			return "", errors.New("not a string interface")
		}
	}
	return "", errors.New("nil interface")
}
