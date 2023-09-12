package elation

import (
	"net/http"
	"strconv"
)

func tokenRequest(w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path != "/token" {
		return false
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"access_token":"foo"}`))

	return true
}

func sliceStrToInt64(in []string) []int64 {
	var out []int64

	for _, v := range in {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}

		out = append(out, i)
	}

	return out
}

func strToInt64(in string) int64 {
	i, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		return 0
	}

	return i
}
