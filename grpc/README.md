# gRPC

I build a client and a server and connect them using grpc and protocol buffers.

## Generate code from .protof file
```
protoc -I protos protos/currency.proto --go_out=protos/. --go-grpc_out=protos/.
```
## Use grpcurl to test server before making client link in below
[grpcurl](https://github.com/fullstorydev/grpcurl)

## Some Usefull commands
```
grpcurl -plaintext 127.0.0.1:9092 list
grpcurl -plaintext 127.0.0.1:9092 list
grpcurl -plaintext -d '{"Base": "RIAL", "Destination": "USD"}' 127.0.0.1:9092 Currency.GetRate
```
For the first two command you should enable reflection.

