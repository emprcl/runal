# Runal

![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/emprcl/runal) ![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/emprcl/runal/build.yml)

:notebook: **[User Manual](https://empr.cl/runal/)**

Runal is a simple creative coding environment for text and ascii art, that runs in the terminal. It works similarly as [processing](https://processing.org/) or [p5js](https://p5js.org/) and can either be programmed with [JavaScript](https://developer.mozilla.org/en-US/docs/Web/JavaScript), or used as a [Go](https://go.dev/) package.

**_Runal is a work-in-progress. It has only been tested on Linux and the API should not be considered as stable until it reaches 1.0._**

_Feel free to [open an issue](https://github.com/emprcl/runal/issues/new)._

![signls screenshot](/docs/screenshot.png)

## Installation

### Quick-install

On **linux** or **macOS**, you can run this quick-install bash script:
```sh
curl -sSL empr.cl/get/runal | bash
```

### Manual installation

#### Linux & macOS

[Download the last release](https://github.com/emprcl/runal/releases) for your platform.

Then:
```sh
# Extract files
mkdir -p runal && tar -zxvf runal_VERSION_PLATFORM.tar.gz -C runal
cd runal

# Run runal
./runal

# Run runal demo
./runal -demo
```

#### Windows

> _We recommend using Windows Terminal with a good monospace font like Iosevka to display Signls correctly on Windows._

Unzip the last [windows release](https://github.com/emprcl/runal/releases) and, in the same directory, run:
```winbatch
; Run runal
.\runal.exe

; Run runal demo
.\runal.exe -demo
```

### Go install

If your a developer using Go, you can use the `go install` command:
```
go install github.com/emprcl/runal@latest
```

### Build it yourself

You'll need [go 1.23](https://go.dev/dl/) minimum.
Although you should be able to build it for either **linux**, **macOS** or **Windows**, it has only been tested on **linux**.

```sh
# Linux
make GOLANG_OS=linux build

# macOS
make GOLANG_OS=darwin build

# Windows
make GOLANG_OS=windows build

# Raspberry Pi OS
make GOLANG_OS=linux GOLANG_ARCH=arm64 build
```


## Usage

### JavaScript runtime

You can use JavaScript for scripting your sketch. Your js file should contain a `setup` and a `draw` method. Both methods take a single argument (here `c`) representing a canvas object that holds all the available primitives:
```js
// sketch.js

function setup(c) {}

function draw(c) {}
```

You can add an extra method `onKey` to catch keyboard events:
```js
function onKey(c, key) {}
````

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
	runal.Run(context.Background(), setup, draw, nil)
}

func setup(c *runal.Canvas) {}

func draw(c *runal.Canvas) {}
```

Then, simply build it:
```
go run sketch.go
```

## Documentation

Check the [API reference](https://empr.cl/runal/#reference).
You can also check some examples in the [examples directory](https://github.com/emprcl/runal/tree/main/examples).

## Acknowledgments

Runal uses a few awesome packages:
 - [dop251/goja](https://github.com/dop251/goja) for the JavaScript engine
 - [fsnotify/fsnotify](https://github.com/fsnotify/fsnotify) for watching file changes in realtime
 - [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) for handling colors
