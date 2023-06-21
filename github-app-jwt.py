#!/usr/bin/env python3

# From https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/generating-a-json-web-token-jwt-for-a-github-app#example-using-python-to-generate-a-jwt

import jwt
import time
import os

# Get the App ID
app_id = os.environ['GH_APP_ID']

# Open PEM
signing_key = jwt.jwk_from_pem(os.environ['GH_APP_PRIVATE_KEY'].encode('utf-8'))

payload = {
    # Issued at time
    'iat': int(time.time()),
    # JWT expiration time (10 minutes maximum)
    'exp': int(time.time()) + 600,
    # GitHub App's identifier
    'iss': app_id
}

# Create JWT
jwt_instance = jwt.JWT()
encoded_jwt = jwt_instance.encode(payload, signing_key, alg='RS256')

print('jwt=' + encoded_jwt)
