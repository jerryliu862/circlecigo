package util

import (
	"encoding/json"
	"strconv"
)

func MapToJson(o interface{}) []byte {
	r, err := json.Marshal(o)

	if err != nil {
		return nil
	}

	return r
}

func StringSliceToMap(list []string) map[string]string {
	m := make(map[string]string)

	for _, v := range list {
		if _, ok := m[v]; !ok {
			m[v] = ""
		}
	}

	return m
}

func ContainString(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}

	return false
}

func ContainInt(list []int, d int) bool {
	for _, v := range list {
		if v == d {
			return true
		}
	}

	return false
}

func IntSliceToString(list []int) string {
	var res string

	for i, v := range list {
		if i != 0 {
			res += ", "
		}
		res += strconv.Itoa(v)
	}

	return res
}
