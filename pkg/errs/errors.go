package errs

import "fmt"


func Wrap(action string, err error) error {
	return fmt.Errorf("%s: %w", action, err)
}
