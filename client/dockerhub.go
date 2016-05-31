package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

// Tag is a docker hub image tag
type Tag struct {
	Layer string `json:"layer"`
	Name  string `json:"name"`
}

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

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("URL '%s' returned status %d", u.String(), resp.StatusCode)
	}

	tags := []Tag{}
	d := json.NewDecoder(resp.Body)
	if err = d.Decode(&tags); err != nil {
		return nil, err
	}

	t := make([]string, len(tags))
	for _, tag := range tags {
		t = append(t, tag.Name)
	}
	return t, nil
}
