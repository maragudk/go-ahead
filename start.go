package ahead

import (
	"fmt"
)

type StartOptions struct{}

func Start(opts StartOptions) error {
	fmt.Println("Started.")
	return nil
}
