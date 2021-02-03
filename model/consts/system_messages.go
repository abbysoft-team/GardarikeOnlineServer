package consts

import "fmt"

var (
	// Messages

	MessageCharacterAuthorized = func(name string) string {
		return fmt.Sprintf("\"%s\" enters the world!", name)
	}
)
