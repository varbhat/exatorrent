exatorrent comes with beautiful performant Web Client written in svelte+Typescript but also comes with documented API so that custom clients(TUI, android app ,desktop app,etc.) can be created for it.

Here lies the documented API reference.

## Authentication

exatorrent has authentication and if the connection is not authenticated, it cannot make API request.

exatorrent has user system. user is structure as defined below.

```go
type User struct {
    Username string	
    Password string 
    Token string
    UserType int // 0 for User,1 for Admin,-1 for Disabled
    CreatedAt time.Time
}
```

Length of Username and Password fields of User must be more than 5. CreatedAt field stores the creation time of User. UserType field denotes whether the User is Admin / User / Disabled.
if UserType field is -1 , User is Disabled and cannot connect to server. if UserType field is 0 , then the User is Normal User and can connect to server and perform normal operations. if UserType field is 1 ,then the User is Admin User and can connect to server. Admin User can also perform Administration Operations in addition to normal Operations.  Token field contains unique random  UUID token issued to User which can be used by User to Authenticate.

exatorrent looks whether `session_token` cookie is present in request and if the cookie is present, it validates value of `session_token` as `Token` , and if `Token` is valid, User will be authenticated.

if `session_token` cookie is not present , exatorrent looks whether Query named "token" is present in request's URL and if query named "token"  is present , it validates value of `token` query as `Token` , and if `Token` is valid , User will be authenticated.

If both `session_token` cookie and `token` query string are not present OR, value presented in cookie or query is invalid, `basic_auth` will be issued,  and if Basic Authentication's values filled by connection are valid, then the User will be authenticated.

After User is authenticated, `session_token` cookie with `Token` value will be set by exatorrent.

Briefly , Connection can be authenticated by Cookie with token ,or URL Query with token , or by Basic Auth(Username and Password).


exatorrent's WebSocket API endpoint `/api/socket` is protected by Authentication and only Authenticated User(Normal User or Admin User) can connect to WebSocket API. Note that Disabled User cannot connect to WebSocket API. Also, Only ONE WebSocket connection is allowed per User(existing WebSocket connection of User gets disconnected on new WebSocket connection of User)

Endpoint `/api/auth` is also provided where you can `POST` json request 

```json
{
  "data1": "yourusernamehere",
  "data2": "yourpasswordhere"
}
```

to get `Token` and `UserType` of the User. Response will be of the form ,

```json
{
  "usertype": "user",
  "session": "uniqueuuidsession"
}
```

# WebSocket API

exatorrent has WebSocket API which provides API to control exatorrent . Note that there is seperate HTTP API Endpoint to retrieve / stream Torrent Files for which websocket will not be used and is documented later.

Only Authenticated User who is not `disabled` can connect to WebSocket API . Only ONE WebSocket connection is allowed per User(existing WebSocket connection of User gets disconnected on new WebSocket connection of User).

Communication in exatorrent  WebSocket  is done through JSON . exatorrent WebSocket Only accepts JSON requests of the form,

```json
{
  "command": "commandstring",
  "data1": "data1string",
  "data2": "data2string",
  "data3": "data3string",
  "aop": 0
}
```

if json request is not of the form mentioned above ,  json 

```json
{
  "type": "resp",
  "state": "error",
  "msg":"invalid command"
}
```

will be sent back as response.


command field of json request represents command that needs to be done. `data1` , `data2`  and `data3` fields of json request represents data.  `aop` field specifies whether json request is admin request or not. if `aop` field is 1 and user is admin, then exatorrent considers it as admin command (note that if users who are not admins send this,they are not valuated). better omit `aop` field if request is not admin request.


## Stream Requests 

Some Requests namely `getalltorrents` , `gettorrents` , `gettorrentinfo`  are stream requests . They continuously send data regularly with 5 seconds gap in between. This is useful , say for showing progress of Torrent . There can be only 1 Stream at a time , so if you request new stream , old stream stops . you can also stop stream by sending request `{"command":"stopstream"}` . This stops stream if any do exist.

## Reference
Please Read [wshandler](https://github.com/varbhat/exatorrent/blob/main/internal/core/socket.go#L108) to get more details about API



