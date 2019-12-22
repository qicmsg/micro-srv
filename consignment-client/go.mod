module consignment-client

go 1.13

replace consignment-service => ../consignment-service

require (
	consignment-service v0.0.0-00010101000000-000000000000
	github.com/micro/go-micro v1.18.0
	github.com/micro/micro v1.18.0 // indirect
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553
)
