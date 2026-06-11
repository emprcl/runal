package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
)

// base64 encoded delimiter that should appear at the top
// of a javascript sketch appended to the runal binary
const embedDelimiter = "Ly9ydW5hbDplbWJlZAo="

func readEmbeddedSketch() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(executable)
	if err != nil {
		return "", err
	}

	delimiter, err := base64.StdEncoding.DecodeString(embedDelimiter)
	if err != nil {
		return "", err
	}

	parts := strings.Split(string(content), string(delimiter))
	if len(parts) < 2 {
		return "", nil
	}

	return parts[len(parts)-1], nil
}

// createEmbeddedExecutable creates a copy of the current executable
// (named 'outFile') with the contents of 'appendFile' appended to it
func createEmbeddedExecutable(outFilename string, appendFilename string) error {
	// read the contents of the appendFile
	sketch, err := os.ReadFile(appendFilename)
	if err != nil {
		return fmt.Errorf("%w reading %s", err, appendFilename)
	}

	// create an io.Reader and an io.Writer for io.Copy
	execName, err := os.Executable()
	if err != nil {
		return err
	}
	execFile, err := os.Open(execName)
	if err != nil {
		return fmt.Errorf("%w opening %s", err, execName)
	}
	defer execFile.Close()

	newFile, err := os.Create(outFilename)
	if err != nil {
		return fmt.Errorf("%w creating %s", err, outFilename)
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, execFile)
	if err != nil {
		return fmt.Errorf("%w copying %s to %s", err, execFile.Name(), newFile.Name())
	}

	delimiter, err := base64.StdEncoding.DecodeString(embedDelimiter)
	if err != nil {
		return err
	}

	appStr := fmt.Sprintf("%s\n%s", delimiter, sketch)
	_, err = newFile.WriteString(appStr)
	if err != nil {
		return fmt.Errorf("%w appending sketch to %s", err, newFile.Name())
	}

	// make the permissions on the new executable the same as the old one
	srcInfo, err := os.Stat(execName)
	if err != nil {
		return fmt.Errorf("%w getting permissions for %s", err, execName)
	}
	err = os.Chmod(outFilename, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("%w setting permissions for %s", err, outFilename)
	}

	return nil
}
