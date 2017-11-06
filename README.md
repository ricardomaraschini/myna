#### MYNA

Hey you,

Some considerations I would like to do before you actually go through this 
code. When I develop prototypes like this one I do like to use an unorthodox
approach so I am always learning new ways of achieving miserable failures,
please keep that in mind and lets see where I failed(or not).

This project is also bundled into a docker image and kept on docker hub, the
easy way of running it would be through this image and you can find a procedure
to do it right down here on this file.

I have tried to use the minimal amount of third party packages as possible, the
exceptions are the json schema validation, the message bird api client and the
activemq package iirc.

By default everything is kept in memory and once the process dies we lost all
pending sms messages. I have implemented also ActiveMQ support to avoid this
problem but it is not enabled by default, I rather prefer to keep things
simple.


#### Architecture

The following graph shows what the architecture looks like. The `queue` box
may be a memory buffer(the default) or an active mq connection.


```
       +          ^
       |          |
       |          |
 _V_   |          |
 @.@   v          +
(\_/)
 m-m----------------------+          +------------------------+
 |                        |          |                        |
 |        webserver       |    +---> |         queue          |
 |                        |    |     |                        |
 +-----+------------------+    |     +-----------+------------+
       |                       |                 |
       |          ^            |                 | +----<---+
       v          |            |                 | | writer |
                  |            |                 | +----^---+
 +----------------+-------+    |                 |
 |                        |    |                 |
 |     jsonvalidation     |    |                 |
 |                        |    |                 |
 +-----+------------------+    |                 |
       |                       |      +----<---+ |
       |          ^            |      | reader | |
       v          |            |      +---->---+ v
                  |            |
 +----------------+-------+    |     +------------------------+
 |                        |    |     |                        |
 |         store          +----+     |         sender         |
 |                        |          |                        |
 +------------------------+          +------------------------+
```


#### Running in docker

This is the easy way as you just need a message bird api key.

```
$ docker pull ricardomaraschini/myna
$ docker run -p8080:8080 --env MBKEY=<api-access-key> ricardomaraschini/myna
```

#### POST body

To send an sms please issue a POST request to /message with the body content
described by the json schema below

```
{
	"title": "sms",
	"type": "object",
	"properties": {
		"recipient": {
			"type": "string",
			"pattern": "^\\+[1-9]{1}[0-9]{3,14}$"
		},
		"originator": {
			"type": "string",
			"minLength": 1
		},
		"message": {
			"type": "string",
			"minLength": 1,
			"maxLength": 39015
		}
	},
	"additionalProperties": false,
	"required": ["recipient", "originator", "message"]
}
```

#### Compile locally

To compile you need `go version go1.9 linux/amd64`. Inside the project
directory run:

```
$ make
```

Should be enough

#### Running locally

Inside the project directory(MBKEY is the message bird api key):

```
$ MBKEY=<api-access-key> ./myna
```

#### Create docker image

Inside the project directory:

```
$ make docker
```
