package factorcfg

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func setupEnv() {
	os.Setenv("FACTORCFG_A", "hey===stuff=things=")
	os.Setenv("FACTORCFG_B", "5")
	os.Setenv("FACTORCFG_C", "1.6")
	os.Setenv("FACTORCFG_D", "hey,there,guy")
	os.Setenv("FACTORCFG_E", "true")
	os.Setenv("FACTORCFG_NESTED_A", "there")
	os.Setenv("FACTORCFG_NESTED_B", "4")
}

func TestEnv(t *testing.T) {
	setupEnv()
	Convey("Given a loader using Env", t, func() {
		loader := NewLoader()
		loader.Use(NewEnv())

		Convey("We should find the environment in our spec", func() {
			spec := &cfgTest{}
			err := loader.Load(spec)
			So(err, ShouldBeNil)
			So(spec, ShouldResemble, cfgExample)
		})
	})
}

// var (
// 	envTplTest = `{{ range $_, $v := . }}{{index $v.Tags "env"}}={{$v.String}} ({{$v.Type}})
//   {{index $v.Tags "doc"}}
// {{ end }}`
// 	envTplResult = []byte(`FACTORCFG_A="hey" (string)
//   doc of A
// FACTORCFG_B=5 (int)

// FACTORCFG_C=1.6 (float32)

// FACTORCFG_D=[]string{"hey", "there", "guy"} ([]string)

// FACTORCFG_E=true (bool)

// FACTORCFG_NESTED_A="there" (string)

// FACTORCFG_NESTED_B=4 (int)
//   doc of Nested.B
// `)
// )

// func TestEnvRender(t *testing.T) {
// 	setupEnv()

// 	Convey("Given a loader using Env", t, func() {
// 		loader := NewLoader()
// 		loader.Use(NewEnv())

// 		Convey("We should render our spec correctly", func() {
// 			tmpl, err := template.New("spec").Parse(envTplTest)
// 			So(tmpl, ShouldNotBeNil)
// 			So(err, ShouldBeNil)
// 			data, err := Render(cfgExample, tmpl)
// 			So(err, ShouldBeNil)
// 			So(string(data), ShouldEqual, string(envTplResult))
// 		})
// 	})
// }
