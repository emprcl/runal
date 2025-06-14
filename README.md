# Runal

![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/emprcl/runal) ![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/emprcl/runal/build.yml)

Runal is a simple creative coding tool for the terminal.
It works similarly as [p5js](https://p5js.org/) and can either be programmed with JavaScript, or used as a Go package.

> :warning: _Runal is a work-in-progress. It has only been tested on Linux and the API should not be considered as stable._

_Feel free to [open an issue](https://github.com/emprcl/runal/issues/new)._

## Installation

### Download executables

[Download the last release](https://github.com/emprcl/runal/releases) for your platform.

Then:
```sh
# Extract files
mkdir -p runal && tar -zxvf runal_VERSION_PLATFORM.tar.gz -C runal
cd runal

# Run runal
./runal -f sketch.js
```

### Go install

If your a developer using Go, you can use the `go install` command:
```
go install github.com/emprcl/runal@latest
```

## Usage

### JavaScript runtime

You can use JavaScript for scripting your sketch. Your js file should contain a `setup` and a `draw` method, like so:
```js
// sketch.js

function setup() {}

function draw() {}
```

And you can then execute the file with:
```sh
./runal -f sketch.js
```

The js file will be automatically reloaded when modified, no need to restart the command.


### Go package

Because Runal is written in Go, you can also use it as a Go package.

```go
// sketch.go
package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/emprcl/runal"
)

func main() {
	runal.Run(context.Background(), setup, draw, onKey)
}

func setup(c *runal.Canvas) {}

func draw(c *runal.Canvas) {}

func onKey(c *runal.Canvas, key string) {}
```

Then, simply build it:
```
go run sketch.go
```
