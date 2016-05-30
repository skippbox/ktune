package client

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

// DockerHub manages Docker hub images
type DockerHub struct {
	Client  *http.Client
	RepoURL string
}

// NewDockerHubClient creates a default docker hub client
func NewDockerHubClient() *DockerHub {

	return &DockerHub{
		Client:  &http.Client{},
		RepoURL: "https://index.docker.io/v1/repositories/",
	}

}

// GetTags return image tags
func (c *DockerHub) GetTags(image string) ([]string, error) {

	u, err := url.Parse(c.RepoURL)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, image, "tags")

	r := &http.Request{
		Method: "GET",
		URL:    u,
		Header: http.Header{"Accept": {"application/json"}},
	}
	//r.Header.Add("Accept", "application/json")

	resp, err := c.Client.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
