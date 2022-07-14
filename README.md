# HADecrypt

HADecrypt is a CLI tool to encrypt/decrypt CMS On-Ramp data. It can also be
operated as a GRPC server for decrypting data through RPC.

## Installation

You may install directly through `go install` if you have go installed.

```bash
go install github.com/hkhait/crypto/cmd/hadecrypt@latest
```

Alternatively, clone the repository and compile from source.

```bash
go build cmd/hadecrypt/main.go
```

## Usage

### Encyption/Decryption

By default, the ciphertext/plaintext are read from standard input. You may store
all input in a plain text file (one ciphertext/plaintext per line) and pipe into
the program. Password is either specified through CLI flag `--password` or
environmental variable `DECRYPT_PASSWORD`.

```bash
hadecrypt decrypt --password=mysupersecretpw123 < ciphertext.txt > plaintext.txt
```

### GRPC Server

Alternatively, you may run the decryption tool as a service for external use. It
expose a single GRPC method to be called remotely.

```protobuf
service CryptoService {
    rpc Decrypt (DecryptRequest) returns (DecryptResponse) {}
}

message DecryptRequest {
    string ciphertext = 1;
}

message DecryptResponse {
    string plaintext = 1;
}
```

The service listen on port 8080 be default. You may configure the listen port
through CLI argument `--port`.
