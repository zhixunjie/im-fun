package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func PrettyJson(str []byte) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, str, "", " "); err != nil {
		return
	}

	fmt.Println(prettyJSON.String())
}

func Max[T int64 | uint64](a, b T) T {
	if a > b {
		return a
	}

	return b
}

func Min[T int64 | uint64](a, b T) T {
	if a < b {
		return a
	}

	return b
}
