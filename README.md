# go-transport-queue

Queue for transporting large batch of messages with specific interval and batch size.

# Usage

For SMTP transport queue:
```
go-transport-queue --port 3000 --interval 1s --batch-size 10 --transport smtp
```

Then you can push messages to REST api:
```
curl -X POST \
  http://hostname:3000/push \
  -H 'content-type: application/json' \
  -d '{
	"recipients":["john.doe@example.com"],
	"subject":"Test email",
	"message":"Hello from mail queue"
}'
```

# Use cases

You have SMTP server with port throttling and you can send 10 messages per second.
```
go-transport-queue --interval 1s --batch-size 10 --transport smtp
```

You want to spread your notification load for FCM server to 1200 messages per minute:
```
go-transport-queue --interval 1s --batch-size 20 --transport fcm
```

# Configuration

## General

* `transport,t` (envvar `TRANSPORT`, required) - type of transport
* `data-path` (envvar `DATA_PATH`, default `/data`) - data storage location
* `batch-size,b` (envvar `BATCH_SIZE`, default `100`) - maximum number of messages per batch
* `interval,i` (envvar `INTERVAL`, default `100ms`) - time interval between batches

## Transports

Transports handle parsing incoming request and message delivery. Currently 3 transport types are supported: `log`, `smtp`, `fcm`

### log transport

This is easy transport for debugging purposes. Request body:
```
{
  "message":"log message"
}
```

### smtp transport

This is easy transport for debugging purposes. Request body:
```
{
  "from":"no-reply@example.com",
  "to":["john.doe@example.com"],
	"subject":"Test email",
	"text":"Hello from mail queue"
	"html":"<b>Hello</b> from mail queue"
}
```

Config variables:
* `smtp-url` (envvar `SMTP_URL`) - SMTP configuration in url format
* `smtp-sender` (envvar `SMTP_SENDER`) - sender of emails (rfc2047)

### fcm transport

Firebase Cloud Messaging transport. Request body:
```
{
	"recipients":["fcm_token1","fcm_token2"],
	"data":{
    "xx":"aa"
  },
	"notification":{
    "title":"transport queue",
    "body":"hello world",
    "icon":"default.png",
    "badge":"1"
    "sound":"default",
    "color":"",
    "click_action":"",
    "body_loc_key":"",
    "body_loc_args":"",
    "title_loc_key":"",
    "title_loc_args":""
  }
}
```

Config variables:
* `fcm-api-key` (envvar `FCM_API_KEY`) - Google FCM Api Key

# Docker

You can start container by running:
```
docker run --rm -p 3000:80 -v /path/to/data:/data go-transport-queue --transport log --batch-size 10 --interval 1s
```

Or by specifying environment variables:
```
docker run --rm -p 3000:80 -e TRANSPORT=log -e BATCH_SIZE=10 -e INTERVAL=1s -v /path/to/data:/data go-transport-queue
```
