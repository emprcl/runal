//go:build js

package canvas

// LoadImage is unsupported on the web build: decoding image files would pull
// the jpeg/png/webp decoders and the mosaic renderer into the wasm binary, and
// there is no local filesystem to read from. Returns nil.
func (c *Canvas) LoadImage(path string) Image {
	logErrorf("LoadImage is not supported on the web build")
	return nil
}
