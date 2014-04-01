package factorcfg

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type optional struct {
	sources []interface{}
}

func Optional(s ...interface{}) optional {
	return optional{sources: s}
}

type EnvFile struct {
	sources []interface{}
}

func NewEnvFile(sources ...interface{}) *EnvFile {
	e := &EnvFile{sources: sources}
	return e
}

func (e *EnvFile) All() (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := e.allSource(m, e.sources)
	return m, err
}

func (e *EnvFile) allSource(m map[string]interface{}, sources []interface{}) error {
	for _, source := range sources {
		switch s := source.(type) {
		case optional:
			err := e.allSource(m, s.sources)
			if err != nil {
				if _, ok := err.(*os.PathError); !ok {
					return err
				}
			}
		case string:
			fp, err := os.Open(s)
			if err != nil {
				return err
			}
			defer fp.Close()
			err = e.intoMap(m, fp)
			if err != nil {
				return err
			}
		case io.Reader:
			err := e.intoMap(m, s)
			if err != nil {
				return err
			}
		default:
			return errors.New("unknown source")
		}
	}

	return nil
}

func (e *EnvFile) Tag() string {
	return "env"
}

func (e *EnvFile) intoMap(m map[string]interface{}, reader io.Reader) (err error) {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")

	for _, fullLine := range lines {
		if !isIgnoredLine(fullLine) {
			key, value, err := parseLine(fullLine)

			if err == nil {
				m[key] = value
			}
		}
	}
	return
}

func parseLine(line string) (key string, value string, err error) {
	if len(line) == 0 {
		err = errors.New("zero length string")
		return
	}

	// ditch the comments (but keep quoted hashes)
	if strings.Contains(line, "#") {
		segmentsBetweenHashes := strings.Split(line, "#")
		quotesAreOpen := false
		segmentsToKeep := make([]string, 0)
		for _, segment := range segmentsBetweenHashes {
			if strings.Count(segment, "\"") == 1 || strings.Count(segment, "'") == 1 {
				if quotesAreOpen {
					quotesAreOpen = false
					segmentsToKeep = append(segmentsToKeep, segment)
				} else {
					quotesAreOpen = true
				}
			}

			if len(segmentsToKeep) == 0 || quotesAreOpen {
				segmentsToKeep = append(segmentsToKeep, segment)
			}
		}

		line = strings.Join(segmentsToKeep, "#")
	}

	// now split key from value
	splitString := strings.SplitN(line, "=", 2)

	if len(splitString) != 2 {
		// try yaml mode!
		splitString = strings.SplitN(line, ":", 2)
	}

	if len(splitString) != 2 {
		err = errors.New("Can't separate key from value")
		return
	}

	// Parse the key
	key = splitString[0]
	if strings.HasPrefix(key, "export") {
		key = strings.TrimPrefix(key, "export")
	}
	key = strings.Trim(key, " ")

	// Parse the value
	value = splitString[1]
	// trim
	value = strings.Trim(value, " ")

	// check if we've got quoted values
	if strings.Count(value, "\"") == 2 || strings.Count(value, "'") == 2 {
		// pull the quotes off the edges
		value = strings.Trim(value, "\"'")

		// expand quotes
		value = strings.Replace(value, "\\\"", "\"", -1)
		// expand newlines
		value = strings.Replace(value, "\\n", "\n", -1)
	}

	return
}

func isIgnoredLine(line string) bool {
	trimmedLine := strings.Trim(line, " \n\t")
	return len(trimmedLine) == 0 || strings.HasPrefix(trimmedLine, "#")
}
