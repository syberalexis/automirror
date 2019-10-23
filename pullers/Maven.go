package pullers

import (
	"database/sql"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Maven struct {
	Url              string
	MetadataFileName string
	DatabaseFile     string
	Artifacts        []Artifact
}

type Artifact struct {
	Group          string
	Id             string
	MinimumVersion string
}

func (m Maven) Pull() {
	m.initDatabase()
	for _, artifact := range m.Artifacts {
		metadata := m.readMetadata(artifact.Group, artifact.Id)
		if len(metadata.Versioning.Versions) != 0 {
			for _, version := range metadata.Versioning.Versions[0].Versions {
				if strings.Compare(version, artifact.MinimumVersion) >= 0 && !m.existsInDatabase(artifact.Group, artifact.Id, version) {
					m.download(metadata.GroupId, metadata.ArtifactId, version)
				}
			}
		}
	}
}

func (m Maven) existsInDatabase(group string, artifact string, version string) bool {
	database, _ := sql.Open("sqlite3", m.DatabaseFile)
	statement, _ := database.Prepare("SELECT id FROM artifact WHERE `name` = ? AND version = ?")
	rows, _ := statement.Query(group+"."+artifact, version)

	return rows.Next()
}

func (m Maven) download(group string, artifact string, version string) {
	println(m.Url + "/" + group + "/" + artifact + "/" + version)

	replacer := strings.NewReplacer(".", "/")
	resp, err := http.Get(strings.Join([]string{m.Url, replacer.Replace(group), replacer.Replace(artifact), version}, "/"))
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	println(body)

	m.insertIntoDatabase(group, artifact, version)
}

func (m Maven) insertIntoDatabase(group string, artifact string, version string) {
	database, _ := sql.Open("sqlite3", m.DatabaseFile)
	statement, _ := database.Prepare("INSERT INTO artifact (`name`, version) VALUES (?, ?)")
	_, err := statement.Exec(group+"."+artifact, version)
	if err != nil {
		log.Fatal(err)
	}
}

func (m Maven) initDatabase() {
	database, _ := sql.Open("sqlite3", m.DatabaseFile)
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS artifact (id INTEGER PRIMARY KEY, `name` TEXT, version TEXT)")
	_, err := statement.Exec()
	if err != nil {
		log.Fatal(err)
	}
}

func (m Maven) readMetadata(group string, artifact string) metadata {
	var metadata metadata
	replacer := strings.NewReplacer(".", "/")

	resp, err := http.Get(strings.Join([]string{m.Url, replacer.Replace(group), replacer.Replace(artifact), m.MetadataFileName}, "/"))
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
