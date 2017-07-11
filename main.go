package main

import (
	"github.com/urfave/cli"
	"os"
	"errors"
	"os/exec"
	"text/template"
	"github.com/flowup/cloudfunc/shim"
	"github.com/spf13/viper"
	"fmt"
	"strconv"
)

type Function struct {
	Name        string `mapstructure:"name"`
	StageBucket string `mapstructure:"bucket"`
	Memory      int `mapstructure:"memory"`
	Timeout     int `mapstructure:"timeout"`
}

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

func main() {
	app := cli.NewApp()
	app.Name = "cloudfunc"
	app.Usage = "deploys cloud function to the google cloud platform"

	app.Commands = []cli.Command{
		{
			Name:    "deploy",
			Aliases: []string{"d"},
			Usage:   "deploys given function to current gcloud account",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "bucket, b",
					Value: "",
					Usage: "sets the bucket name for upload",
				},
			},
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() < 1 {
					return errors.New("please specify the function folder")
				}

				// get name of the function
				source := ctx.Args().First()

				// initialize viper configuration
				viper.AddConfigPath(source) // add the function folder
				viper.SetConfigName("function")
				viper.SetDefault("name", source)
				viper.SetDefault("bucket", ctx.String("bucket"))

				// read viper config
				err := viper.ReadInConfig()
				if err != nil {
					fmt.Println(err)
				}

				function := Function{}
				viper.Unmarshal(&function)

				// check the stage bucket as it's required
				if function.StageBucket == "" {
					return errors.New("please specify the stage bucket for upload (see help)")
				}

				// set the folder name constant
				folderName := "__deploy__" + function.Name

				// create the deployment folder and defer it's removal
				err = os.Mkdir(folderName, os.ModePerm)
				if err != nil {
					return err
				}

				// set required environment variables for correct go build
				os.Setenv("GOOS", "linux")
				os.Setenv("CGO_ENABLED", "0")

				// build the binary of the function
				build := exec.Command("go", "build", "-a", "-installsuffix", "cgo", "-o", folderName+"/main", "./"+source)
				build.Stdout = os.Stdout
				build.Stderr = os.Stdout

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

				// deploy the function
				args := []string{"beta", "functions", "deploy", function.Name}
				args = append(args, function.ToArray()...)
				fmt.Println(args)
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
			},
		},
	}

	app.Run(os.Args)
}
