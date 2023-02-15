package common

const MinBandwidth = 0
const MaxBandwidth = 100

var ProvidersMap = map[string]string{"Topolo": "Topolo", "Rond": "Rond", "Kildy": "Kildy"}
var VoiceProviderMap = map[string]string{"TransparentCalls": "TransparentCalls", "E-Voice": "E-Voice", "JustPhone": "JustPhone"}

const UrlApi = "http://172.16.238.10:8383"
const UrlMMSSystem = UrlApi + "/mms"

const AlphaColumn = 0
const BandwidthColumn = 1
const ResponseTimeColumn = 2
const ProviderColumn = 3
