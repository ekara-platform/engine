package engine

import (
	"strconv"
	"strings"
)

type version struct {
	major int
	minor int
	micro int
	full  string
}

func CreateVersion(full string) (Version, error) {
	result := version{full: full}
	split := strings.Split(full, ".")
	if len(split) > 0 {
		major, err := strconv.Atoi(split[0])
		if err != nil {
			return nil, err
		}
		result.major = int(major)
	}
	if len(split) > 1 {
		minor, err := strconv.Atoi(split[1])
		if err != nil {
			return nil, err
		}
		result.minor = int(minor)
	}
	if len(split) > 2 {
		minor, err := strconv.Atoi(split[2])
		if err != nil {
			return nil, err
		}
		result.micro = int(minor)
	}
	return result, nil
}
