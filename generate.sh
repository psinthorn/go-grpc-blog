!#/bin/bash

# Generate Protofile
protoc auth.proto --go_out=plugins=grpc: .