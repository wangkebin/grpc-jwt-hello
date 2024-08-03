#gRPC Hello World with JWT auth
to recompile protobuf file, run the following command within folder helloworld
```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    helloworld.proto
```

to test
```
~/go/src/grpc-jwt-hello ➤ go run greater_server/main.go myjwt/ed25519-public.pem 
~/go/src/grpc-jwt-hello ➤ go run greater_client/main.go ./myjwt/ed25519-private.pem
```