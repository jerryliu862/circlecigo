package util

import "strings"

func AdmitEmailDomain(domains []string, email string) bool {
	e := strings.Split(email, "@")
	if len(e) != 2 {
		return false
	}

	return ContainString(domains, e[1])
}
