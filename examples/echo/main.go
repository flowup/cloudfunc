package echo

import (
	"github.com/flowup/cloudfunc/api"
)

func main() {
	var input map[string]interface{} = make(map[string]interface{})
	api.GetInput(&input)

	api.Send(&input)
}
