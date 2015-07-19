# clobbermonster
Clobbermonster is a cli load testing tool for persistent AMQP (RabbitMQ) queues written in Golang.

This app is designed mostly to enable users to test/load test a queue consumer, some form of worker processes intended to read the data sent to the RabbitMQ AMQP queue, and perform additional processing on the data.

## Installation
You can clone this repository and use it in a development capacity as you would any other Golang project, or you can download a prebuilt binary from the [releases](https://github.com/camronlevanger/clobbermonster/releases) page, and run it that way. All that is required to use the binary is Golang v1.2.1 or later.

## Usage

```
-f="~": Path of the directory containig the JSON messages to send to the queue
-i=1: The number of seconds to wait before sending messages per interval
-m=100: The number of messages to send per interval
-n="queue": The name of the persistent queue
-t=100000: The total number of messages to send to the queue
-u="amqp://guest:guest@localhost:5672/": AMQP Connection String
-P="guest": STOMP password
-T="test": STOMP topic (/topic/your_topic)
-U="guest": STOMP username
-p="amqp": use amqp or stomp protocol
-s="localhost": STOMP host
```
```
-v=0: log level for V logs
-vmodule=: comma-separated list of pattern=N settings for file-filtered logging
-alsologtostderr=false: log to standard error as well as files
-stderrthreshold=0: logs at or above this threshold go to stderr
-log_backtrace_at=:0: when logging hits line file:N, emit a stack trace
-log_dir="": If non-empty, write log files in this directory
-logtostderr=false: log to standard error instead of files
```

So, for example, the following:

`clobbermonster -p amqp -u amqp://guest:guest@localhost:5672/ -n test_clobbermonster -t 5000 -m 20 -i 1 -f /home/user/my-json-files`

will connect to an AMQP server on `localhost`, and then proceed to send `20` randomly selected messages from the directory `/home/user/my-json-files` (any number of files ending in .json) every `1` second to a queue named `test_clobbermonster` until `5000` messages have been sent.

And,

`clobbermonster -p stomp -s localhost:61613 -U guest -P guest -T /topic/test_topic -n test_clobbermonster -t 5000 -m 20 -i 1 -f /home/user/my-json-files`

does the same thing, but over STOMP instead of AMQP.

## Also

Any improvements are greatly welcomed, please open a pull request.
