#!/bin/bash

# Generate the Ed25519 key pair
private_key=$(openssl genpkey -algorithm ed25519)
public_key=$(echo "$private_key" | openssl pkey -pubout)

# Save the private key to file
echo "$private_key" > ./config/keys/ed25519_private.pem

# Save the public key to file
echo "$public_key" > ./config/keys/ed25519_public.pem
