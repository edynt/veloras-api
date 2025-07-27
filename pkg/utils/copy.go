package utils

import (
	"fmt"

	"github.com/jinzhu/copier"
)

func SafeCopy(dst, src any) error {
	if err := copier.Copy(dst, src); err != nil {
		return fmt.Errorf("copy error: %w", err)
	}
	return nil
}
