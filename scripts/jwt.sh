#!/bin/bash

# Generate a random string with high entropy
secret=$(openssl rand -base64 32)

# Print the secret
echo "Generated secret: $secret"
#!/bin/bash