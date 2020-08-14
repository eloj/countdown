package main

import (
	"fmt"
	"net/url"
	"strconv"
)

type WordParams struct {
	Limit   int `json:"limit"`
	Maxdist int `json:"maxdist"`
}

func validateIntRange(s string, min int, max int) (int, bool) {
	val, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, false
	}
	return int(val), int(val) >= min && int(val) <= max
}

func ParseWordParams(url *url.URL, wp *WordParams) error {
	var ok bool
	for qk, qvs := range url.Query() {
		narg := len(qvs)

		// Skip keys without value(s)
		if narg == 0 {
			continue
		}

		switch qk {
		case "limit":
			if narg != 1 {
				return fmt.Errorf("limit takes one argument")
			}
			if wp.Limit, ok = validateIntRange(qvs[0], 0, 1<<16); !ok {
				return fmt.Errorf("limit out of range error")
			}

		case "maxdist":
			if narg != 1 {
				return fmt.Errorf("maxdist takes one argument")
			}
			if wp.Maxdist, ok = validateIntRange(qvs[0], 0, 31); !ok {
				return fmt.Errorf("maxdist out of range error")
			}

		default:
			return fmt.Errorf("Invalid parameter '%s'", qk)
		}

	}

	return nil
}
