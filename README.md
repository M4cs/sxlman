# sxlman
A command line tool for downloading/updating Skater XL mods from mod.io

# What is this?

This is a tool you use from your command line to keep up to date with SXL Mods. It will allow you to track mods from mod.io and then check periodically for updates and then download those updates for you. It helps keep you playing the most stable, up to date mods without needing to check mod.io constantly. Just focus on what you want and get it simply.

# Installation

There are two ways of installing sxlman.

### The Easy Way

1. Download and install [Go](https://golang.org/doc/install)

2. Run `go get github.com/M4cs/sxlman`

### The Not-Easy Way

1. Download the latest `sxlman.exe` release from [Here](https://github.com/M4cs/sxlman/releases)

2. Move the `sxlman.exe` program to somewhere on your computer.

3. In the Windows Search bar enter "Environment Variables" and select the "Edit System Environment Variables" option

4. Click "Environment Variables"

5. Under User Variables select PATH and then "Edit"

6. Add the folder you placed `sxlman.exe` in

**If this wasn't clear, Google how to add programs to your PATH**

# Getting Started

First you need a [mod.io](https://mod.io) account and API key.

Once you have an API Key which can be accessed from [here](https://mod.io/apikey/widget) you can run `sxlman` from the command line.

On first run, a config file will be created in `C:\Users\YOUR NAME\Documents\sxlman\` called `config.json`. Inside of this file you can place your API Key. You can also choose which Download folder to download mods to, and if you want AutoUpdates to be on everytime you run the program.

# Usage

### Searching for Mods

```
sxlman --search Your Search Here
```

### Tracking A Mod

```
sxlman --track --id MOD_ID
```

### Untracking A Mod

```
sxlman --untrack --id MOD_ID
```

### Listing Tracked Mods

```
sxlman --list
```

### Downloading A Specific Mod

```
sxlman --download --id MOD_ID
```

### Check for Updates Manually

```
sxlman -c
```

