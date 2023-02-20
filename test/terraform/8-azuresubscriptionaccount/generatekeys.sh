!#/bin/bash
openssl req -x509 -newkey rsa:4096 -sha256 -keyout my.key -out my.crt -subj "/CN=test.com" -days 600
openssl pkcs12 -export -name "test.com" -out my.pfx -inkey my.key -in my.crt