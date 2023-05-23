# Notify
A tool to send a notification to your phone when a command has finished running.

Notify uses [ntfy.sh](https://ntfy.sh/) to handle push notifications. The service is free to use, and all you need is a customized topic name that you have subscribed to on the app on your phone. The documentation for the website has examples on how to subscribe to a topic in both the [getting started section](https://docs.ntfy.sh/#step-1-get-the-app) and a section explaining how to [subscribe using the app](https://docs.ntfy.sh/subscribe/phone/). Note that **topic names are public**, so it's wise to choose something that cannot be guessed easily. I'd recommend generating some sort of random value to use as a prefix for your topic like a v4 UUID.

Once you have your topic name, simply set it in your environment using the `NOTIFY_TOPIC` environment variable, or pass it directly to `notify` using the `-t` or `--topic` flag. After that, just type out the command that you want notify to run. I would recommend adding `--` before the command to prevent notify from attempting to parse the other command's arguments.

## Example
```shell
notify --topic "notify-test" -- sleep 10
```

## Installation
Currently, the best way to install `notify` is by using Go directly:
```shell
go install github.com/cdevoogd/notify@latest
```
