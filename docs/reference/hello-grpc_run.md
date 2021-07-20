## hello-grpc run

Launch Hello GRPC server

### Synopsis

Launch Hello GRPC server

```
hello-grpc run [flags]
```

### Options

```
      --api-domain string             Domain used for apiserver (prod: api.appscode.com
      --cors-origin-allow-subdomain   Allow CORS request from subdomains of origin (default true)
      --cors-origin-host string       Allowed CORS origin host e.g, domain[:port] (default "*")
      --enable-cors                   Enable CORS support
  -h, --help                          help for run
      --log-rpc                       log RPC request and response
      --plaintext-addr string         host:port used to serve http json apis (default ":8080")
      --secure-addr string            host:port used to serve secure apis (default ":8443")
      --tls-ca-file string            File containing CA certificate
      --tls-cert-file string          File container server TLS certificate
      --tls-private-key-file string   File containing server TLS private key
```

### SEE ALSO

* [hello-grpc](hello-grpc.md)	 - Hello gRPC by Appscode

