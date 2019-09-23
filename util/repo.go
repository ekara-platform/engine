package util

import "strings"

//repositoryFlavor returns the repository flavor, branchn tag ..., based on the
// presence of '@' into the given url
func RepositoryFlavor(url string) (string, string) {
	if strings.Contains(url, "@") {
		s := strings.Split(url, "@")
		return s[0], s[1]
	}
	return url, ""
}

