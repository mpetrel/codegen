package str

import (
	"fmt"
	"strings"
)

func LowerFirst(rawStr string) string {
	return fmt.Sprintf("%s%s", strings.ToLower(rawStr[:1]), rawStr[1:])
}
