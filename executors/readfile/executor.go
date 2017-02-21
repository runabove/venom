package exec

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fsamin/go-dump"
	"github.com/mitchellh/mapstructure"

	"github.com/runabove/venom"
)

// Name for test readfile
const Name = "readfile"

// New returns a new Test Exec
func New() venom.Executor {
	return &Executor{}
}

// Executor represents a Test Exec
type Executor struct {
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
}

// Result represents a step result
type Result struct {
	Executor    Executor `json:"executor,omitempty" yaml:"executor,omitempty"`
	Content     string   `json:"content,omitempty" yaml:"content,omitempty"`
	Err         string   `json:"error" yaml:"error"`
	TimeSeconds float64  `json:"timeSeconds,omitempty" yaml:"timeSeconds,omitempty"`
	TimeHuman   string   `json:"timeHuman,omitempty" yaml:"timeHuman,omitempty"`
}

// GetDefaultAssertions return default assertions for type exec
func (Executor) GetDefaultAssertions() *venom.StepAssertions {
	return &venom.StepAssertions{Assertions: []string{"result.err ShouldNotExist"}}
}

// Run execute TestStep of type exec
func (Executor) Run(l *log.Entry, aliases venom.Aliases, step venom.TestStep) (venom.ExecutorResult, error) {

	var t Executor
	if err := mapstructure.Decode(step, &t); err != nil {
		return nil, err
	}

	if t.Path == "" {
		return nil, fmt.Errorf("Invalid path")
	}

	start := time.Now()
	result := Result{Executor: t}

	content, errr := t.readfile(t.Path)
	if errr != nil {
		result.Err = errr.Error()
	}
	result.Content = content

	elapsed := time.Since(start)
	result.TimeSeconds = elapsed.Seconds()
	result.TimeHuman = fmt.Sprintf("%s", elapsed)

	return dump.ToMap(result, dump.WithDefaultLowerCaseFormatter())
}

func (e *Executor) readfile(path string) (string, error) {

	fileInfo, _ := os.Stat(path)
	if fileInfo != nil && fileInfo.IsDir() {
		path = filepath.Dir(path) + "/*.yml"
		log.Debugf("path computed:%s", path)
	}

	filesPath, errg := filepath.Glob(path)
	if errg != nil {
		return "", fmt.Errorf("Error reading files on path:%s :%s", path, errg)
	}

	var out string
	for _, f := range filesPath {
		log.Debugf("read %s", f)
		dat, errr := ioutil.ReadFile(f)
		if errr != nil {
			return "", fmt.Errorf("Error while reading file")
		}
		out += string(dat)
	}
	return out, nil
}
