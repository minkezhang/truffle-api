Protobufs
----

Protobuf definitions for the Bene API package.

## Compile

Run from project root

```bash
protoc -I ./ \
  --go_out=proto/go/ \
  --go_opt=paths=import \
  --go_opt=module=github.com/minkezhang/truffle-api/proto/go \
  proto/*.proto
```
