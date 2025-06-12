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
	GracePeriod string
}
