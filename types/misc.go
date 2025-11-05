package types

type Birthdate struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year,omitempty"`
}

type MaskPosition struct {
	Point  string  `json:"point"`
	XShift float32 `json:"x_shift"`
	YShift float32 `json:"y_shift"`
	Scale  float32 `json:"scale"`
}

type Location struct {
	Latitude             float32 `json:"latitude"`
	Longitude            float32 `json:"longitude"`
	HorizontalAccuracy   float32 `json:"horizontal_accuracy,omitempty"`
	LivePeriod           int     `json:"live_period,omitempty"`
	Heading              int     `json:"heading,omitempty"`
	ProximityAlertRadius int     `json:"proximity_alert_radius,omitempty"`
}

type LinkPreviewOptions struct {
	IsDisabled       bool   `json:"is_disabled,omitempty"`
	Url              string `json:"url,omitempty"`
	PreferSmallMedia bool   `json:"prefer_small_media,omitempty"`
	PreferLargeMedia bool   `json:"prefer_large_media,omitempty"`
	ShowAboveText    bool   `json:"show_above_text,omitempty"`
}

type SuggestedPostPrice struct {
	Currency string
	Amount   int
}

type PaidMediaPurchased struct {
	From             User
	PaidMediaPayload string `json:"paid_media_payload"`
}
