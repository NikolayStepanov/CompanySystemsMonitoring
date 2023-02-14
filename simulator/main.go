package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const minResponseTime = 30
const maxResponseTime = 2000

const minConnectionStability = 600
const maxConnectionStability = 1000

const minVoicePurity = 0
const maxVoicePurity = 92

const minVoiceCallMedian = 3
const maxVoiceCallMedian = 60

const minTTFB = 2
const maxTTFB = 980

const minBandwidth = 0
const maxBandwidth = 100

const minEmailDeliveryTime = 0
const maxEmailDeliveryTime = 600

const smsFilename = "sms.data"
const voiceFilename = "voice.data"
const emailFilename = "email.data"
const billingFilename = "billing.data"
const accendentListFilename = "accendents.data"

var firstSMSRowForCorrupt int
var secondSMSRowForCorrupt int

var firstVoiceRowForCorrupt int
var secondVoiceRowForCorrupt int

var firstEmailRowForCorrupt int
var secondEmailRowForCorrupt int

var MMSCollection []MMSItem
var SupportCollection []SupportItem
var AccendentCollection []AccendentItem

type MMSItem struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

type SupportItem struct {
	Topic         string `json:"topic"`
	ActiveTickets int    `json:"active_tickets"`
}

type AccendentItem struct {
	Topic  string `json:"topic"`
	Status string `json:"status"`
}

const accendentStatusActive = "active"
const accendentStatusClosed = "closed"

var AccendentTopics = []string{
	"SMS delivery in EU",
	"MMS connection stability",
	"Voice call connection purity",
	"Checkout page is down",
	"Support overload",
	"Buy phone number not working in US",
	"API Slow latency",
}

func init() {
	rand.Seed(time.Now().UnixNano())

	firstSMSRowForCorrupt = rand.Intn(70)
	fmt.Printf("First SMS row for currupt %d\n", firstSMSRowForCorrupt+1)

	secondSMSRowForCorrupt = rand.Intn(90)
	fmt.Printf("Second SMS row for currupt %d\n", secondSMSRowForCorrupt+1)

	firstVoiceRowForCorrupt = rand.Intn(70)
	fmt.Printf("First Voice row for currupt %d\n", firstVoiceRowForCorrupt+1)

	secondVoiceRowForCorrupt = rand.Intn(90)
	fmt.Printf("Second Voice row for currupt %d\n", secondVoiceRowForCorrupt+1)

	firstEmailRowForCorrupt = rand.Intn(70)
	fmt.Printf("First Email row for currupt %d\n", firstEmailRowForCorrupt+1)

	secondEmailRowForCorrupt = rand.Intn(90)
	fmt.Printf("Second Email row for currupt %d\n", secondEmailRowForCorrupt+1)
}

func main() {
	shuffleSmsData()

	MMSCollection = shuffleMMSData()

	shuffleVoiceData()
	shuffleEmailData()
	shuffleBillingData()

	SupportCollection = shuffleSupportData()
	AccendentCollection = shuffleAccendentData()

	listenAndServeHTTP()
}

func shuffleSmsData() {
	var data string
	for i, country := range getCountriesList() {
		row := strings.Join([]string{
			country,
			getRandomBandwidthInString(),
			getRandomResponseTimeInString(),
			getSmsProviderByCountry(country),
		}, ";") + "\n"

		if i == firstSMSRowForCorrupt || i == secondSMSRowForCorrupt {
			row = strings.Replace(row, ";", "", rand.Intn(4))
			row = strings.Replace(row, "R", "", rand.Intn(3))
			row = strings.Replace(row, "C", "", rand.Intn(3))

			fmt.Println("SMS row corrupted")
		}

		data += row
	}

	err := ioutil.WriteFile(getFilapathByFilename(smsFilename), []byte(data), 0644)
	if err != nil {
		fmt.Printf("Error in write sms data: %s", err.Error())
	}
}

func shuffleMMSData() []MMSItem {
	data := make([]MMSItem, 0)
	for _, country := range getCountriesList() {
		data = append(
			data,
			MMSItem{
				Country:      country,
				Provider:     getMMSProviderByCountry(country),
				Bandwidth:    getRandomBandwidthInString(),
				ResponseTime: getRandomResponseTimeInString(),
			},
		)
	}

	return data
}

func shuffleVoiceData() {
	var data string
	for i, country := range getCountriesList() {
		row := strings.Join([]string{
			country,
			getRandomBandwidthInString(),
			getRandomResponseTimeInString(),
			getVoiceCallProviderByCountry(country),
			getRandomConnectionStability(),
			getRandomTTFB(),
			getRandomVoicePurity(),
			getRandomMedianOfCallsTime(),
		}, ";") + "\n"

		if i == firstVoiceRowForCorrupt || i == secondVoiceRowForCorrupt {
			row = strings.Replace(row, ";", "", rand.Intn(4))
			row = strings.Replace(row, "R", "", rand.Intn(3))
			row = strings.Replace(row, "C", "", rand.Intn(3))

			fmt.Println("Voice row corrupted")
		}

		data += row
	}

	err := ioutil.WriteFile(getFilapathByFilename(voiceFilename), []byte(data), 0644)
	if err != nil {
		fmt.Printf("Error in write sms data: %s", err.Error())
	}
}

func shuffleEmailData() {
	var data string
	providersList := getEmailProvidersList()
	i := 0
	for _, country := range getCountriesList() {
		for _, provider := range providersList {
			row := strings.Join([]string{
				country,
				provider,
				getRandomEmailDeliveryTime(),
			}, ";") + "\n"

			if i == firstEmailRowForCorrupt || i == secondEmailRowForCorrupt {
				row = strings.Replace(row, ";", "", rand.Intn(4))
				row = strings.Replace(row, "A", "", rand.Intn(3))
				row = strings.Replace(row, "a", "", rand.Intn(3))
				row = strings.Replace(row, "O", "", rand.Intn(3))
				row = strings.Replace(row, "o", "", rand.Intn(3))
				row = strings.Replace(row, "M", "", rand.Intn(3))
				row = strings.Replace(row, "m", "", rand.Intn(3))
				row = strings.Replace(row, "P", "", rand.Intn(3))
				row = strings.Replace(row, "p", "", rand.Intn(3))

				fmt.Println("Email row corrupted")
			}

			data += row
			i++
		}
	}

	err := ioutil.WriteFile(getFilapathByFilename(emailFilename), []byte(data), 0644)
	if err != nil {
		fmt.Printf("Error in write email data: %s", err.Error())
	}
}

func shuffleBillingData() {
	data := ""
	for i := 0; i < 6; i++ {
		value := getRandomIntBetweenValues(0, 150)
		if value > 50 {
			value = 1
		} else {
			value = 0
		}

		data = data + fmt.Sprintf("%d", value)
		// create customer
		// purchase
		// payout
		// recurring
		// fraud control
		// checkout page
	}

	err := ioutil.WriteFile(getFilapathByFilename(billingFilename), []byte(data), 0644)
	if err != nil {
		fmt.Printf("Error in write sms data: %s", err.Error())
	}
}

func shuffleSupportData() []SupportItem {
	data := make([]SupportItem, 0)
	for _, topic := range getSupportTopicsList() {
		data = append(data, SupportItem{Topic: topic, ActiveTickets: getRandomSupportTickets()})
	}

	return data
}

func shuffleAccendentData() []AccendentItem {
	collection := make([]AccendentItem, 0)
	status := ""
	for _, topic := range AccendentTopics {
		if getRandomIntBetweenValues(0, 2) == 1 {
			status = accendentStatusActive
		} else {
			status = accendentStatusClosed
		}

		collection = append(collection, AccendentItem{Topic: topic, Status: status})
	}

	return collection
}

func getCountriesList() []string {
	return []string{"RU", "US", "GB", "FR", "BL", "AT", "BG", "DK", "CA", "ES", "CH", "TR", "PE", "NZ", "MC"}
}

func getSmsProviderByCountry(country string) string {
	smsProviderMap := map[string]string{
		"RU": "Topolo",
		"US": "Rond",
		"GB": "Topolo",
		"FR": "Topolo",
		"BL": "Kildy",
		"AT": "Topolo",
		"BG": "Rond",
		"DK": "Topolo",
		"CA": "Rond",
		"ES": "Topolo",
		"CH": "Topolo",
		"TR": "Rond",
		"PE": "Topolo",
		"NZ": "Kildy",
		"MC": "Kildy",
	}

	return smsProviderMap[country]
}

func getMMSProviderByCountry(country string) string {
	smsProviderMap := map[string]string{
		"RU": "Topolo",
		"US": "Rond",
		"GB": "Topolo",
		"FR": "Topolo",
		"BL": "Kildy",
		"AT": "Topolo",
		"BG": "Rond",
		"DK": "Topolo",
		"CA": "Rond",
		"ES": "Topolo",
		"CH": "Topolo",
		"TR": "Rond",
		"PE": "Topolo",
		"NZ": "Kildy",
		"MC": "Kildy",
	}

	return smsProviderMap[country]
}

func getVoiceCallProviderByCountry(country string) string {
	voiceProviderMap := map[string]string{
		"RU": "TransparentCalls",
		"US": "E-Voice",
		"GB": "TransparentCalls",
		"FR": "TransparentCalls",
		"BL": "E-Voice",
		"AT": "TransparentCalls",
		"BG": "E-Voice",
		"DK": "JustPhone",
		"CA": "JustPhone",
		"ES": "E-Voice",
		"CH": "JustPhone",
		"TR": "TransparentCalls",
		"PE": "JustPhone",
		"NZ": "JustPhone",
		"MC": "E-Voice",
	}

	return voiceProviderMap[country]
}

func getEmailProvidersList() []string {
	return []string{
		"Gmail",
		"Yahoo",
		"Hotmail",
		"MSN",
		"Orange",
		"Comcast",
		"AOL",
		"Live",
		"RediffMail",
		"GMX",
		"Protonmail",
		"Yandex",
		"Mail.ru",
	}
}

func getSupportTopicsList() []string {
	return []string{
		"SMS",
		"MMS",
		"Email",
		"Billing",
		"Create account",
		"API",
		"Marketing",
		"Privacy",
		"GDPR",
		"Other",
	}
}

func getRandomSupportTickets() int {
	return getRandomIntBetweenValues(0, 8)
}

func getFilapathByFilename(filename string) string {
	return "data/" + filename
}

func getRandomBandwidthInString() string {
	return strconv.Itoa(getRandomIntBetweenValues(minBandwidth, maxBandwidth))
}

func getRandomResponseTimeInString() string {
	return strconv.Itoa(getRandomIntBetweenValues(minResponseTime, maxResponseTime))
}

func getRandomConnectionStability() string {
	stability := getRandomIntBetweenValues(minConnectionStability, maxConnectionStability)

	return fmt.Sprintf("%.2f", float32(stability)/1000)
}

func getRandomTTFB() string {
	return strconv.Itoa(getRandomIntBetweenValues(minTTFB, maxTTFB))
}

func getRandomVoicePurity() string {
	return strconv.Itoa(getRandomIntBetweenValues(minVoicePurity, maxVoicePurity))
}

func getRandomMedianOfCallsTime() string {
	return strconv.Itoa(getRandomIntBetweenValues(minVoiceCallMedian, maxVoiceCallMedian))
}

func getRandomEmailDeliveryTime() string {
	return strconv.Itoa(getRandomIntBetweenValues(minEmailDeliveryTime, maxEmailDeliveryTime))
}

func getRandomIntBetweenValues(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func listenAndServeHTTP() {
	router := mux.NewRouter()

	router.HandleFunc("/mms", handleMMS)
	router.HandleFunc("/support", handleSupport)
	router.HandleFunc("/accendent", handleAccendent)
	router.HandleFunc("/test", handleTest).Methods("GET", "OPTIONS")

	http.ListenAndServe("127.0.0.1:8383", router)
}

func handleMMS(w http.ResponseWriter, r *http.Request) {
	response(w, r, MMSCollection)
}

func handleSupport(w http.ResponseWriter, r *http.Request) {
	response(w, r, SupportCollection)
}

func handleAccendent(w http.ResponseWriter, r *http.Request) {
	response(w, r, AccendentCollection)
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Write(testBodyResponse)
}

func response(w http.ResponseWriter, r *http.Request, responseStruct interface{}) {
	response, _ := json.Marshal(responseStruct)

	w.Write(response)
}

var testBodyResponse = []byte(`
{
    "status": true,
    "data_app": {
      "sms": [
        [
          {
            "country": "Saint Barthélemy",
            "bandwidth": "90",
            "response_time": "1615",
            "provider": "Kildy"
          },
          {
            "country": "New Zealand",
            "bandwidth": "67",
            "response_time": "505",
            "provider": "Kildy"
          },
          {
            "country": "Monaco",
            "bandwidth": "18",
            "response_time": "1022",
            "provider": "Kildy"
          },
          {
            "country": "United States",
            "bandwidth": "68",
            "response_time": "561",
            "provider": "Rond"
          },
          {
            "country": "Bulgaria",
            "bandwidth": "9",
            "response_time": "1155",
            "provider": "Rond"
          },
          {
            "country": "Canada",
            "bandwidth": "6",
            "response_time": "1629",
            "provider": "Rond"
          },
          {
            "country": "Turkey",
            "bandwidth": "5",
            "response_time": "1689",
            "provider": "Rond"
          },
          {
            "country": "Russian Federation",
            "bandwidth": "83",
            "response_time": "1070",
            "provider": "Topolo"
          },
          {
            "country": "United Kingdom",
            "bandwidth": "49",
            "response_time": "241",
            "provider": "Topolo"
          },
          {
            "country": "France",
            "bandwidth": "76",
            "response_time": "346",
            "provider": "Topolo"
          },
          {
            "country": "Austria",
            "bandwidth": "40",
            "response_time": "38",
            "provider": "Topolo"
          },
          {
            "country": "Denmark",
            "bandwidth": "56",
            "response_time": "450",
            "provider": "Topolo"
          },
          {
            "country": "Spain",
            "bandwidth": "93",
            "response_time": "906",
            "provider": "Topolo"
          },
          {
            "country": "Peru",
            "bandwidth": "92",
            "response_time": "1670",
            "provider": "Topolo"
          }
        ],
        [
          {
            "country": "Austria",
            "bandwidth": "40",
            "response_time": "38",
            "provider": "Topolo"
          },
          {
            "country": "Bulgaria",
            "bandwidth": "9",
            "response_time": "1155",
            "provider": "Rond"
          },
          {
            "country": "Canada",
            "bandwidth": "6",
            "response_time": "1629",
            "provider": "Rond"
          },
          {
            "country": "Denmark",
            "bandwidth": "56",
            "response_time": "450",
            "provider": "Topolo"
          },
          {
            "country": "France",
            "bandwidth": "76",
            "response_time": "346",
            "provider": "Topolo"
          },
          {
            "country": "Monaco",
            "bandwidth": "18",
            "response_time": "1022",
            "provider": "Kildy"
          },
          {
            "country": "New Zealand",
            "bandwidth": "67",
            "response_time": "505",
            "provider": "Kildy"
          },
          {
            "country": "Peru",
            "bandwidth": "92",
            "response_time": "1670",
            "provider": "Topolo"
          },
          {
            "country": "Russian Federation",
            "bandwidth": "83",
            "response_time": "1070",
            "provider": "Topolo"
          },
          {
            "country": "Saint Barthélemy",
            "bandwidth": "90",
            "response_time": "1615",
            "provider": "Kildy"
          },
          {
            "country": "Spain",
            "bandwidth": "93",
            "response_time": "906",
            "provider": "Topolo"
          },
          {
            "country": "Turkey",
            "bandwidth": "5",
            "response_time": "1689",
            "provider": "Rond"
          },
          {
            "country": "United Kingdom",
            "bandwidth": "49",
            "response_time": "241",
            "provider": "Topolo"
          },
          {
            "country": "United States",
            "bandwidth": "68",
            "response_time": "561",
            "provider": "Rond"
          }
        ]
      ],
      "mms": [
        [
          {
            "country": "Saint Barthélemy",
            "provider": "Kildy",
            "bandwidth": "35",
            "response_time": "1621"
          },
          {
            "country": "New Zealand",
            "provider": "Kildy",
            "bandwidth": "7",
            "response_time": "1785"
          },
          {
            "country": "Monaco",
            "provider": "Kildy",
            "bandwidth": "6",
            "response_time": "969"
          },
          {
            "country": "United States",
            "provider": "Rond",
            "bandwidth": "41",
            "response_time": "256"
          },
          {
            "country": "Bulgaria",
            "provider": "Rond",
            "bandwidth": "56",
            "response_time": "564"
          },
          {
            "country": "Canada",
            "provider": "Rond",
            "bandwidth": "12",
            "response_time": "1518"
          },
          {
            "country": "Turkey",
            "provider": "Rond",
            "bandwidth": "82",
            "response_time": "1459"
          },
          {
            "country": "Russian Federation",
            "provider": "Topolo",
            "bandwidth": "3",
            "response_time": "1689"
          },
          {
            "country": "United Kingdom",
            "provider": "Topolo",
            "bandwidth": "63",
            "response_time": "1818"
          },
          {
            "country": "France",
            "provider": "Topolo",
            "bandwidth": "83",
            "response_time": "955"
          },
          {
            "country": "Austria",
            "provider": "Topolo",
            "bandwidth": "28",
            "response_time": "305"
          },
          {
            "country": "Denmark",
            "provider": "Topolo",
            "bandwidth": "100",
            "response_time": "384"
          },
          {
            "country": "Spain",
            "provider": "Topolo",
            "bandwidth": "37",
            "response_time": "1238"
          },
          {
            "country": "Switzerland",
            "provider": "Topolo",
            "bandwidth": "89",
            "response_time": "381"
          },
          {
            "country": "Peru",
            "provider": "Topolo",
            "bandwidth": "83",
            "response_time": "1714"
          }
        ],
        [
          {
            "country": "Austria",
            "provider": "Topolo",
            "bandwidth": "28",
            "response_time": "305"
          },
          {
            "country": "Bulgaria",
            "provider": "Rond",
            "bandwidth": "56",
            "response_time": "564"
          },
          {
            "country": "Canada",
            "provider": "Rond",
            "bandwidth": "12",
            "response_time": "1518"
          },
          {
            "country": "Denmark",
            "provider": "Topolo",
            "bandwidth": "100",
            "response_time": "384"
          },
          {
            "country": "France",
            "provider": "Topolo",
            "bandwidth": "83",
            "response_time": "955"
          },
          {
            "country": "Monaco",
            "provider": "Kildy",
            "bandwidth": "6",
            "response_time": "969"
          },
          {
            "country": "New Zealand",
            "provider": "Kildy",
            "bandwidth": "7",
            "response_time": "1785"
          },
          {
            "country": "Peru",
            "provider": "Topolo",
            "bandwidth": "83",
            "response_time": "1714"
          },
          {
            "country": "Russian Federation",
            "provider": "Topolo",
            "bandwidth": "3",
            "response_time": "1689"
          },
          {
            "country": "Saint Barthélemy",
            "provider": "Kildy",
            "bandwidth": "35",
            "response_time": "1621"
          },
          {
            "country": "Spain",
            "provider": "Topolo",
            "bandwidth": "37",
            "response_time": "1238"
          },
          {
            "country": "Switzerland",
            "provider": "Topolo",
            "bandwidth": "89",
            "response_time": "381"
          },
          {
            "country": "Turkey",
            "provider": "Rond",
            "bandwidth": "82",
            "response_time": "1459"
          },
          {
            "country": "United Kingdom",
            "provider": "Topolo",
            "bandwidth": "63",
            "response_time": "1818"
          },
          {
            "country": "United States",
            "provider": "Rond",
            "bandwidth": "41",
            "response_time": "256"
          }
        ]
      ],
      "voice_call": [
        {
          "country": "Russian Federation",
          "bandwidth": "63",
          "response_time": "1842",
          "provider": "TransparentCalls",
          "connection_stability": 0.88,
          "ttfb": 409,
          "voice_purity": 26,
          "median_of_calls_time": 7
        },
        {
          "country": "United States",
          "bandwidth": "90",
          "response_time": "1414",
          "provider": "E-Voice",
          "connection_stability": 0.83,
          "ttfb": 755,
          "voice_purity": 32,
          "median_of_calls_time": 39
        },
        {
          "country": "United Kingdom",
          "bandwidth": "16",
          "response_time": "894",
          "provider": "TransparentCalls",
          "connection_stability": 0.91,
          "ttfb": 925,
          "voice_purity": 74,
          "median_of_calls_time": 5
        },
        {
          "country": "France",
          "bandwidth": "77",
          "response_time": "1039",
          "provider": "TransparentCalls",
          "connection_stability": 0.63,
          "ttfb": 478,
          "voice_purity": 0,
          "median_of_calls_time": 25
        },
        {
          "country": "Saint Barthélemy",
          "bandwidth": "61",
          "response_time": "394",
          "provider": "E-Voice",
          "connection_stability": 0.97,
          "ttfb": 544,
          "voice_purity": 87,
          "median_of_calls_time": 33
        },
        {
          "country": "Austria",
          "bandwidth": "81",
          "response_time": "371",
          "provider": "TransparentCalls",
          "connection_stability": 0.66,
          "ttfb": 481,
          "voice_purity": 84,
          "median_of_calls_time": 23
        },
        {
          "country": "Bulgaria",
          "bandwidth": "5",
          "response_time": "573",
          "provider": "E-Voice",
          "connection_stability": 0.7,
          "ttfb": 283,
          "voice_purity": 57,
          "median_of_calls_time": 36
        },
        {
          "country": "Denmark",
          "bandwidth": "47",
          "response_time": "1721",
          "provider": "JustPhone",
          "connection_stability": 0.74,
          "ttfb": 263,
          "voice_purity": 90,
          "median_of_calls_time": 57
        },
        {
          "country": "Canada",
          "bandwidth": "40",
          "response_time": "1115",
          "provider": "JustPhone",
          "connection_stability": 0.64,
          "ttfb": 30,
          "voice_purity": 25,
          "median_of_calls_time": 25
        },
        {
          "country": "Spain",
          "bandwidth": "4",
          "response_time": "1587",
          "provider": "E-Voice",
          "connection_stability": 0.83,
          "ttfb": 385,
          "voice_purity": 0,
          "median_of_calls_time": 43
        },
        {
          "country": "Turkey",
          "bandwidth": "28",
          "response_time": "1954",
          "provider": "TransparentCalls",
          "connection_stability": 0.69,
          "ttfb": 58,
          "voice_purity": 47,
          "median_of_calls_time": 44
        },
        {
          "country": "Peru",
          "bandwidth": "97",
          "response_time": "781",
          "provider": "JustPhone",
          "connection_stability": 0.92,
          "ttfb": 206,
          "voice_purity": 74,
          "median_of_calls_time": 3
        },
        {
          "country": "New Zealand",
          "bandwidth": "44",
          "response_time": "1730",
          "provider": "JustPhone",
          "connection_stability": 0.64,
          "ttfb": 845,
          "voice_purity": 7,
          "median_of_calls_time": 58
        },
        {
          "country": "Monaco",
          "bandwidth": "78",
          "response_time": "1856",
          "provider": "E-Voice",
          "connection_stability": 0.81,
          "ttfb": 542,
          "voice_purity": 6,
          "median_of_calls_time": 42
        }
      ],
      "email": {
        "Austria": [
          [
            {
              "country": "AT",
              "provider": "Hotmail",
              "delivery_time": 95
            },
            {
              "country": "AT",
              "provider": "Yandex",
              "delivery_time": 105
            },
            {
              "country": "AT",
              "provider": "Mail.ru",
              "delivery_time": 169
            }
          ],
          [
            {
              "country": "AT",
              "provider": "AOL",
              "delivery_time": 393
            },
            {
              "country": "AT",
              "provider": "GMX",
              "delivery_time": 471
            },
            {
              "country": "AT",
              "provider": "Comcast",
              "delivery_time": 495
            }
          ]
        ],
        "Bulgaria": [
          [
            {
              "country": "BG",
              "provider": "Yandex",
              "delivery_time": 0
            },
            {
              "country": "BG",
              "provider": "Comcast",
              "delivery_time": 3
            },
            {
              "country": "BG",
              "provider": "GMX",
              "delivery_time": 5
            }
          ],
          [
            {
              "country": "BG",
              "provider": "Hotmail",
              "delivery_time": 375
            },
            {
              "country": "BG",
              "provider": "MSN",
              "delivery_time": 467
            },
            {
              "country": "BG",
              "provider": "AOL",
              "delivery_time": 544
            }
          ]
        ],
        "Canada": [
          [
            {
              "country": "CA",
              "provider": "Comcast",
              "delivery_time": 22
            },
            {
              "country": "CA",
              "provider": "GMX",
              "delivery_time": 90
            },
            {
              "country": "CA",
              "provider": "Hotmail",
              "delivery_time": 93
            }
          ],
          [
            {
              "country": "CA",
              "provider": "Live",
              "delivery_time": 474
            },
            {
              "country": "CA",
              "provider": "Mail.ru",
              "delivery_time": 521
            },
            {
              "country": "CA",
              "provider": "MSN",
              "delivery_time": 553
            }
          ]
        ],
        "Denmark": [
          [
            {
              "country": "DK",
              "provider": "Gmail",
              "delivery_time": 135
            },
            {
              "country": "DK",
              "provider": "Yahoo",
              "delivery_time": 246
            },
            {
              "country": "DK",
              "provider": "Hotmail",
              "delivery_time": 256
            }
          ],
          [
            {
              "country": "DK",
              "provider": "Yandex",
              "delivery_time": 464
            },
            {
              "country": "DK",
              "provider": "Mail.ru",
              "delivery_time": 519
            },
            {
              "country": "DK",
              "provider": "Comcast",
              "delivery_time": 562
            }
          ]
        ],
        "France": [
          [
            {
              "country": "FR",
              "provider": "MSN",
              "delivery_time": 6
            },
            {
              "country": "FR",
              "provider": "Hotmail",
              "delivery_time": 51
            },
            {
              "country": "FR",
              "provider": "Comcast",
              "delivery_time": 75
            }
          ],
          [
            {
              "country": "FR",
              "provider": "Gmail",
              "delivery_time": 469
            },
            {
              "country": "FR",
              "provider": "RediffMail",
              "delivery_time": 484
            },
            {
              "country": "FR",
              "provider": "Live",
              "delivery_time": 526
            }
          ]
        ],
        "Monaco": [
          [
            {
              "country": "MC",
              "provider": "Yandex",
              "delivery_time": 17
            },
            {
              "country": "MC",
              "provider": "Live",
              "delivery_time": 102
            },
            {
              "country": "MC",
              "provider": "AOL",
              "delivery_time": 127
            }
          ],
          [
            {
              "country": "MC",
              "provider": "Protonmail",
              "delivery_time": 532
            },
            {
              "country": "MC",
              "provider": "RediffMail",
              "delivery_time": 534
            },
            {
              "country": "MC",
              "provider": "Mail.ru",
              "delivery_time": 559
            }
          ]
        ],
        "New Zealand": [
          [
            {
              "country": "NZ",
              "provider": "Comcast",
              "delivery_time": 51
            },
            {
              "country": "NZ",
              "provider": "Gmail",
              "delivery_time": 81
            },
            {
              "country": "NZ",
              "provider": "RediffMail",
              "delivery_time": 85
            }
          ],
          [
            {
              "country": "NZ",
              "provider": "AOL",
              "delivery_time": 577
            },
            {
              "country": "NZ",
              "provider": "Hotmail",
              "delivery_time": 589
            },
            {
              "country": "NZ",
              "provider": "Orange",
              "delivery_time": 592
            }
          ]
        ],
        "Peru": [
          [
            {
              "country": "PE",
              "provider": "Gmail",
              "delivery_time": 5
            },
            {
              "country": "PE",
              "provider": "Comcast",
              "delivery_time": 23
            },
            {
              "country": "PE",
              "provider": "Protonmail",
              "delivery_time": 148
            }
          ],
          [
            {
              "country": "PE",
              "provider": "Orange",
              "delivery_time": 540
            },
            {
              "country": "PE",
              "provider": "Yahoo",
              "delivery_time": 554
            },
            {
              "country": "PE",
              "provider": "Yandex",
              "delivery_time": 590
            }
          ]
        ],
        "Russian Federation": [
          [
            {
              "country": "RU",
              "provider": "AOL",
              "delivery_time": 32
            },
            {
              "country": "RU",
              "provider": "GMX",
              "delivery_time": 34
            },
            {
              "country": "RU",
              "provider": "MSN",
              "delivery_time": 236
            }
          ],
          [
            {
              "country": "RU",
              "provider": "Orange",
              "delivery_time": 491
            },
            {
              "country": "RU",
              "provider": "Gmail",
              "delivery_time": 546
            },
            {
              "country": "RU",
              "provider": "Mail.ru",
              "delivery_time": 570
            }
          ]
        ],
        "Saint Barthélemy": [
          [
            {
              "country": "BL",
              "provider": "Yandex",
              "delivery_time": 0
            },
            {
              "country": "BL",
              "provider": "AOL",
              "delivery_time": 48
            },
            {
              "country": "BL",
              "provider": "Live",
              "delivery_time": 108
            }
          ],
          [
            {
              "country": "BL",
              "provider": "Hotmail",
              "delivery_time": 440
            },
            {
              "country": "BL",
              "provider": "Yahoo",
              "delivery_time": 526
            },
            {
              "country": "BL",
              "provider": "MSN",
              "delivery_time": 553
            }
          ]
        ],
        "Spain": [
          [
            {
              "country": "ES",
              "provider": "Yandex",
              "delivery_time": 6
            },
            {
              "country": "ES",
              "provider": "Protonmail",
              "delivery_time": 42
            },
            {
              "country": "ES",
              "provider": "MSN",
              "delivery_time": 81
            }
          ],
          [
            {
              "country": "ES",
              "provider": "Live",
              "delivery_time": 356
            },
            {
              "country": "ES",
              "provider": "Hotmail",
              "delivery_time": 542
            },
            {
              "country": "ES",
              "provider": "Orange",
              "delivery_time": 592
            }
          ]
        ],
        "Switzerland": [
          [
            {
              "country": "CH",
              "provider": "AOL",
              "delivery_time": 29
            },
            {
              "country": "CH",
              "provider": "Hotmail",
              "delivery_time": 134
            },
            {
              "country": "CH",
              "provider": "Yandex",
              "delivery_time": 250
            }
          ],
          [
            {
              "country": "CH",
              "provider": "Live",
              "delivery_time": 447
            },
            {
              "country": "CH",
              "provider": "Yahoo",
              "delivery_time": 558
            },
            {
              "country": "CH",
              "provider": "Protonmail",
              "delivery_time": 580
            }
          ]
        ],
        "Turkey": [
          [
            {
              "country": "TR",
              "provider": "Yahoo",
              "delivery_time": 90
            },
            {
              "country": "TR",
              "provider": "GMX",
              "delivery_time": 98
            },
            {
              "country": "TR",
              "provider": "AOL",
              "delivery_time": 129
            }
          ],
          [
            {
              "country": "TR",
              "provider": "Gmail",
              "delivery_time": 487
            },
            {
              "country": "TR",
              "provider": "MSN",
              "delivery_time": 494
            },
            {
              "country": "TR",
              "provider": "RediffMail",
              "delivery_time": 494
            }
          ]
        ],
        "United Kingdom": [
          [
            {
              "country": "GB",
              "provider": "AOL",
              "delivery_time": 83
            },
            {
              "country": "GB",
              "provider": "Mail.ru",
              "delivery_time": 170
            },
            {
              "country": "GB",
              "provider": "Gmail",
              "delivery_time": 175
            }
          ],
          [
            {
              "country": "GB",
              "provider": "Yandex",
              "delivery_time": 545
            },
            {
              "country": "GB",
              "provider": "Yahoo",
              "delivery_time": 560
            },
            {
              "country": "GB",
              "provider": "MSN",
              "delivery_time": 571
            }
          ]
        ],
        "United States": [
          [
            {
              "country": "US",
              "provider": "Gmail",
              "delivery_time": 72
            },
            {
              "country": "US",
              "provider": "MSN",
              "delivery_time": 82
            },
            {
              "country": "US",
              "provider": "Live",
              "delivery_time": 120
            }
          ],
          [
            {
              "country": "US",
              "provider": "Protonmail",
              "delivery_time": 558
            },
            {
              "country": "US",
              "provider": "Hotmail",
              "delivery_time": 565
            },
            {
              "country": "US",
              "provider": "Mail.ru",
              "delivery_time": 597
            }
          ]
        ]
      },
      "billing": {
        "create_customer": true,
        "purchase": false,
        "payout": true,
        "reccuring": true,
        "fraud_control": true,
        "checkout_page": false
      },
      "support": [
        3,
        144
      ],
      "incident": [
        {
          "topic": "SMS delivery in EU",
          "status": "active"
        },
        {
          "topic": "MMS connection stability",
          "status": "active"
        },
        {
          "topic": "Checkout page is down",
          "status": "active"
        },
        {
          "topic": "Support overload",
          "status": "active"
        },
        {
          "topic": "Voice call connection purity",
          "status": "closed"
        },
        {
          "topic": "Buy phone number not working in US",
          "status": "closed"
        },
        {
          "topic": "API Slow latency",
          "status": "closed"
        }
      ]
    },
    "error": ""
  }	
`)
