package goenv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/moisespsena/go-error-wrap"
)

var versionRe, _ = regexp.Compile(`\D(\d+)`)
var versionRe2, _ = regexp.Compile(`^(\D+)\d.*$`)

type GoVersion struct {
	Title string
	//UpdatedAt time.Time
	key         string
	DownloadUrl string
}

func (v *GoVersion) Key() string {
	if v.key == "" {
		r := versionRe.FindAllString(v.Title, 4)
		for i, v := range r {
			iv, _ := strconv.Atoi(v[1:])
			r[i] = fmt.Sprintf("%05d", iv)
		}
		v.key = strings.ToLower(versionRe2.ReplaceAllString(v.Title, "$1")) + "-" + strings.Join(r, "")
	}
	return v.key
}

type GoVersions struct {
}

func (v *GoVersions) Available() ([]*GoVersion, error) {
	cmd := exec.Command("go", "version")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Start()
	if err != nil {
		return nil, errwrap.Wrap(err, "Cmd.Start: Get current go version")
	}
	err = cmd.Wait()
	if err != nil {
		return nil, errwrap.Wrap(err, "Cmd.Wait: Get current go version")
	}
	parts := strings.Split(strings.TrimSpace(out.String()), " ")
	osinfo := strings.Replace(parts[len(parts)-1], "/", "-", 1)

	r, err := http.Get("https://api.github.com/repos/golang/go/milestones?sort=title&order=desc&state=closed")
	if err != nil {
		return nil, errwrap.Wrap(err, "Request")
	}
	defer r.Body.Close()
	d := json.NewDecoder(r.Body)
	var versions []*GoVersion
	if err = d.Decode(&versions); err != nil {
		return nil, errwrap.Wrap(err, "JSON Decode")
	}

	tr := &http.Transport{
		MaxIdleConnsPerHost: 1024,
		TLSHandshakeTimeout: 0 * time.Second,
	}
	client := &http.Client{Transport: tr}

	var newVersions []*GoVersion

	for _, v := range versions {
		v.DownloadUrl = "https://dl.google.com/go/" + strings.ToLower(v.Title) + "." + osinfo + ".tar.gz"
		r, err = client.Head(v.DownloadUrl)
		if err == nil && r.StatusCode == 200 {
			newVersions = append(newVersions, v)
		}
		r.Body.Close()
	}

	versions = newVersions

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Key() > versions[j].Key()
	})

	return versions, nil
}
