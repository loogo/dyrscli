package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

type scaffold struct {
	Path  string
	Topic string
	Host  string
}

func (scaffold *scaffold) Scaffold() {
	var fileName string
	if scaffold.Topic == "all" {
		baseURL := fmt.Sprintf("http://%s/subjects/", scaffold.Host)
		resp, err := http.Get(baseURL)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		var subjects []string
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&subjects)

		for _, sub := range subjects {
			if strings.HasSuffix(sub, "-value") {
				fileName = strings.TrimSuffix(sub, "-value") + ".avrc"
				downloadSchema(scaffold, sub, fileName)
			}
		}
	} else {
		fileName = scaffold.Topic + ".avrc"
		downloadSchema(scaffold, scaffold.Topic+"-value", fileName)
	}
}

func downloadSchema(scaffold *scaffold, name string, fileName string) {
	baseURL := fmt.Sprintf("http://%s/subjects/%s/", scaffold.Host, name)
	//get versions
	resp, err := http.Get(baseURL + "versions")

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	var versions []int
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&versions)

	lastVersion := versions[len(versions)-1]

	resp, err = http.Get(baseURL + "versions/" + strconv.Itoa(lastVersion) + "/schema")
	if err != nil {
		log.Fatalln(err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	_ = os.Mkdir(scaffold.Path, 0644)
	fullName := path.Join(scaffold.Path, fileName)
	_ = ioutil.WriteFile(fullName, bodyBytes, 0644)
}
