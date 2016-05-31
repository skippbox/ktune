package worker

import (
	"fmt"
	"math"
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
	for _, tag := range tags {
		// latest is not a version
		if tag == "latest" {
			continue
		}

		fmt.Printf("%s:%s\n", image, tag)
		if getTagValue(tag) > getTagValue(latest) {
			latest = tag
		}
	}

	return latest, nil

}

func getTagValue(tag string) int {

	tagchops := strings.Split(tag, ".")

	// get first 4 chops
	if len(tagchops) > 4 {
		tagchops = tagchops[:4]
	}

	value := 0
	for i, chop := range tagchops {

		// if an item have a subitem, remove it
		// we only get the first chunk
		chop = strings.Split(chop, "-")[0]

		// remove non numbers
		reg := regexp.MustCompile("[^\\d]")
		chop = reg.ReplaceAll([]byte(chop), "")

		val, err := strconv.Atoi(chop)
		if err != nil {
			// should be a integer number
			panic(err)
		}

		value += int(val * math.Pow(100, float64((len(tagchops)-i))))
	}

	return value
}
