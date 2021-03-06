package pullers

import (
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/configs"
	"github.com/syberalexis/automirror/utils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// Maven object to pull Java artifacts from a Maven repo with mvn unix config
type Maven struct {
	Source           string
	Destination      string
	MetadataFileName string     `toml:"metadata_file_name"`
	POMFile          string     `toml:"pom_file"`
	DatabaseFile     string     `toml:"database_file"`
	Artifacts        []Artifact `toml:"artifact"`
}

// Artifact to pull
type Artifact struct {
	Group          string
	ID             string
	MinimumVersion string `toml:"minimum_version"`
}

// NewMaven method to construct Maven
func NewMaven(config configs.EngineConfig) (interface{}, error) {
	var maven Maven
	err := configs.Parse(&maven, config.Config)
	if err != nil {
		return nil, err
	}
	return maven, nil
}

// Pull artifacts from Maven repo
// Inherits public method to launch pulling process
// Return number of downloaded artifacts
func (m Maven) Pull(log *log.Logger) (int, error) {
	counter := 0
	replacer := strings.NewReplacer(".", "/")

	for _, artifact := range m.Artifacts {
		group := replacer.Replace(artifact.Group)
		artifactID := replacer.Replace(artifact.ID)
		metadata, err := m.readMetadata(group, artifactID)
		if err != nil {
			return counter, err
		}

		if len(metadata.Versioning.Versions) != 0 {
			for _, version := range metadata.Versioning.Versions[0].Versions {
				isExistInDB, err := utils.ExistsInDatabase(m.DatabaseFile, fmt.Sprintf("%s.%s:%s", artifact.Group, artifact.ID, version))
				if err != nil {
					return counter, err
				}
				if strings.Compare(version, artifact.MinimumVersion) >= 0 && !isExistInDB {
					err = m.downloadWithDependencies(artifact.Group, artifact.ID, version, log)
					if err != nil {
						return counter, err
					}
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
		GroupID:      "automirror",
		ArtifactID:   "automirror",
		Version:      "0.0.0",
		Dependencies: []dependencies{
			{
				Dependencies: []dependency{
					{
						GroupID:    group,
						ArtifactID: artifact,
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

// Private method to get archive list of artifact to clone
func (m Maven) downloadWithDependencies(group string, artifact string, version string, log *log.Logger) error {
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
		fmt.Sprintf("-Dmaven.repo.local=%s", m.Destination),
	)

	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	err = cmd.Run()
	if err != nil {
		return err
	}

	return utils.InsertIntoDatabase(m.DatabaseFile, fmt.Sprintf("%s.%s:%s", group, artifact, version), "true")
}

// Private method to read Maven Metadata File from Repo
// One file per artifacts
// Return the Metadata structure
func (m Maven) readMetadata(group string, artifact string) (metadata, error) {
	var metadata metadata

	resp, err := http.Get(strings.Join([]string{m.Source, group, artifact, m.MetadataFileName}, "/"))
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
	GroupID    string     `xml:"groupId"`
	ArtifactID string     `xml:"artifactId"`
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
	GroupID      string         `xml:"groupId"`
	ArtifactID   string         `xml:"artifactId"`
	Version      string         `xml:"version"`
	Dependencies []dependencies `xml:"dependencies"`
}

type dependencies struct {
	XMLName      xml.Name     `xml:"dependencies"`
	Dependencies []dependency `xml:"dependency"`
}

type dependency struct {
	XMLName    xml.Name `xml:"dependency"`
	GroupID    string   `xml:"groupId"`
	ArtifactID string   `xml:"artifactId"`
	Version    string   `xml:"version"`
}
