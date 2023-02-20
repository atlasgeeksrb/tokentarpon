# TokenTarpon
## About
TokenTarpon is an API providing tokenization with optional encryption. Calling applications can use the service to avoid storing sensitive information such as PII in their own datastores. For example, post a string and receive a UUID, and then store the UUID rather than storing the sensitive text value. When the text value is needed, post the UUID and receive the text value. Callers can also send an array of values, and retrieve an array of values.

## TODO
- wrap up Docker files; needs SSL, automated testing, mongodb configuration and config settings
- call/use updateCheckSum
- finish writing tests
- allow user to send encryption options
- implement audit log
- provide API documentation
- build a simple demo front end
- implement JWT

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
