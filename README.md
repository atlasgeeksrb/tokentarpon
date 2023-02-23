# TokenTarpon
## About
TokenTarpon is a service providing tokenization with optional encryption. Calling applications can use the service to avoid storing sensitive information such as PII in their own datastores. 

Post a string and receive a UUID, and then store the UUID rather than storing the sensitive text value. When the text value is needed, post the UUID and receive the text value. Callers can also send an array of values, and retrieve an array of values.

## How to Use
- use Git to pull latest
- update settings
  - modify docker-compose.yml
    - set a strong password for the root mongodb user
    - if desired, change the host's port from 8092
  - modify config.json
    - set the EncryptionKey, used when encrypting stored values
    - set the mongodb settings (myuser, mypassword, mydb)
  - modify mongoinit.js
    - change myuser, mypassword, mydb to match the settings in config.json
- start containers
  `docker compose up`

## TODO
- duplicate tokens are generated
- provide Postman tests, Swagger docs
- wrap up Docker files; needs SSL, automated testing, mongo config settings
- call/use updateCheckSum
- accept a flag to indicate whether to encrypt stored values
- finish writing tests
- allow user to send encryption options
- implement audit log with checksums
- provide API documentation
- build a simple demo front end
- implement JWT
- provide account-domain mapping, with flooded, blocked
- provide admin endpoints
  - page and search logs
  - change account settings & mappings
  - endpoint for changing the encryption key

## Routes
- PUT a token /tokens/:domainId
- GET a token /tokens/:domainId/:id
- GET token value /tokens/:domainId/:id/value
- DELETE the specified token /tokens/:domainId/:id
- PUT several tokens /tokens/:domainId/:id
- GET tokens for domain tokens/tokens/:domainId
- POST a query to get multiple token values /tokens/:domainId/values

The last two routes, to get tokens and get token values, optionally take start and limit parameters in the querystring for pagination.

## Modifying
Please refer to the LICENSE
