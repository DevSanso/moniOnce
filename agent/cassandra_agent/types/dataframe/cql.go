package dataframe

type TracesSession struct {
	SessionID string
	Client    string
	Command   string
	Coordinator string
	CoordiantorPort int
	Duration int
	Request  int
	Started_at string

	Parameters map[string]string
}

type SystemViewQueries struct {
	ThreadId string
	QueueMicroSec uint64
	RunningMicroSec uint64
	Text string
}

type SystemViewClients struct {
	Address string
	DriverName string
	ConnectionStage string
	Hostname string
	RequestCnt int64
	Keyspace string
	Username string
}

type SystemViewSystemLog struct {
	Timestamp string
	Level string
	Logger string
	Message string
}