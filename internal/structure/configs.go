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
NOTIF_TIMES: "1, 2, 3, 4, 7, 30"
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
	NotifTimes  []int
}

type EmailConfig struct {
	BaseURL         string
	Endpoint        string
	EmailTemplateID string
	APIKey          string
}
