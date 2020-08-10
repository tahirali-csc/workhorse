package buildlogs

type LogWriter interface {
	Write(data []byte)
}
