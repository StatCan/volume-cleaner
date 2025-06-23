package structure

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
