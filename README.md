# Slack-api 
Библиотека для работы с api слака на golang
## Установка
```
go get gitlab.mobbtech.com/root-team/automation/slack
```
## Создание клиента
```go
var (
    token = "slack_auth_token"
    client = slack.NewClient(token, &http.Client{})
)
```

## Поиск пользователя по email
```go
var userLoader = slack.NewUserLoader(client)

var user slack.User
user, err = userLoader.LoadUserByEmail("user@gmail.com")
```

## Отправка сообщения
В качестве `RecipientID` можно использовать как ID юзера, так и канала.
```go
var messageSender = slack.NewMessageSender(client)
if err = messageSender.SendMessage(slack.Message{
    RecipientID: "CHANNEL_ID",
    Text: "Hello world!",
    AsUser:      true,
}); nil != err {
    //error
}
```
 Чтобы отвечать в тред, необходимо указать таймштамп треда
 ```go
 var messageSender = slack.NewMessageSender(client)
 if err = messageSender.SendMessage(slack.Message{
     RecipientID: "CHANNEL_ID",
     Text: "Hello world!",
     Timestamp: 100500,
     AsUser:      true,
 }); nil != err {
     //error
 }
 ```