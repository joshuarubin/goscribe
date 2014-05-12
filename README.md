goscribe [![wercker status](https://app.wercker.com/status/1b0a41def3a5dc3d25770d8b0e7ae909/s/ "wercker status")](https://app.wercker.com/project/bykey/1b0a41def3a5dc3d25770d8b0e7ae909) [![Coverage Status](https://coveralls.io/repos/joshuarubin/goscribe/badge.png?branch=master)](https://coveralls.io/r/joshuarubin/goscribe?branch=master)
========

Go Audio Transcription Web App

## Environment

The following environment variables are expected.

* `$BASE_URL` *MUST* be a publicly accessible URL, not something on `localhost` if you intend to test callbacks through the app itself.
    * You can use [ngrok](https://ngrok.com/) or [Runscope Passageway](https://www.runscope.com/docs/passageway) to make a localhost server public.
    * Alternatively, you can set callbacks to use [Runscope Request Capture](https://www.runscope.com/docs/request-capture) or [RequestBin](http://requestb.in/) if you just want to ensure that a callback is made and to inspect its data.
* For testing, you may also set `$TELAPI_BASE_HOST` to override the default (`api.telapi.com`). This is useful if you want to use [Runscope](https://www.runscope.com) to proxy and capture outgoing requests. Requests will always be built with the `https://` scheme.

```bash
export BASE_URL=”http://<public_server_location>”
export TELAPI_ACCOUNT_SID=”<your_telapi_account_sid>”
export TELAPI_AUTH_TOKEN=”<your_telapi_auth_token>”
export AUTH_USER=”<http_basic_auth_user>”
export AUTH_PASS=”<http_basic_auth_pass>”
```

## Heroku

For deployment on Heroku, in addition to setting the environment variables, the Go buildpack has to be set.

```bash
heroku config:add BUILDPACK_URL=https://github.com/kr/heroku-buildpack-go.git
```
