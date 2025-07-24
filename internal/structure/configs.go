package structure

/*
Example configs

controller:

NAMESPACE: "anray-liu"
TIME_LABEL: "volume-cleaner/unattached-time"
NOTIF_LABEL: "volume-cleaner/notification-count"
TIME_FORMAT: "2006-01-02_15-04-05Z"

scheduler:

NAMESPACE: "anray-liu"
TIME_LABEL: "volume-cleaner/unattached-time"
NOTIF_LABEL: "volume-cleaner/notification-count"
GRACE_PERIOD: "180"
TIME_FORMAT: "2006-01-02_15-04-05Z"
DRY_RUN: "true"
NOTIF_TIMES: "1, 2, 3, 4, 7, 30"

BASE_URL: "https://api.notification.canada.ca",
ENDPOINT: "/v2/notifications/email",
EMAIL_TEMPLATE_ID: "Random Template",
API_KEY: "Random APIKEY",
*/

type ControllerConfig struct {
	Namespace  string
	TimeLabel  string
	NotifLabel string
	TimeFormat string
}

type SchedulerConfig struct {
	Namespace   string
	TimeLabel   string
	NotifLabel  string
	TimeFormat  string
	GracePeriod int
	DryRun      bool
	NotifTimes  []int
	EmailCfg    EmailConfig
}

type EmailConfig struct {
	BaseURL         string
	Endpoint        string
	EmailTemplateID string
	APIKey          string
}

// For internal use

// Represents the main request body structure for sending Email Notifications with GC Notify
type RequestBody struct {
	EmailAddress    string          `json:"email_address"`
	TemplateID      string          `json:"template_id"`
	Personalisation Personalisation `json:"personalisation"`
}

// Represents the variables used in the email template when calling GC Notify
type Personalisation struct {
	Name         string `json:"name"`
	VolumeName   string `json:"volume_name"`
	GracePeriod  string `json:"grace_period"`
	DeletionDate string `json:"deletion_date"`
}
