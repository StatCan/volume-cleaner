package structure

type ControllerConfig struct {
	Namespace  string
	Label      string
	TimeFormat string
}

type SchedulerConfig struct {
	Namespace   string
	Label       string
	TimeFormat  string
	GracePeriod int
	DryRun      bool
}

type EmailConfig struct {
	BaseURL         string
	Endpoint        string
	EmailTemplateID string
	APIKey          string
}
