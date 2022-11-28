package utils

import (
	"fmt"
	"strings"
)

func TitleBreak(title string) {
	fmt.Printf("\n>> %s\n%s\n", title, strings.Repeat("=", 50))
}
