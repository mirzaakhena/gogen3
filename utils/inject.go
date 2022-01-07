package utils

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

const injectedCodeLocation = "//!"

// InjectToInteractor ...
func InjectToInteractor(interactorFilename, injectedCode string) ([]byte, error) {

	existingFile := interactorFilename

	// open interactor file
	file, err := os.Open(existingFile)
	if err != nil {
		return nil, err
	}

	needToInject := false

	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	for scanner.Scan() {
		row := scanner.Text()

		// check the injected code in interactor
		if strings.TrimSpace(row) == injectedCodeLocation {

			needToInject = true

			// inject code
			buffer.WriteString(injectedCode)
			buffer.WriteString("\n")

			continue
		}

		buffer.WriteString(row)
		buffer.WriteString("\n")
	}

	// if no injected marker found, then abort the next step
	if !needToInject {
		return nil, nil
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	// rewrite the file
	if err := ioutil.WriteFile(existingFile, buffer.Bytes(), 0644); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func InjectCodeAtTheEndOfFile(filename, templateCode string) ([]byte, error) {

	// reopen the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	for scanner.Scan() {
		row := scanner.Text()

		buffer.WriteString(row)
		buffer.WriteString("\n")
	}

	// write the template in the end of file
	buffer.WriteString(templateCode)
	buffer.WriteString("\n")

	return buffer.Bytes(), nil

}
