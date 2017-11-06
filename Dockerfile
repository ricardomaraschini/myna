FROM	golang:1.9
WORKDIR	/go/src/github.com/ricardomaraschini/myna
COPY	. .
COPY	assets/* /assets/
RUN	go install github.com/ricardomaraschini/myna
RUN	mv /go/bin/myna /
CMD	[ "/myna" ]
