#### MYNA

Hey you,

Some considerations I would like to do before you actually go through this 
code. When I develop prototypes like this one I do like to use an unorthodox
approach so I am always learning new ways of achieve miserable failure,
please keep that in mind and lets see where I failed(or not). :-)

This project is also bundled into a docker image and kept on docker hub, the
easiest way of running it would be through this image.

I have tried to use the minimal amount of third party packages as possible, the
exceptions are json schema validation, message bird api client and activemq
package(iirc).

By default everything is kept in memory and once the process dies we lost all
pending sms messages. I have implemented also ActiveMQ support to avoid this
problem, it is not enabled by default though.


#### Architecture

The following graph shows what the architecture looks like. The `queue` box
may be a memory buffer(the default) or an ActiveMQ connection.


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

This is the easiest way as you only need a Message Bird api key.

```
$ docker pull ricardomaraschini/myna
$ docker run -p8080:8080 --env MBKEY=<api-access-key> ricardomaraschini/myna
```

#### POST body

To send a sms please issue a POST request to /message with the body content
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

#### Running locally

Inside the project directory(MBKEY is the Message Bird API key):

```
$ MBKEY=<api-access-key> ./myna
```

#### Create docker image

Inside the project directory:

```
$ make docker
```
