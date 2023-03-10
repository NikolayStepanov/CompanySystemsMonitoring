# Ceтeвой многопоточный сервис для Statuspage

## Дипломная работа Skillbox профессия Go-разработчик

Вы пришли работать разработчиком в компанию занимающуюся провайдингом современных средств коммуникации.

Компания предоставляет инструменты и API для автоматизации работы систем SMS, MMS, Голосовых звонков и Email. География клиентов распространяется на 160 стран и компания быстро растёт. Требуется всё больше ресурсов со стороны службы поддержки и было принято решение снизить количество заявок с помощью создания страниц информирования клиентов о текущем состоянии систем. С помощью этих страниц компания планирует снизить количество однотипных
вопросов и высвободить время агентов службы поддержки для решения более сложных задач.

Поэтому каждое подразделение компании самостоятельно контролирует работу поставщиков услуг в автоматизированном режиме храня эти данные. Ваша задача —финализировать проект, объединив эти данные и разработав небольшой сетевой сервис, который будет принимать запросы по сети и возвращать данные о состоянии систем компании.

## Функционал сервиса

### Сбор данных о системе SMS

Реализована функцию получения данных о состоянии системы SMS из файла формата CSV.

Файл CSV содержит 4 поля (alpha-2 код страны, пропускная способность канала от 0% до 100%, среднее время ответа в ms, название компании провайдера)

Пример содержимого файла:

```
BG;42;501;Rond
```

### Сбор данных о системе MMS

Реализована функцию получения данных о состоянии системы MMS опрашивая API системы через HTTP GET запрос

Получаемый ответ от API в формате json срез структур следующего вида:

```go
type MMSData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}
```

### Сбор данных о системе Voice Call

Реализована функцию получения данных о состоянии системы Voice Call из файла формата CSV.

Файл csv с данными Voice Call содержит 8 полей: alpha-2 код страны, текущая нагрузка в процентах, среднее время ответа, провайдер, стабильность соединения, TTFB, чистота связи, медиана длительности звонка.

Пример содержимого файла:

```
MC;25;1274;E-Voice;0.51;46;68;24
```

### Сбор данных о системе Email

Реализована функцию получения данных о состоянии системы Email из файла формата CSV. Файл содержит 3 поля: alpha-2 код страны, провайдер, среднее время доставки писем в ms.

```
RU;Gmail;581
```

### Сбор данных о системе Billing

Реализована функцию получения данных о состоянии системы Billing из файла содержащего битовую маску состояния систем. Каждый бит отвечает за состояние отдельной системы. Системы по порядку: создание клиента, оплата, выплата, платежи по подписке, контроль мошенничества, страница оплаты.

Результаты булевых операций сохранены в структуру вида:

```go
type BillingData struct {
	CreateCustomer bool `json:"create_customer"`
	Purchase       bool `json:"purchase"`
	Payout         bool `json:"payout"`
	Recurring      bool `json:"reccuring"`
	FraudControl   bool `json:"fraud_control"`
	CheckoutPage   bool `json:"checkout_page"`
}
```

### Сбор данных о системе Support

Написана функцию, которая отправляет GET запрос к API, разбирает полученный ответ в срез структур и возвращает его

Ответ API в формате json срез структур следующего вида:

```go
type SupportData struct {
	Topic         string `json:"topic"`
	ActiveTickets int    `json:"active_tickets"`
}
```

### Сбор данных о системе истории инцидентов

Написана функцию, которая отправляет GET запрос к API, разбирает полученный ответ в срез структур и возвращает его

Ответ от API в формате json срез структур следующего вида:

```go
type IncidentData struct {
	Topic  string `json:"topic"`
	Status string `json:"status"`
}
```

### Состояния систем

Полученные наборы данных о системах подготавливаются для возврата браузеру в нужном формате и количестве. Для отображения данные отсортированы и дополнительно отфильтрованы и показываются клиентам в удобном для просмотра виде.

Ответ сервиса в формате json. Структура вида:

```go
type ResultT struct {
	Status bool       `json:"status"`
	Data   ResultSetT `json:"data"`
	Error  string     `json:"error"`
}

type ResultSetT struct {
	SMS       [][]SMSData              `json:"sms"`
	MMS       [][]MMSData              `json:"mms"`
	VoiceCall []VoiceCallData          `json:"voice_call"`
	Email     map[string][][]EmailData `json:"email"`
	Billing   BillingData              `json:"billing"`
	Support   []int                    `json:"support"`
	Incidents []IncidentData           `json:"incident"`
}
```

## Запуск сервисов

```makefile
make build && make up
```

## Просмотр данных о состоянии систем

1. Запустить сервисы (simulator, monitoring)
2. Результат сервиса в виде json можно увидеть по адресу http://127.0.0.1:8080/api или открыть html страничку в каталоге web/index.html [status_page](web/index.html)

##  Используемые концепции в сервисе

- Веб-сервис на Go
- Работа с фреймворком [gorilla/mux](https://github.com/gorilla/mux)
- Работа с файловой системой **json, csv**
- Graceful Shutdown
- Многопоточность
