package pullers

import (
	"automirror/configs"
	"automirror/utils"
	"database/sql"
	"fmt"
	"github.com/BurntSushi/toml"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Python struct {
	Url          string
	Folder       string
	DatabaseFile string `toml:"database_file"`
}

func BuildPython(pullerConfig configs.PullerConfig) Puller {
	var config Python
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

func (p Python) Pull() int {
	p.initDatabase()
	return p.readRepository("/simple/")
}

// Private method to get archive list of artifact to download
func (p Python) readRepository(subpath string) int {
	counter := 0
	resp, err := http.Get(p.Url + subpath)
	if err != nil {
		log.Fatal(err)
	}

	z := html.NewTokenizer(resp.Body)
	finished := false
	for !finished {
		next := z.Next()
		switch {
		case next == html.ErrorToken:
			// End of the document, we're done
			finished = true
			break
		case next == html.StartTagToken:
			token := z.Token()
			if token.Data == "a" && len(token.Attr) > 0 && token.Attr[0].Val != "../" {
				if strings.HasSuffix(token.Attr[0].Val, "/") {
					counter += p.readRepository(token.Attr[0].Val)
				} else {
					match := p.match(token.Attr[0].Val)
					if match != "" && !p.existsInDatabase(match) {
						p.download(subpath, token.Attr[0].Val)
						counter++
					}
				}
			}
		}
	}
	return counter
}

func (p Python) match(url string) string {
	re := regexp.MustCompile("^.*/(.+\\.(tar.gz|whl|zip|bz2|tar.bz2))#?.*$")
	match := re.FindStringSubmatch(url)
	if match != nil {
		return match[1]
	}
	return ""
}

// Private method to download artifacts
func (p Python) download(subpath string, url string) {
	match := p.match(url)
	if match != "" {
		file := strings.Join([]string{p.Folder, strings.Replace(subpath, "/simple", "", 1), match}, "")

		if err := utils.FileDownloader(url, file); err != nil {
			panic(err)
		}
		fmt.Printf("%s successfully pulled !\n", file)

		p.insertIntoDatabase(match)
	} else {
		log.Print(url + " not matched")
	}
}

// Private method to check if version of artifact is already downloaded
// Return a boolean
func (p Python) existsInDatabase(archive string) bool {
	database, _ := sql.Open("sqlite3", p.DatabaseFile)
	statement, _ := database.Prepare("SELECT id FROM package WHERE `name` = ?")
	rows, _ := statement.Query(archive)

	return rows.Next()
}

// Private method to insert downloaded artifact info into Database
func (p Python) insertIntoDatabase(archive string) {
	database, _ := sql.Open("sqlite3", p.DatabaseFile)
	statement, _ := database.Prepare("INSERT INTO package (`name`) VALUES (?)")

	_, err := statement.Exec(archive)
	if err != nil {
		log.Fatal(err)
	}
}

// Private method to initialize SQLite Database
func (p Python) initDatabase() {
	database, _ := sql.Open("sqlite3", p.DatabaseFile)
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS package (id INTEGER PRIMARY KEY, `name` TEXT)")

	_, err := statement.Exec()
	if err != nil {
		log.Fatal(err)
	}
}
