package akizuki

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/porfirion/trie"
)

type defaultDetector struct {
	path    string
	visited *trie.Trie[struct{}]
}

const defaultDatabasePath = "akizuki.json"

func NewDefaultDetector(databasePath string) (*defaultDetector, error) {

	visited := []string{}
	if _, err := os.Stat(databasePath); err == nil {
		f, err := os.Open(databasePath)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, &visited)
		if err != nil {
			return nil, err
		}
	}

	ret := defaultDetector{
		path:    databasePath,
		visited: trie.BuildPrefixesOnly(visited...),
	}
	return &ret, nil
}

func (d defaultDetector) NewPages(urls []string) ([]string, error) {
	ret := []string{}
	for _, url := range urls {
		_, visited := d.visited.GetByString(url)
		if !visited {
			ret = append(ret, url)
		}
	}
	return ret, nil
}
func (d defaultDetector) AddPages(urls []string) error {
	for _, url := range urls {
		d.visited.PutString(url, struct{}{})
	}

	visited := []string{}

	d.visited.Iterate(func(prefix []byte, value struct{}) {
		visited = append(visited, string(prefix))
	})

	j, err := json.Marshal(visited)
	if err != nil {
		return err
	}

	f, err := os.Create(d.path)
	if err != nil {
		return err
	}
	defer f.Close()
	n, err := f.Write(j)
	if err != nil {
		return err
	}
	if n != len(j) {
		return errors.New("failed to write")
	}

	return nil
}
