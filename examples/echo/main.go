package main

import (
	"github.com/flowup/cloudfunc/api"
)

func main() {
	var input map[string]interface{} = make(map[string]interface{})

	cloudFunction := api.NewCloudFunc()
	req, err := cloudFunction.GetRequest()
	if err != nil {
		panic(err)
	}

	err = req.BindBody(&input)
	if err != nil {
		panic(err)
	}

	cloudFunction.SendResponse(&input)
}
