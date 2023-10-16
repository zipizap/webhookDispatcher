# WebhookDispatcher


- Should work transparently with webhook secret-token as the headers and body are forwarded unchanged
- Hot-reload of rules.yaml file changes
- Works with webhooks of content-type `application/json`, untested with `application/x-www-form-urlencoded` 


## Usage

- Edit `rules.yaml`
- Spawn webhookDispatcher and let it run in the foreground
- [opt] Update `rules.yaml` at any time - its reloaded on each IncommingWebhook

## Local development

```
docker run -p 8081:80 kennethreitz/httpbin

./go_run.sh

clear; curl -kvX POST http://localhost:8080/post -H 'Content-Type: application/json' -d '{"myprop1":"myval1","myprop2":"myval2"}'

```

### TODO & DONE

[x] ./rules.yaml
[x] conditions-in-same-rule AND'ed
[x] improved log of rule evaluation
[x] forward_destination -> forward_destinations []
[ ] webpage
