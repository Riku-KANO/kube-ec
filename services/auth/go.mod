module github.com/Riku-KANO/kube-ec/services/auth

go 1.25

require (
	github.com/Riku-KANO/kube-ec/pkg v0.0.0
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
)

replace github.com/Riku-KANO/kube-ec/pkg => ../../pkg
