package manifest

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type Manifest struct {
	Balancers Balancers
	Queues    Queues
	Services  Services
	Tables    Tables

	root string
}

func Load(data []byte) (*Manifest, error) {
	var m Manifest

	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func LoadFile(path string) (*Manifest, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	m, err := Load(data)
	if err != nil {
		return nil, err
	}

	root, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return nil, err
	}

	m.root = root

	return m, nil
}

func (m *Manifest) Path(sub string) (string, error) {
	if m.root == "" {
		return "", fmt.Errorf("path undefined for a manifest with no root")
	}

	return filepath.Join(m.root, sub), nil
}

func message(w io.Writer, format string, args ...interface{}) {
	if w != nil {
		w.Write([]byte(fmt.Sprintf(format, args...) + "\n"))
	}
}
