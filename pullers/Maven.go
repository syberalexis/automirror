package pullers

import (
	"automirror/configs"
	"database/sql"
	"encoding/xml"
	"fmt"
	"github.com/BurntSushi/toml"
	"io"

	"io/ioutil"
	"log"
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

func BuildMaven(pullerConfig configs.PullerConfig) Puller {
	var config Maven
	tomlFile, err := ioutil.ReadFile(pullerConfig.Config)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := toml.Decode(string(tomlFile), &config); err != nil {
		log.Fatal(err)
	}

	config.Url = pullerConfig.Source
	config.Folder = pullerConfig.Destination
	return config
}

// Inherits public method to launch pulling process
// Return number of downloaded artifacts
func (m Maven) Pull() int {
	replacer := strings.NewReplacer(".", "/")
	m.initDatabase()
	counter := 0

	for _, artifact := range m.Artifacts {
		group := replacer.Replace(artifact.Group)
		artifactId := replacer.Replace(artifact.Id)
		metadata := m.readMetadata(group, artifactId)

		if len(metadata.Versioning.Versions) != 0 {
			for _, version := range metadata.Versioning.Versions[0].Versions {
				if strings.Compare(version, artifact.MinimumVersion) >= 0 && !m.existsInDatabase(artifact.Group, artifact.Id, version) {
					m.downloadWithDependencies(artifact.Group, artifact.Id, version)
					counter++
				}
			}
		}
	}

	return counter
}

func (m Maven) createPOM(group string, artifact string, version string) {
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

	file, _ := os.Create(m.POMFile)
	xmlWriter := io.Writer(file)

	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("  ", "    ")
	if err := enc.Encode(project); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

// Private method to get archive list of artifact to download
func (m Maven) downloadWithDependencies(group string, artifact string, version string) {
	m.createPOM(group, artifact, version)

	cmd := exec.Command("mvn", "clean", "compile", "dependency:sources", "dependency:resolve", "-f", m.POMFile, "-DdownloadSources=true", "-DdownloadJavadocs=true", "-Dmaven.repo.local="+m.Folder)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	m.insertIntoDatabase(group, artifact, version)
}

// Private method to check if version of artifact is already downloaded
// Return a boolean
func (m Maven) existsInDatabase(group string, artifact string, version string) bool {
	database, _ := sql.Open("sqlite3", m.DatabaseFile)
	statement, _ := database.Prepare("SELECT id FROM artifact WHERE `name` = ? AND version = ?")
	rows, _ := statement.Query(group+"."+artifact, version)

	return rows.Next()
}

// Private method to insert downloaded artifact info into Database
func (m Maven) insertIntoDatabase(group string, artifact string, version string) {
	database, _ := sql.Open("sqlite3", m.DatabaseFile)
	statement, _ := database.Prepare("INSERT INTO artifact (`name`, version) VALUES (?, ?)")

	_, err := statement.Exec(group+"."+artifact, version)
	if err != nil {
		log.Fatal(err)
	}
}

// Private method to initialize SQLite Database
func (m Maven) initDatabase() {
	database, _ := sql.Open("sqlite3", m.DatabaseFile)
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS artifact (id INTEGER PRIMARY KEY, `name` TEXT, version TEXT)")

	_, err := statement.Exec()
	if err != nil {
		log.Fatal(err)
	}
}

// Private method to read Maven Metadata File from Repo
// One file per artifacts
// Return the Metadata structure
func (m Maven) readMetadata(group string, artifact string) metadata {
	var metadata metadata

	resp, err := http.Get(strings.Join([]string{m.Url, group, artifact, m.MetadataFileName}, "/"))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = xml.Unmarshal(body, &metadata)
	if err != nil {
		log.Fatal(err)
	}

	return metadata
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
