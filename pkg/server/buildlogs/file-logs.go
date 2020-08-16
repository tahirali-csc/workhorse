package buildlogs

import (
	"os"
	"path"
	"workhorse/pkg/server/api"

	"github.com/google/uuid"
)

type fileContainerLogs struct {
	file *os.File
	Path string
}

func (fl *fileContainerLogs) Write(data []byte) {
	fl.file.WriteString(string(data))
}

func (fl *fileContainerLogs) GetLocation() string {
	return fl.file.Name()
}

func (fl *fileContainerLogs) createFile(baseDir string) (*os.File, error) {
	folderName := uuid.New()
	logPath := path.Join(baseDir, folderName.String())
	os.MkdirAll(logPath, 0755)
	return os.Create(path.Join(logPath, "logs.txt"))
}

func NewContainerLogsWriter(serverConfig api.ServerConfig) ContainerLogsWriter {
	fileLogs := &fileContainerLogs{}
	fileLogs.file, _ = fileLogs.createFile(serverConfig.ContainerLogsFolder)
	return fileLogs
}

// func NewFileLogs(baseDir string) *fileLogs {
// 	file, _ := createFile(baseDir)

// 	return &FileLogs{
// 		file: file,
// 		Path: file.Name(),
// 	}
// }

// func (fl *FileLogs) Write(data []byte) {
// 	fl.file.WriteString(string(data))
// }

// func createFile(baseDir string) (*os.File, error) {
// 	folderName := uuid.New()
// 	logPath := path.Join(baseDir, folderName.String())
// 	log.Print("Full Path::::", logPath)
// 	os.MkdirAll(logPath, 0755)
// 	return os.Create(path.Join(logPath, "logs.txt"))
// }
