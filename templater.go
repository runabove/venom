package venom

import (
	"gopkg.in/yaml.v2"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// Templater contains templating values on a testsuite
type Templater struct {
	Values map[string]string
}

func newTemplater(values map[string]string) *Templater {
	if values == nil {
		values = make(map[string]string)
	}
	return &Templater{Values: values}
}

// Add add data to templater
func (tmpl *Templater) Add(prefix string, values map[string]string) {
	if tmpl.Values == nil {
		tmpl.Values = make(map[string]string)
	}
	for k, v := range values {
		tmpl.Values[prefix+"."+k] = v
	}
}

// Apply apply vars on string
func (tmpl *Templater) Apply(step TestStep) (TestStep, error) {

	log.Debugf("templater> before: %#v", step)

	// Using yaml to encode/decode, it generates map[interface{}]interface{} typed data that json does not like
	s, err := yaml.Marshal(step)
	if err != nil {
		return nil, fmt.Errorf("templater> Error while marshaling: %s", err)
	}
	sb := string(s)

	for k, v := range tmpl.Values {
		sb = strings.Replace(sb, "{{."+k+"}}", v, -1)
	}

	var t TestStep
	if err := yaml.Unmarshal([]byte(sb), &t); err != nil {
		return nil, fmt.Errorf("templater> Error while unmarshal: %s", err)
	}

	log.Debugf("templater> after: %+v", t)

	return t, nil
}
