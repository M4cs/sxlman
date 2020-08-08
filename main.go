package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("SXL Tracker", "Track SXL Package with Mod.io")
	trackID := parser.String("i", "id", &argparse.Options{Help: "ID to associate with. Required for everything but Search"})
	searchQuery := parser.String("s", "search", &argparse.Options{Help: "Search for Mods"})
	untrack := parser.Flag("u", "untrack", &argparse.Options{Help: "Pass to Untrack a Mod"})
	track := parser.Flag("t", "track", &argparse.Options{Help: "Pass to Untrack a Mod"})
	u := parser.Flag("c", "update", &argparse.Options{Help: "Check for Updates"})
	dl := parser.Flag("d", "download", &argparse.Options{Help: "Download a mod"})
	list := parser.Flag("l", "list", &argparse.Options{Help: "List Tracked Packages"})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	var config Config
	firstRun := false
	user, _ := user.Current()
	if _, err := os.Stat(user.HomeDir + "\\Documents\\sxlman"); os.IsNotExist(err) {
		err := os.MkdirAll(user.HomeDir+"\\Documents\\sxlman", 0777)
		if err != nil {
			fmt.Println("Could not create config directory!")
		}
		config = createNewConfig()
		file, _ := json.MarshalIndent(config, "", " ")
		firstRun = true
		_ = ioutil.WriteFile(user.HomeDir+"\\Documents\\sxlman\\config.json", file, 0644)
		fmt.Println("Please Enter API Key for mod.io in " + user.HomeDir + "\\Documents\\sxlman\\config.json" + "!")
	}
	if _, err := os.Stat(user.HomeDir + "\\Documents\\sxlman\\config.json"); os.IsNotExist(err) {
		config = createNewConfig()
		file, _ := json.MarshalIndent(config, "", " ")
		firstRun = true
		_ = ioutil.WriteFile(user.HomeDir+"\\Documents\\sxlman\\config.json", file, 0644)
		fmt.Println("Please Enter API Key for mod.io in " + user.HomeDir + "\\Documents\\sxlman\\config.json" + "!")
	}
	configFile, err := os.Open(user.HomeDir + "\\Documents\\sxlman\\config.json")
	if err != nil {
		log.Fatal("Error Opening Config File", err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		log.Fatal("Error Parsing Config File: ", err.Error())
	}
	if config.AutoUpdate && !firstRun && (len(config.TrackedPackages) > 0) {
		checkForUpdates(config)
	}
	if *list {
		listTracked(config)
	}
	if *u {
		checkForUpdates(config)
	}
	if *trackID != "" {
		if *untrack {
			untrackMods(*trackID, config)
		} else if *track {
			trackMods(*trackID, config)
		} else if *dl {
			downloadMod(*trackID, config)
		} else {
			fmt.Println("Please pass --track or --untrack with the Mod ID!")
		}
	}
	if *searchQuery != "" {
		res := searchMods(*searchQuery, config)
		displaySearchResults(res)
	}
	if (*track || *untrack || *dl) && *trackID == "" {
		fmt.Println("Please Pass An ID to Track! Use --search to find one!")
	}
}

func createNewConfig() Config {
	user, _ := user.Current()
	config := Config{}
	config.APIKey = ""
	config.DownloadFolder = user.HomeDir + "\\Downloads"
	config.TrackedPackages = []Mods{}
	config.AutoUpdate = true
	return config
}

func updateConfig(config Config) {
	user, _ := user.Current()
	file, _ := json.MarshalIndent(config, "", " ")
	_ = ioutil.WriteFile(user.HomeDir+"\\Documents\\sxlman\\config.json", file, 0644)
}
