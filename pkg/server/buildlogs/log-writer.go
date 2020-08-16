package buildlogs

type ContainerLogsWriter interface {
	Write(data []byte)
	GetLocation() string
}
