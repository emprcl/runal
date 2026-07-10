module github.com/emprcl/runal/cmd/runal-wasm

go 1.26.4

replace github.com/emprcl/runal => ../../

replace github.com/emprcl/runal/x/js => ../../x/js

require (
	github.com/emprcl/runal v0.0.0-00010101000000-000000000000
	github.com/emprcl/runal/x/js v0.0.0-00010101000000-000000000000
)
