module github.com/Kisotu/daraja-client-go/cmd/daraja-cli

go 1.25.9

require github.com/Kisotu/daraja-client-go v0.0.0

require (
	github.com/google/uuid v1.6.0 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
)

replace github.com/Kisotu/daraja-client-go => ../..
