default: build plan

deps:
	go install github.com/hashicorp/terraform

build:
	go build -o terraform-provider-fakewebservices .

test:
	go test -v .

plan:
	@terraform plan