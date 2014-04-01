package factorcfg

import (
	. "github.com/smartystreets/goconvey/convey"

	"bytes"
	"testing"
)

const (
	testEnvFile = `FACTORCFG_A=hey
FACTORCFG_B=5
FACTORCFG_C=1.6
FACTORCFG_D=hey,there,guy
FACTORCFG_E=true
FACTORCFG_NESTED_A=there
FACTORCFG_NESTED_B=4`
)

func TestEnvFile(t *testing.T) {
	Convey("Given a loader using EnvFile", t, func() {
		loader := NewLoader()
		buf := bytes.NewBufferString(testEnvFile)
		loader.Use(NewEnvFile(buf))

		Convey("We should find the env file in our spec", func() {
			spec := &cfgTest{}
			err := loader.Load(spec)
			So(err, ShouldBeNil)
			So(spec, ShouldResemble, cfgExample)
		})
	})
}
