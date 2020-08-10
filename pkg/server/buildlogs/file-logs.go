package buildlogs

import (
	"os"
	"path"

	"github.com/google/uuid"
)

type FileLogs struct {
	file *os.File
	Path string
}

func NewFileLogs(baseDir string) *FileLogs {
	file, _ := createFile(baseDir)

	return &FileLogs{
		file: file,
		Path: file.Name(),
	}
}

func (fl *FileLogs) Write(data []byte) {
	fl.file.WriteString(string(data))
}

func createFile(baseDir string) (*os.File, error) {
	folderName := uuid.New()
	logPath := path.Join(baseDir, folderName.String())
	os.MkdirAll(logPath, 0755)
	return os.Create(path.Join(logPath, "logs.txt"))
}
