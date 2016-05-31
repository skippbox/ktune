package worker

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/odacremolbap/grisou/client"
)

// GetLatestTag return images used in deployments
func GetLatestTag(client *client.DockerHub, image string) (string, error) {

	tags, err := client.GetTags(image)
	if err != nil {
		return "", err
	}

	latest := ""
	latestValue := []int{0, 0, 0}

	for _, tag := range tags {

		// latest is not a version
		if tag == "latest" {
			continue
		}
		tagValue := getTagValue(tag)

		if isGreaterThan(tagValue, latestValue) {
			latest = tag
			latestValue = tagValue
		}
	}

	return latest, nil

}

func getTagValue(tag string) []int {

	tagchops := strings.Split(tag, ".")
	version := []int{0, 0, 0}

	r := regexp.MustCompile("[^\\d]*(\\d*)?[^\\d]*")
	for i := range tagchops {
		// if '-' is present get only the leftmost part
		// so '2-cross' will be 2
		// remove
		c := r.FindStringSubmatch(tagchops[i])[1]
		v, _ := strconv.Atoi(c)
		version[i] = v
	}

	return version[:3]
}

func isGreaterThan(tag1, tag2 []int) bool {
	if len(tag1) != 3 || len(tag2) != 3 {
		panic(fmt.Sprintf("Tag values slices must be 3 items long: '%v' ; '%v'", tag1, tag2))
	}

	for i := 0; i < 3; i++ {
		if tag1[i] > tag2[i] {
			return true
		}
		if tag2[i] > tag1[i] {
			return false
		}
	}

	return false

}
