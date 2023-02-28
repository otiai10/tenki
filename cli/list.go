package cli

import (
	"fmt"

	"github.com/otiai10/tenki/tenki"
)

func List() error {
	areas := tenki.ListAreas()
	for _, area := range areas {
		fmt.Println(area)
	}
	return nil
}
