#!/bin/bash
SERVER_CN=localhost

# Private files: ca.key, server.key, server.pem, server.crt
# Shared files: ca.crt (for client), server.csr, (for CA)

# Step 1: Generate Cert. Authority + Trust Cert. (ca.crt)
openssl genrsa -passout pass:1111 -des3 -out ca.key 4096
openssl req -passin pass:1111 -new -x509 -days 365 -key ca.key -out ca.crt -subj "/CN=${SERVER_CN}"

# Step 2: Generate the Server Private Key (server.key)
openssl genrsa -passout pass:1111 -des3 -out server.key 4096

# Step 3: Get a certificate signing request from the CA (server.csr)
openssl req -passin pass:1111 -new -key server.key -out server.csr -subj "/CN=${SERVER_CN}"

# Step 4: Sign the certificate with the CA we created (self-signing) - server.crt
openssl x509 -req -passin pass:1111 -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt

# Step 5: Convert the server certificate to .pem format for gRPC - server.pem
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out server.pem