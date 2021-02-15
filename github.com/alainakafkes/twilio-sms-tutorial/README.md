# Twilio

Blog post here: <https://www.twilio.com/blog/2017/09/send-text-messages-golang.html>

Code sample here: <https://github.com/alainakafkes/twilio-sms-tutorial>

Dashboard here: <https://www.twilio.com/console/sms/dashboard>

Docs here: <https://www.twilio.com/docs/sms/send-messages>

## Example `curl`

```shell
curl 'https://api.twilio.com/2010-04-01/Accounts/[account id]/Messages.json' -X POST \
--data-urlencode 'To=[phone number]' \
--data-urlencode 'MessagingServiceSid=[sid]' \
--data-urlencode 'Body=hello world' \
-u [account id]:[auth token]
```
