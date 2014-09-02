package mydebug

import (
	"fmt"
)

func PrintDebug(log string, ptr interface{}) {
	fmt.Printf("Debug: %s; Pointer: %p\n", log, ptr)
}
