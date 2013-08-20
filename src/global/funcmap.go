package global

import (
	"time"

	. "github.com/achun/template"
)

func init() {
	FuncsMap["now"] = func() string {
		return time.Now().Format("2006-01-02 15:04:05")
	}
	FuncsMap["time"] = func() time.Time {
		return time.Now()
	}

}
