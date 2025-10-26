module github.com/Riku-KANO/kube-ec/services/auth

go 1.25

require (
	github.com/Riku-KANO/kube-ec/pkg v0.0.0
	github.com/Riku-KANO/kube-ec/proto v0.0.0
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	google.golang.org/grpc v1.62.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240123012728-ef4313101c80 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

replace github.com/Riku-KANO/kube-ec/proto => ../../proto

replace github.com/Riku-KANO/kube-ec/pkg => ../../pkg
