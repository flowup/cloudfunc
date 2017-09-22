package main

import (
	"github.com/flowup/cloudfunc/api"
	"fmt"
	"os"
	"net/http"
	"io"
	"strings"
)

type RepoConfig struct {
	Owner string `json:"owner"`
	Repo string `json:"repo"`
}

type Result struct {
	Files []string `json:"files"`
	Log []string `json:"log"`
}

/*
Testing config
{
	"owner": "flowup",
    "repo": "cloudfunc"
}
*/

var (
	globalLog = ""
)

func main() {
	config := &RepoConfig{}
	cloudFunc := api.NewCloudFunc()
	req, err := cloudFunc.GetRequest()
	if err != nil {
		cloudFunc.SendResponse("Err: " + err.Error())
		return
	}

	if err := req.BindBody(&config); err != nil {
		cloudFunc.SendResponse("Err: " + err.Error())
		return
	}

	// create the zip file
	zipFile, err := os.Create("/tmp/repo.zip")
	if err != nil {
		cloudFunc.SendResponse("Err: " + err.Error())
		return
	}
	defer zipFile.Close()

	// download the repo
	res, err := http.Get(fmt.Sprintf("https://github.com/%s/%s/archive/master.zip", config.Owner, config.Repo))
	if err != nil {
		cloudFunc.SendResponse("Err: " + err.Error())
		return
	}
	defer res.Body.Close()

	// copy downloaded zip to the file
	_, err = io.Copy(zipFile, res.Body)
	if err != nil {
		cloudFunc.SendResponse("Err: " + err.Error())
		return
	}

	os.MkdirAll("/tmp/repo/", os.ModePerm)
	files, err := Unzip("/tmp/repo.zip", "/tmp/repo")
	if err != nil {
		cloudFunc.SendResponse("Err: " + err.Error())
		return
	}

	var args []string
	for _, dirname := range allPackagesInFS("/tmp/repo/...") {
		args = append(args, dirname)
	}

	for _, dir := range args {
		lintDir(dir)
	}

	cloudFunc.SendResponse(&Result{
		Files: files,
		Log: strings.Split(globalLog, "\n"),
	})
}
