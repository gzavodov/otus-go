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

	variables := os.Environ()
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, fileInfo.Name())

		s, err := m.readFirstLine(filePath)
		if err != nil {
			return nil, err
		}

		variables = append(variables, fmt.Sprintf("%s=%s", fileInfo.Name(), s))
	}

	cmd := exec.Command(executablePath)
	cmd.Env = variables
	output, err := cmd.Output()

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
