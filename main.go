package main

import (
	"github.com/urfave/cli"
	"os"
	"errors"
	"github.com/spf13/viper"
)

func main() {
	app := cli.NewApp()
	app.Name = "cloudfunc"
	app.Usage = "deploys cloud function to the google cloud platform"
	app.Version = "1.0.0"

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
				viper.SetDefault("source", source)

				// read viper config
				err := viper.ReadInConfig()
				if err != nil {
					return err
				}

				function := &Function{}
				err = viper.Unmarshal(function)
				if err != nil {
					return err
				}

				// check the stage bucket as it's required
				if function.StageBucket == "" {
					return errors.New("please specify the stage bucket for upload (see help)")
				}

				return Deploy(function)
			},
		},
	}

	app.Run(os.Args)
}
