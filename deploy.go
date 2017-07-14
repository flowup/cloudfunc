package main

import (
	"os"
	"os/exec"
	"github.com/flowup/cloudfunc/shim"
	"fmt"
	"html/template"
	"strconv"
)

// Function represents a configuration of a single function
type Function struct {
	Source      string `mapstructure:"source"`
	Name        string `mapstructure:"name"`
	StageBucket string `mapstructure:"bucket"`
	Memory      int `mapstructure:"memory"`
	Timeout     int `mapstructure:"timeout"`
}

// ToArray creates an array of function arguments that are passed into the
// deployment command
func (f *Function) ToArray() []string {
	res := []string{}

	if f.StageBucket != "" {
		res = append(res, "--stage-bucket", f.StageBucket)
	}

	if f.Memory != 0 {
		res = append(res, "--memory", strconv.Itoa(f.Memory))
	}

	if f.Timeout != 0 {
		res = append(res, "--timeout", strconv.Itoa(f.Timeout))
	}

	res = append(res, "--trigger-http")

	return res
}

// Deploy performs deployment of a given function configuration
func Deploy(function *Function) error {
	// set the folder name constant
	folderName := "__deploy__" + function.Name

	// create the deployment folder and defer it's removal
	err := os.Mkdir(folderName, os.ModePerm)
	if err != nil {
		return err
	}

	// set required environment variables for correct go build
	os.Setenv("GOOS", "linux")
	os.Setenv("CGO_ENABLED", "0")

	// build the binary of the function
	build := exec.Command("go", "build", "-a", "-installsuffix", "cgo", "-o", folderName+"/main", "./"+function.Source)
	build.Stdout = os.Stdout
	build.Stderr = os.Stdout

	fmt.Println("Building sources")
	err = build.Run()
	if err != nil {
		return err
	}

	tbytes, err := shim.Asset("shim/index.js")
	if err != nil {
		return err
	}

	shimFile, err := os.Create(folderName + "/index.js")
	if err != nil {
		return err
	}

	t := template.Must(template.New("shim").Parse(string(tbytes)))
	t.Execute(shimFile, &function)

	// change dir to the deployment folder and back
	os.Chdir(folderName)


	fmt.Println("Deploying function", function.Name, "to bucket", function.StageBucket)
	// deploy the function
	args := []string{"beta", "functions", "deploy", function.Name}
	args = append(args, function.ToArray()...)
	deploy := exec.Command("gcloud", args...)
	deploy.Stdout = os.Stdout
	deploy.Stderr = os.Stdout

	err = deploy.Run()
	if err != nil {
		return err
	}

	shimFile.Close()
	os.Chdir("..")
	return os.RemoveAll(folderName)
}
