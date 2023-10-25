# WebhookDispatcher

## What is this?

This daemon receives **IncommingWebhooks** on localhost:8080, matches them against **Rule(s)**, and **forwards the webhooks** to destination(s)


The daemon shows the **IncommingWebhook**, **Rules matching** and **ForwardedWebhooks** in a webpage at http://localhost:8080/web with full-details of:
- the received IncommingWebhook http-request (all headers, body, etc)
- the Rules matching result (with any errors)
- the ForwardedWebhook http-response (all headers, body, etc) 


The **Rules** are defined in `./rules.yaml` [file](./rules.yaml), with each **Rule** having its own 0-or-more `jmespath_conditions`, and 0-or-more `forward_destinations`. Each **IncommingWebhook** can match with 0-or-more **Rules**.

A **Rule** matches an **IncommingWebhook**: 
- if all conditions are true
- or if there are no conditions 
- (ie, if any condition of a Rule is false, then that Rule does not match)

When a **Rule** matches, the IncommingWebhook headers and body are copied and forwarded in new http-requests to the `forward_destinations`

The `./rules.yaml` file is hot-reloaded, so it can be updated while the daemon is running.



Other features:
- Made to work with generic webhooks of content-type `application/json`
- Works transparently with webhook secret-tokens (like Github `X-Hub-Signature`), as the webhook headers and body are forwarded unchanged
- Quickly tested with github webhooks


## Usage

- Spawn webhookDispatcher and let it run in the foreground
  ```
  # Inside directory containing rules.yaml

  #LOGLEVEL can be one of: "panic", "fatal", "error", "warning", "info" (default), "debug", "trace"
  export LOGLEVEL=info
  ./webhookDispatcher
  ```

- Open http://localhost:8080/web  
  Shows history and details of Incommingwebhooks, Rules, ForwardedWebhook(s)


- Edit the `rules.yaml` file
  - file-changes are reloaded automatically (on every new IncommingWebhook)
  - the rules match against IncommingWebhook, the json-representation of each IncommingWebhook - see it in http://localhost:8080/web 
  - create your own rules `jmespath_conditions` that *must* return boolean `true` or `false` (or an error will be shown in logs/webpage to help troubleshoot) 
  - to help fine-tune `jmespath_conditions`, use the web https://play.jmespath.org/ with IncommingWebhook json
  - the `forward_destinations` can be tested using online-http-inspectors like [https://beeceptor.com](https://beeceptor.com/resources/http-echo/index.html) or https://webhook.site
  - when all the updates are saved in the `rules.yaml` file, resend a webhook and check in webpage how the Rule matched


_______


# Local development

```
# Shell 1: basic http-echo-server, can be used to send it ForwardedWebhook(s) 
docker run -p 8081:80 kennethreitz/httpbin

# Shell 2: webhookDispatcher: receives IncommingWebhook, and for each Rule matched will send ForwardedWebhook
./go_run.sh

# Shell 3: Manually simulate sending an IncommingWebhook to webhookDispatcher
clear; curl -kvX POST http://localhost:8080/post -H 'Content-Type: application/json' -d '{"myprop1":"myval1","myprop2":"myval2"}'
```

### Test with zzz.localhost.run:80 --> localhost:8080

```
# Shell 4: remote-port-forward tunnel from zzz.localhost.run:80 --> localhost:8080
# Easy: no client-install, no registration
# Some headers added, but at first-sight seems to keep original headers/body
ssh -R 80:localhost:8080 localhost.run
# -> will show the zzz.localhost.run temp-domain (ex: 138aeba52799cd.lhr.life )

# Shell 5: Manually simulate sending an IncommingWebhook to zzz.localhost.run:80 --> localhost:8080 
clear; curl -kvX POST http://138aeba52799cd.lhr.life:80/post -H 'Content-Type: application/json' -d '{"myprop1":"myval1","myprop2":"myval2"}'
```

### Test with Github-signed-webhook

```
#
#
#  Github-signed-webhook --> zzz.localhost.run:80 --> localhost:8080 --> localhost:8000
#  
#  |-------- (1) ---------|                       |------ (3) ------|
#                          |-------- (2) --------|                   |------ (4) -----|



# ----- (4) -----
# Shell-A: Listen and validate Github-signed-webhook on localhost:8000 
export GITHUB_WEBHOOK_SECRET=sillySecret
python test_stuff/validate_webhook_secret/script.py



# ----- (3) -----
# Edit rules.yaml , to add this rule:
- rule_name: github-signed-webhook
  jmespath_conditions:
    correct-method: Method == 'POST'
    header-exists: Header."X-Hub-Signature" != null
  forward_destinations:
    validateWebhookSecretProgram:
      url: http://localhost:8000/webhook

# Shell-B:
./go_run.sh

# Open http://localhost:8080/web



# ----- (2) -----
# Shell-C: 
ssh -R 80:localhost:8080 localhost.run
# -> will show the zzz.localhost.run temp-domain (ex: 138aeba52799cd.lhr.life )


# ----- (1) -----
# Github: add webhook
# - PayloadUrl: http://zzz.localhost.run    (from (2) )
# - Content-type: application/json
# - Secret: sillySecret  
# - Send me everything


# ----- Test it end-2-end ;) -----
# Github create a webhook: 
#  - go to Weebhooks/Recent-Deliveries
#  - <click one from the list> and then inside click "Redeliver" button
#  
#
# In webpage http://localhost:8080/web  check that:
# - there is a new line for **IncommingWebhook and rule "github-signed-webhook"** - click it to see details
# - in **Rule** check that `EvaluationResult` is `true`
# - in **ForwardedWebhook(s)** check that `validateWebhookSecretProgram.Status` is `200 OK`
# - in **ForwardedWebhook(s)** check that `validateWebhookSecretProgram.Body` is `Github Webhook received and validated`

```






### TODO & DONE

- [x] ./rules.yaml
- [x] conditions-in-same-rule AND'ed
- [x] improved log of rule evaluation
- [x] forward_destination -> forward_destinations []
- [x] IncommingWebhookRequest.Timestamp
- [x] env-var LOGLEVEL (default )
- [x] Works with empty-body
- [x] Works with all HTTP-verbs (GET, POST, PUT, etc)
- [x] Works with Content-type `application/json` 
- [x] Works with all signature headers
... and many others...
