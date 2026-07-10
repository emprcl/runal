//go:build js

package canvas

func (c *Canvas) LoadImage(path string) Image {
	logErrorf("LoadImage is not supported on the web build")
	return nil
}
