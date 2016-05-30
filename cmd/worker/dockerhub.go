package worker

import "github.com/odacremolbap/grisou/client"

// GetLatestTag return images used in deployments
func GetLatestTag(client *client.DockerHub, image string) (string, error) {

	_, err := client.GetTags(image)

	if err != nil {
		return "", err
	}

	return "latest", nil

}
