package notification_service

type MessagePayload struct {
	AppID            string                `json:"app_id"`
	Headings         PayloadHeading        `json:"headings"`
	Contents         PayloadContent        `json:"contents"`
	IncludeAliases   PayloadIncludeAliases `json:"include_aliases"`
	TargetChannel    string                `json:"target_channel"`
	AndroidChannelID string                `json:"android_channel_id"`
}

type PayloadHeading struct {
	EN string `json:"en"`
}

type PayloadContent struct {
	EN string `json:"en"`
}

type PayloadIncludeAliases struct {
	ExternalID []string `json:"external_id"`
}
