# CloudFunc

CloudFunc is a command line tool (cli) that deploys cloud functions with ease. It uses node.js shim to wrap the Go binary
as the Cloud Functions only support Node at the moment.

> Only us-central region is supported while Google Cloud Functions are in Beta

## Installation

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

Single `deploy` command is exposed by the `cloudfunc``. This allows to target a folder with `main` package that will be
deployed to cloud functions. Name of the function will be derived from the name of the folder.

> You need to also target the storage bucket that will be used to store contents of your function

```
cloudfunc deploy myfunction --bucket mybucket
```

Where:
- `myfunction` is the folder with your function
- `mybucket` is the name of your gcloud bucket

## Configuration

Additional configurations can be done using `function.json` file within your function `main` package (look at the `/example` folder).

```
{
  "name": "myfunctionname",
  "bucket": "mybucketname",
  "memory": 128,
  "timeout": 3
}
```