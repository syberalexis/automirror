package pullers

import (
	"automirror/configs"
	"automirror/utils"
	"encoding/xml"
	"fmt"
	"github.com/BurntSushi/toml"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Object data structure
type Maven struct {
	Url              string
	Folder           string
	MetadataFileName string     `toml:"metadata_file_name"`
	POMFile          string     `toml:"pom_file"`
	DatabaseFile     string     `toml:"database_file"`
	Artifacts        []Artifact `toml:"artifact"`
}

type Artifact struct {
	Group          string
	Id             string
	MinimumVersion string `toml:"minimum_version"`
}

func BuildMaven(pullerConfig configs.PullerConfig) (Puller, error) {
	var config Maven
	tomlFile, err := ioutil.ReadFile(pullerConfig.Config)
	if err != nil {
		return nil, err
	}
	if _, err := toml.Decode(string(tomlFile), &config); err != nil {
		return nil, err
	}

	config.Url = pullerConfig.Source
	config.Folder = pullerConfig.Destination
	return config, nil
}

// Inherits public method to launch pulling process
// Return number of downloaded artifacts
func (m Maven) Pull() (int, error) {
	counter := 0
	replacer := strings.NewReplacer(".", "/")

	err := utils.InitializeDatabase(m.DatabaseFile, "CREATE TABLE IF NOT EXISTS artifact (id INTEGER PRIMARY KEY, `name` TEXT, version TEXT)")
	if err != nil {
		return counter, err
	}

	for _, artifact := range m.Artifacts {
		group := replacer.Replace(artifact.Group)
		artifactId := replacer.Replace(artifact.Id)
		metadata, err := m.readMetadata(group, artifactId)
		if err != nil {
			return counter, err
		}

		if len(metadata.Versioning.Versions) != 0 {
			for _, version := range metadata.Versioning.Versions[0].Versions {
				isExistInDB, err := utils.ExistsInDatabase(m.DatabaseFile, "SELECT id FROM artifact WHERE `name` = ? AND version = ?", fmt.Sprintf("%s.%s", group, artifact), version)
				if err != nil {
					return counter, err
				}
				if strings.Compare(version, artifact.MinimumVersion) >= 0 && !isExistInDB {
					err = m.downloadWithDependencies(artifact.Group, artifact.Id, version)
					counter++
				}
			}
		}
	}

	return counter, nil
}

// private function to create the POM file
func (m Maven) createPOM(group string, artifact string, version string) error {
	project := project{
		ModelVersion: "4.0.0",
		GroupId:      "automirror",
		ArtifactId:   "automirror",
		Version:      "0.0.0",
		Dependencies: []dependencies{
			{
				Dependencies: []dependency{
					{
						GroupId:    group,
						ArtifactId: artifact,
						Version:    version,
					},
				},
			},
		},
	}

	file, err := os.Create(m.POMFile)
	if err != nil {
		return err
	}
	xmlWriter := io.Writer(file)

	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("  ", "    ")
	return enc.Encode(project)
}

// Private method to get archive list of artifact to download
func (m Maven) downloadWithDependencies(group string, artifact string, version string) error {
	err := m.createPOM(group, artifact, version)
	if err != nil {
		return err
	}

	cmd := exec.Command(
		"mvn",
		"clean",
		"compile",
		"dependency:sources",
		"dependency:resolve",
		"-f",
		m.POMFile,
		"-DdownloadSources=true",
		"-DdownloadJavadocs=true",
		fmt.Sprintf("-Dmaven.repo.local=%s", m.Folder),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return utils.InsertIntoDatabase(m.DatabaseFile, "INSERT INTO artifact (`name`, version) VALUES (?, ?)", fmt.Sprintf("%s.%s", group, artifact), version)
}

// Private method to read Maven Metadata File from Repo
// One file per artifacts
// Return the Metadata structure
func (m Maven) readMetadata(group string, artifact string) (metadata, error) {
	var metadata metadata

	resp, err := http.Get(strings.Join([]string{m.Url, group, artifact, m.MetadataFileName}, "/"))
	if err != nil {
		return metadata, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return metadata, err
	}

	err = xml.Unmarshal(body, &metadata)
	return metadata, err
}

// Metadata Structure
type metadata struct {
	XMLName    xml.Name   `xml:"metadata"`
	GroupId    string     `xml:"groupId"`
	ArtifactId string     `xml:"artifactId"`
	Versioning versioning `xml:"versioning"`
}

type versioning struct {
	XMLName     xml.Name   `xml:"versioning"`
	Latest      string     `xml:"latest"`
	Release     string     `xml:"release"`
	Versions    []versions `xml:"versions"`
	LastUpdated string     `xml:"lastUpdated"`
}

type versions struct {
	XMLName  xml.Name `xml:"versions"`
	Versions []string `xml:"version"`
}

// POM structure
type project struct {
	XMLName      xml.Name       `xml:"project"`
	ModelVersion string         `xml:"modelVersion"`
	GroupId      string         `xml:"groupId"`
	ArtifactId   string         `xml:"artifactId"`
	Version      string         `xml:"version"`
	Dependencies []dependencies `xml:"dependencies"`
}

type dependencies struct {
	XMLName      xml.Name     `xml:"dependencies"`
	Dependencies []dependency `xml:"dependency"`
}

type dependency struct {
	XMLName    xml.Name `xml:"dependency"`
	GroupId    string   `xml:"groupId"`
	ArtifactId string   `xml:"artifactId"`
	Version    string   `xml:"version"`
}
