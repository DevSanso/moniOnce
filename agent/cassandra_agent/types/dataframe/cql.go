package dataframe

type TracesSession struct {
	SessionID string
	Client    string
	Command   string
	Coordinator string
	Duration int
	Request  int
	Started_at string

	ParametersJson string
}