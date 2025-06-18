package structure

/*
Example configs

controller:

NAMESPACE: "anray-liu"
LABEL: "volume-cleaner/unattached-time"
TIME_FORMAT: "2006-01-02_15-04-05Z"

scheduler:

NAMESPACE: "anray-liu"
LABEL: "volume-cleaner/unattached-time"
GRACE_PERIOD: "180"
TIME_FORMAT: "2006-01-02_15-04-05Z"
DRY_RUN: "true"
*/

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
