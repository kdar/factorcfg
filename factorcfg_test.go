package factorcfg

type cfgTestNested struct {
	A string `env:"FACTORCFG_NESTED_A"`
	B int    `doc:"doc of Nested.B" env:"FACTORCFG_NESTED_B"`
}

type cfgTest struct {
	A          string   `doc:"doc of A" env:"FACTORCFG_A"`
	B          int      `env:"FACTORCFG_B"`
	C          float32  `env:"FACTORCFG_C"`
	D          []string `env:"FACTORCFG_D"`
	E          bool     `env:"FACTORCFG_E"`
	unexported string

	Nested cfgTestNested
}

var cfgExample = &cfgTest{
	A: "hey===stuff=things=",
	B: 5,
	C: 1.6,
	D: []string{"hey", "there", "guy"},
	E: true,
	Nested: cfgTestNested{
		A: "there",
		B: 4,
	},
}
