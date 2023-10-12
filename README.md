## Diagram

![image](https://github.com/farismfirdaus/plant-nursery/assets/62270078/9e3d6915-2b74-4172-91c6-380bf948fbc6)

## Endpoint

https://www.postman.com/mission-explorer-62696199/workspace/plant-nursery/overview

## Getting Started

### Dependencies

Please install each of below to your machine:

- Go [installation instruction](https://go.dev/doc/install)
- Postgres [installation instruction](https://www.postgresql.org/download/)

### Initial Setup

- Clone repository
  ```
  $ git clone git@github.com:farismfirdaus/plant-nursery.git
  $ cd plant-nursery
  ```

- Setup env variables
  ```
  $ cp .env.sample .env
  ```
  
- Setup cert
  ```
  $ make cert
  ```

- Migrate database
  ```
  $ make migrate
  ```

- Seed database
  ```
  $ make seed
  ```
  
- Run application
  ```
  $ make run
  ```

## Todo

- [ ] Logging. need more configuration awful log
- [x] Unit test: customer service
- [x] Unit test: cart service
- [x] Unit test: plant service
- [x] Unit test: order service
- [ ] Unit test: all controller
- [ ] Unit test: all repository
- [ ] Unit test: utils
- [ ] Prettify response struct

## User Stories

- [x] As a customer I should be able to create a Garden Enthusiasts account;
- [x] As a customer I should be able to view the list of available plants;\
A customer should be able to purchase a plant;
- [x] As a customer I should be able to select one or multiple plants to add to their cart;\
They should be able to specify the quantity for each plant they choose;\
Once they've selected everything they want, they should be able to place an order;\
The system should track the total cost;
- [x] As a customer I should be able to view Purchase History:\
Users should be able to see a history of all the plants they've ordered;\
Each entry in the purchase history should show the date of purchase, list of plants bought, and the total amount paid;
