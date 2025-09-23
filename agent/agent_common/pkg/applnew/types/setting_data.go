package types

type SettingData struct {
	ObjectId  int
	LogConfig struct {
		Level string
		Dir  string
		Size int64
	}

	ConfigDb struct {
		Drvier string
		Dsn    string
	}

	ConfigType string
}