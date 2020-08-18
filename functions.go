package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func displaySearchResults(res Results) {
	for _, mod := range res.Data {
		fmt.Println("Name: " + mod.Name)
		fmt.Println("ID: " + fmt.Sprint(mod.ID))
		fmt.Println("")
	}
	fmt.Println("Run sxlman --track --id {ID} to track a mod!")
}

func untrackMods(modID string, config Config) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/629/mods?id="+modID+"&api_key="+config.APIKey, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	b, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode == 404 {
		fmt.Println("No Results Found")
		os.Exit(1)
	}
	var mod Results
	err = json.Unmarshal(b, &mod)
	parseID, _ := strconv.Atoi(modID)
	if mod.Data[0].ID == parseID {
		untrackMod(mod.Data[0], config)
	} else {
		fmt.Println("Mod Not Found!")
	}
}

func installMod(fp string, config Config) {
	user, _ := user.Current()
	if config.GamePath == "" {
		fmt.Println("You need to add your game path to the sxlman config file! Found here:", user.HomeDir+"\\Documents\\sxlman\\config.json")
		os.Exit(1)
	}
	files, err := Unzip(fp, path.Join(config.GamePath, "Mods"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Installed Files:\n", strings.Join(files, "\n"))
	fmt.Println("\nTo Uninstall, Remove the Folder From:", path.Join(config.GamePath, "Mods"))
}

func trackMods(modID string, config Config) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/629/mods?id="+modID+"&api_key="+config.APIKey, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	b, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode == 404 {
		fmt.Println("No Results Found")
		os.Exit(1)
	}
	var mod Results
	err = json.Unmarshal(b, &mod)
	parseID, _ := strconv.Atoi(modID)
	if mod.Data[0].ID == parseID {
		trackMod(mod.Data[0], config)
	} else {
		fmt.Println("Mod Not Found!")
	}
}

func searchMods(query string, config Config) Results {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/629/mods?_q="+strings.ReplaceAll(query, " ", "+")+"&api_key="+config.APIKey, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	b, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode == 404 {
		fmt.Println("No Results Found")
		os.Exit(1)
	}
	var mod Results
	err = json.Unmarshal(b, &mod)
	return mod
}

func storeModToConfig(config Config) {
	user, _ := user.Current()
	file, _ := json.MarshalIndent(config, "", " ")
	_ = ioutil.WriteFile(user.HomeDir+"\\Documents\\sxlman\\config.json", file, 0644)
}

func trackMod(mod Mods, config Config) {
	idx := -1
	for i, m := range config.TrackedPackages {
		if m == mod {
			idx = i
			break
		}
	}
	if idx == -1 {
		config.TrackedPackages = append(config.TrackedPackages, mod)
		storeModToConfig(config)
		fmt.Println("Tracking Mod: " + mod.Name)
	} else {
		fmt.Println("Already Tracking Mod!")
	}
}

func updateModsFromOld(mod Mods, config Config) {
	for idx, mods := range config.TrackedPackages {
		if mods.ID == mod.ID {
			config.TrackedPackages[idx] = mod
		}
	}
	storeModToConfig(config)
}

func listTracked(config Config) {
	fmt.Println("Tracking:")
	for _, m := range config.TrackedPackages {
		fmt.Println(m.Name)
		fmt.Println("ID:", m.ID)
		fmt.Println("Current Version:", m.Modfile.Version)
		fmt.Println("")
	}
}

func downloadMod(modID string, config Config) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/629/mods?id="+modID+"&api_key="+config.APIKey, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	b, _ := ioutil.ReadAll(res.Body)
	var modRes Results
	_ = json.Unmarshal(b, &modRes)
	out, _ := os.Create(filepath.Join(config.DownloadFolder, modRes.Data[0].Name+"-"+modRes.Data[0].Modfile.Version+".zip"))
	defer out.Close()
	resp, err := http.Get(modRes.Data[0].Modfile.DownloadInfo.URL)
	if err != nil {
		fmt.Println("Couldn't get download for:", modRes.Data[0].Name)
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Could not download new version of:", modRes.Data[0].Name)
	} else {
		fmt.Println("Downloaded Version of", modRes.Data[0].Name)
		fmt.Println("Available At:", filepath.Join(config.DownloadFolder, modRes.Data[0].Name+"-"+modRes.Data[0].Modfile.Version+".zip"))
	}
}

func checkForUpdates(config Config) {
	counter := 0
	for _, mod := range config.TrackedPackages {
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/629/mods?id="+strconv.Itoa(mod.ID)+"&api_key="+config.APIKey, nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		b, _ := ioutil.ReadAll(res.Body)
		var modRes Results
		_ = json.Unmarshal(b, &modRes)
		if modRes.Data[0].LastUpdate > mod.LastUpdate {
			out, _ := os.Create(filepath.Join(config.DownloadFolder, mod.Name+"-"+modRes.Data[0].Modfile.Version+".zip"))
			defer out.Close()
			resp, err := http.Get(modRes.Data[0].Modfile.DownloadInfo.URL)
			if err != nil {
				fmt.Println("Couldn't get download for:", mod.Name)
			}
			defer resp.Body.Close()
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				fmt.Println("Could not download new version of:", mod.Name)
			} else {
				fmt.Println("Downloaded Updated Version of", mod.Name)
				fmt.Println("Available At:", filepath.Join(config.DownloadFolder, mod.Name+"-"+modRes.Data[0].Modfile.Version+".zip"))
			}
			updateModsFromOld(modRes.Data[0], config)
		}
	}
	fmt.Println("Updated", counter, "Mods")
}

func untrackMod(mod Mods, config Config) {
	idx := -1
	for i, m := range config.TrackedPackages {
		if m == mod {
			idx = i
			break
		}
	}
	if idx >= 0 {
		name := config.TrackedPackages[idx].Name
		config.TrackedPackages[idx] = config.TrackedPackages[len(config.TrackedPackages)-1]
		config.TrackedPackages[len(config.TrackedPackages)-1] = Mods{}
		config.TrackedPackages = config.TrackedPackages[:len(config.TrackedPackages)-1]
		storeModToConfig(config)
		fmt.Println("Untracked Mod:", name)
	} else {
		fmt.Println("Mod is not being tracked!")
	}
}

// Unzip function
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
