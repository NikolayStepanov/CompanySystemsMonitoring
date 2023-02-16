package common

const MinBandwidth = 0
const MaxBandwidth = 100

var ProvidersMap = map[string]string{"Topolo": "Topolo", "Rond": "Rond", "Kildy": "Kildy"}
var VoiceProvidersMap = map[string]string{"TransparentCalls": "TransparentCalls", "E-Voice": "E-Voice", "JustPhone": "JustPhone"}
var EmailProvidersMap = map[string]string{"Gmail": "Gmail", "Yahoo": "Yahoo", "Hotmail": "Hotmail",
	"MSN": "MSN", "Orange": "Orange", "Comcast": "Comcast", "AOL": "AOL", "Live": "Live", "RediffMail": "RediffMail",
	"GMX": "GMX", "Protonmail": "Protonmail", "Yandex": "Yandex", "Mail.ru": "Mail.ru"}

const UrlApi = "http://172.16.238.10:8383"
const UrlMMSSystem = UrlApi + "/mms"

const (
	AlphaColumn = iota
	BandwidthColumn
	ResponseTimeColumn
	ProviderColumn
)
