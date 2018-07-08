GO=go
DST=dst

all: build

build: 
	@mkdir -p $(DST)
	$(GO) build main.go
	@mv main dst

get:
	go get github.com/aws/aws-sdk-go/aws
	go get github.com/aws/aws-sdk-go/aws/session
	go get github.com/aws/aws-sdk-go/service/ec2
	go get github.com/aws/aws-sdk-go/aws/ec2metadata

clean:
	$(RM) -rf dst
