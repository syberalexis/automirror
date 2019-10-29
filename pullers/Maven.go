package pullers

import (
	"automirror/configs"
	"automirror/utils"
	"database/sql"
	"encoding/xml"
	"github.com/BurntSushi/toml"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Object data structure
type Maven struct {
	Url              string
	Folder           string
	MetadataFileName string     `toml:"metadata_file_name"`
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
				if strings.Compare(version, artifact.MinimumVersion) >= 0 && !m.existsInDatabase(group, artifactId, version) {
					m.getArtifactArchives(group, artifactId, version)
					counter++
				}
			}
		}
	}

	return counter
}

// Private method to get archive list of artifact to download
func (m Maven) getArtifactArchives(group string, artifact string, version string) {
	resp, err := http.Get(strings.Join([]string{m.Url, group, artifact, version}, "/"))
	if err != nil {
		log.Fatal(err)
	}

	z := html.NewTokenizer(resp.Body)
	for {
		next := z.Next()
		switch {
		case next == html.ErrorToken:
			// End of the document, we're done
			return
		case next == html.StartTagToken:
			token := z.Token()
			if token.Data == "a" && len(token.Attr) > 0 && token.Attr[0].Val != "../" {
				m.download(group, artifact, version, token.Attr[0].Val)
			}
		}
	}

	m.insertIntoDatabase(group, artifact, version)
}

// Private method to download artifacts
func (m Maven) download(group string, artifact string, version string, archive string) {
	url := strings.Join([]string{m.Url, group, artifact, version, archive}, "/")
	file := strings.Join([]string{m.Folder, group, artifact, version, archive}, "/")

	if err := utils.FileDownloader(url, file); err != nil {
		panic(err)
	}
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
