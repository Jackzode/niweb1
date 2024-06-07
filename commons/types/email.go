package types

type EmailCodeContent struct {
	SourceType string `json:"source_type"`
	Email      string `json:"e_mail"`
	UserID     string `json:"user_id"`
	// Used for unsubscribe notification
	NotificationSources []string `json:"notification_source,omitempty"`
	// Used for third-party login account binding
	BindingKey string `json:"binding_key,omitempty"`
}
