# go-slack
Slack api client for golang applications

## Installing
```
go get github.com/kryabinin/go-slack
```

## Example
Get user information
```go
client := slack.NewClient("token")

user, err := client.GetUserByEmail(context.Background(), "test@mail.com")
if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
}

fmt.Println("user's name is: ", user.Name)
```

