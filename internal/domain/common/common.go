package common

var (
	ProvidersMap      = map[string]string{"Topolo": "Topolo", "Rond": "Rond", "Kildy": "Kildy"}
	VoiceProvidersMap = map[string]string{"TransparentCalls": "TransparentCalls", "E-Voice": "E-Voice", "JustPhone": "JustPhone"}
	EmailProvidersMap = map[string]string{"Gmail": "Gmail", "Yahoo": "Yahoo", "Hotmail": "Hotmail",
		"MSN": "MSN", "Orange": "Orange", "Comcast": "Comcast", "AOL": "AOL", "Live": "Live", "RediffMail": "RediffMail",
		"GMX": "GMX", "Protonmail": "Protonmail", "Yandex": "Yandex", "Mail.ru": "Mail.ru"}
)
var (
	UrlApi            = "http://172.16.238.10:8383"
	UrlMMSSystem      = UrlApi + "/mms"
	UrlSupportSystem  = UrlApi + "/support"
	UrlIncidentSystem = UrlApi + "/accendent"
)

const (
	MinBandwidth = 0
	MaxBandwidth = 100
)

const (
	AlphaColumn = iota
	BandwidthColumn
	ResponseTimeColumn
	ProviderColumn
)
