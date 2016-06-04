package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// TagVersion contains tag parsing
type TagVersion struct {
	Mayor    int
	Minor    int
	Build    int
	Revision int
}

// Tag is a docker hub image tag
type Tag struct {
	Layer   string     `json:"layer"`
	Name    string     `json:"name"`
	Version TagVersion `json:"-"`
}

// ImageData has image info
type ImageData struct {
	Tags []Tag
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

// GetImageData return image tags
func (c *DockerHub) GetImageData(image string) (*ImageData, error) {

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

	img := ImageData{}
	d := json.NewDecoder(resp.Body)
	if err = d.Decode(&img.Tags); err != nil {
		return nil, err
	}

	for i := range img.Tags {
		img.Tags[i].parseVersion()
	}

	return &img, nil
}

// parseVersion parses a tag, and tries to match 1.2.3-4 as Mayor.Minor.Build-Revision
func (t *Tag) parseVersion() {
	tagchops := strings.Split(t.Name, ".")

	r := regexp.MustCompile("[^\\d]*(\\d*)?[^\\d]*(\\d*)?")

	for i := range tagchops {
		switch i {
		case 0:
			c := r.FindStringSubmatch(tagchops[i])
			v, _ := strconv.Atoi(c[1])
			t.Version.Mayor = v
		case 1:
			c := r.FindStringSubmatch(tagchops[i])
			v, _ := strconv.Atoi(c[1])
			t.Version.Minor = v
		case 3:
			c := r.FindStringSubmatch(tagchops[i])
			v, _ := strconv.Atoi(c[1])
			t.Version.Build = v
			v, _ = strconv.Atoi(c[2])
			t.Version.Revision = v
		}
	}
}

// GetLatestTag returns the latest tag
// if a tag "latest" is found it returns "latest"
func (i *ImageData) GetLatestTag() string {

	latest := Tag{
		Version: TagVersion{0, 0, 0, 0},
	}

	for _, t := range i.Tags {
		if t.Name == "latest" {
			return t.Name
		}

		// TODO refactor to a function

		if t.Version.Mayor < latest.Version.Mayor {
			continue
		}

		if t.Version.Mayor > latest.Version.Mayor {
			latest = t
			continue
		}

		// equal mayor

		if t.Version.Minor < latest.Version.Minor {
			continue
		}

		if t.Version.Mayor > latest.Version.Mayor {
			latest = t
			continue
		}

		// equal minor

		if t.Version.Build < latest.Version.Build {
			continue
		}

		if t.Version.Build > latest.Version.Build {
			latest = t
			continue
		}

		// equal build

		if t.Version.Revision < latest.Version.Revision {
			continue
		}

		if t.Version.Revision > latest.Version.Revision {
			latest = t
			continue
		}
	}

	return latest.Name

}
