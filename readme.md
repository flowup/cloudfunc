# CloudFunc

CloudFunc is a command line tool (cli) that deploys [Google Cloud Functions](https://cloud.google.com/functions/) with ease.
It uses node.js shim to wrap the Go binary as the Cloud Functions only support Node at the moment.

> Only us-central region is supported while Google Cloud Functions are in Beta

## Installation

> Tested on Golang 1.8

Please follow [gcloud sdk](https://cloud.google.com/sdk/downloads) installation notes in case you don't have the `gcloud` command.

You will also need to update and install beta features at the moment
```
gcloud components update &&
gcloud components install beta
```

Finally download the `cloudfunc` with its SDK
```
go get -u github.com/flowup/cloudfunc/...
```

## Code Example

This code example shows a simple cloud function that sends back whatever JSON it receives.

```go
package main

import (
	"github.com/flowup/cloudfunc/api"
)

func main() {
	var input map[string]interface{} = make(map[string]interface{})
	api.GetInput(&input)

	api.Send(&input)
}
```

## Usage

Single `deploy` command is exposed by the `cloudfunc`. This allows to target a folder with `main` package that will be
deployed to cloud functions. Name of the function will be derived from the name of the folder.

> You need to also target the storage bucket that will be used to store contents of your function

```
cloudfunc deploy myfunction --bucket mybucket
```

Where:
- `myfunction` is the folder with your function
- `mybucket` is the name of your gcloud bucket

## Configuration

Additional configurations can be done using `function.json` file within your function `main` package (look at the `/examples` folder).

```
{
  "name": "myfunctionname",
  "bucket": "mybucketname",
  "memory": 128,
  "timeout": 3
}
```

In case your bucket is specified within the `function.json` file, you can simply deploy with:

```
cloudfunc deploy example # where example is your function folder
```
