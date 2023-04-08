module github.com/gusaul/grpcox

go 1.12

require (
	github.com/fullstorydev/grpcurl v1.3.2
	github.com/gorilla/mux v1.7.0
	github.com/jhump/protoreflect v1.5.0
	go.uber.org/zap v1.24.0
	google.golang.org/grpc v1.21.0
)

replace go.uber.org/zap => github.com/santakdalai90/zap v0.0.0-20230408061401-268eb0bfc7e2
