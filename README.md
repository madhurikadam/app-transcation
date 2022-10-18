# app-transcation
- Microservice for handling the cusomter cards accounts and transcations
### Overview
- Each cardholder (customer) has an account with their data.
- For each operation done by the customer a transaction is created and associated with their
respective account.
- Each transaction has a specific type (normal purchase, withdrawal, credit voucher or
purchase with installments)
- Transactions of type purchase and withdrawal are registered with negative amounts, while
transactions of credit voucher are registered with positive value.

## Requirements
 - go sdk 1.18 
 - go mod 
 - docker 
 - docker-compose
 - postgres
 - make tool

## Local Setup 
### pre requisite
- ```docker, docker-compose, go, make```
### To up the dependency 
- This will create the postgres docker container
> `make up`

### To build the app binary and run binary
> `make build`

### To run the service
> `make dev`

## Authors
Madhuri Kadam
madhurikadam300@gmail.com