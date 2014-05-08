goscribe [![wercker status](https://app.wercker.com/status/1b0a41def3a5dc3d25770d8b0e7ae909/s/ "wercker status")](https://app.wercker.com/project/bykey/1b0a41def3a5dc3d25770d8b0e7ae909) [![Coverage Status](https://coveralls.io/repos/joshuarubin/goscribe/badge.png?branch=master)](https://coveralls.io/r/joshuarubin/goscribe?branch=master)
========

Go Audio Transcription Web App

## Environment

The following environment variables are expected.

`$BASE_URL` *MUST* be a publicly accessible URL, not something on `localhost` or else the transcription server callback will never be made.
You can use [ngrok](https://ngrok.com/) or [Runscope Passageway](https://www.runscope.com/docs/passageway) to make a localhost server public.

```bash
export BASE_URL=”http://<public_server_location>”
export TELAPI_ACCOUNT_SID=”<your_telapi_account_sid>”
```

## Heroku

For deployment on Heroku, in addition to setting the environment variables, the Go buildpack has to be set.

```bash
heroku config:add BUILDPACK_URL=https://github.com/kr/heroku-buildpack-go.git
```
