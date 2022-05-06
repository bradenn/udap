// Copyright (c) 2022 Braden Nicholson

package pulse

import "fmt"

func LogGlobal(format string, args ...any) {
	fmt.Println(fmt.Sprintf(format, args...))
}
