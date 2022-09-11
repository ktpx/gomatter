# gomatter

Command line Utility to send notification/massages via Mattermost webhooks.

```
$ gomatter -m "hello, world!" -c world
```

Send text to channel "world".  Default channel is "town-square".

You can set the webhook URL and default channel in the code, or
specifiy via environment variable MM_WEBHOOKURL (or pass to 
command via -w option).

Messages can be read from STDIN with -r.  Check help for more
info.




