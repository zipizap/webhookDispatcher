from flask import Flask, request
import os
import hmac
import hashlib


# INSTALL:
#   pip install flask
#
# USAGE:
#   GITHUB_WEBHOOK_SECRET=zzz 
#   python script.py
#
#   Forward gitlab-webhooks to POST http://localhost:8000/webhook

try:
    github_webhook_secret = os.environ["GITHUB_WEBHOOK_SECRET"]
except KeyError:
    print("Error: required env-var GITHUB_WEBHOOK_SECRET is not set")
    exit(1)

app = Flask(__name__)



@app.route('/webhook', methods=['POST'])
def respond():
    # Get the X-Hub-Signature header value
    signature = request.headers.get('X-Hub-Signature')

    # Create a hash using the payload and the secret
    secret = github_webhook_secret.encode()

    sha_name, signature = signature.split('=')
    if sha_name != 'sha1':
        return 'GithubWebhook: Invalid signature format.', 400

    mac = hmac.new(secret, msg=request.data, digestmod=hashlib.sha1)

    if not hmac.compare_digest(str(mac.hexdigest()), str(signature)):
        print('GithubWebhook: Invalid signature. 400')
        return 'GithubWebhook: Invalid signature.', 400

    # If we reach this point, the call was correctly signed
    print('Github Webhook received and validated. 200')
    return 'Github Webhook received and validated', 200

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=8000)
