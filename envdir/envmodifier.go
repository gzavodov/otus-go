package envdir

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// EnvModifier is simple utility allows run program with custom environment variables.
type EnvModifier struct {
}

// Run executes program defined by executablePath with environment variables defined in dirPath
// Returns program output or error
func (m *EnvModifier) Run(dirPath string, executablePath string) ([]byte, error) {

	dirInfo, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return nil, err
	}

	if !dirInfo.IsDir() {
		return nil, fmt.Errorf("the path '%s' is not a directory", dirPath)
	}

	fileInfos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	variables := make(map[string]string)

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, fileInfo.Name())

		s, err := m.readFirstLine(filePath)
		if err != nil {
			return nil, err
		}

		variables[fileInfo.Name()] = s
	}

	for key, value := range variables {
		//fmt.Printf("os.Setenv('%s', '%s')\n", key, value)
		os.Setenv(key, value)
	}

	//fmt.Printf("exec.Command('%s')\n", executablePath)
	cmd := exec.Command(executablePath)

	output, err := cmd.Output()

	for key := range variables {
		os.Unsetenv(key)
	}

	return output, err
}

func (m *EnvModifier) readFirstLine(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	s, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}

	return s, nil
}
