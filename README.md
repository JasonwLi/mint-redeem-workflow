# Mint and Redeem Workflow Engine
This engine is a simple workflow engine that is able to take post requests to the `/mint` and `/redeem` endpoints and execute them against a mocked brale api. The execution is done through the usage of Cadence. An open source workflow engine that is developed by Uber. 

Cadence workflows give our workflow engine automatice state persistence which can protect us from failovers and provides an unbroken trail for the request. It also has builting idempotency handling allowing us to prevent double spends in our flow. In this engine I use the unique request IDs as the workflow and the idem ids for cadence and brale. This means each request will only automatically spin up one workflow and that workflow will only ever be able to make that request once to brale, preventing double spends.

Cadence also provides excellet horizontally scaling capabilities as the workflows are lightweight and can be put to sleep during long processing actions freeing up computing resources. 

The built in Cadence UI gives us good debugging and event logging helping engineers to manually solve issues as the last solution. 

The use of cadence makes our workflow engine easily extendable since the only work required is to code new workflow logic and register them with the worker. 

The downside of Cadence is that the startup and infra costs are quite high. The system is complicate to get started with without prior experience requiring extensive domain knowledge. It also forces the operator to use the specified Cadence stack such as having to maintain a Cassandra db on all of the environments. Its also not very suitable for tasks that require quick turn arounds, i.e a Cadence workflow wouldnt be used for creating an user that doesn't require any external calls. 

The service db is a simple sqllite db to mimic a postgres rds in production. I chose a traditional RDS because the records that represent the requests are extremely light weight and we have no other relations with any other records, this rules out document dbs as we do a lot more writes than reads. 

The workflow db is cassandra which is chosen by Cadence by default.

The code is modularized into different big chuncks. The api layer, the service layer, the workflow layer, the infra layer, the model/db layer and the deps/config layer. This allows us to add to each section independently without having to worry too much about breaking changes affecting other areas of the codebase. 

## What could be improved if given more time
1. The codebase lacks validation in many places. If given more time all params for every method would be a strongly typed struct represented by a `valueobject` That is initializable with simple types supported by golang but performs validations on the value. i.e using common.Address for validating evm addresses. This would help with validating API requests as well.
2. Using better mocks. I would use dependency ejection a bit more efficiently when it comes to my api handlers. This would allow me to directly inject mocked calls into the tests rather than having to define methods to be mocked as package level variables to be overriden by tests. This would allow me to directly unit test my activities code better. 
3. The code base is not as organised as i would like. There are many shared configs being duplicated(mainly relating to workflow setup) I would define a separate workflow config package to manage these. 
4. The brale client itself is somewhat lacking as it actually lacks the capabilities to make REST calls as everything is mocked.
5. I would also introduce more logging in the code to help for debugging purposes in a production environment.
6. The main method responsible for starting both the worker and API should be separated so if one panics it doesnt kill the other process. I chose to them together to combine the different DBs sqllite creates, otherwise workflows wouldnt be able to access the same records as the api. Of course this wouldnt be a problem with proper infra setup.

### How to run this repo
0. Clone the repo onto your machine
1. ensure you have go installed on your machine. At least > 1.19. You can install go using following the instructions here: https://go.dev/doc/install
2. ensure you have docker installed on your machine. We need docker to run our Cadence workflows for persistence.
https://docs.docker.com/engine/install/
3. In the root repo run `docker-compose up` you should see 6 containers spinning up under the namespace `mint-redeem-workflow` Wait for ~2 before continuing to the next step as the startup can take awhile
4. We now need to register a Cadence domain. If you are running a m1 machine or later use this command: `docker run --platform linux/amd64 --network=host --rm ubercadence/cli:master --do test-domain2 domain register -rd 1
` if you are on a pre m1 machine just run `docker run --network=host --rm ubercadence/cli:master --do test-domain2 domain register -rd 1`
5. Check that your domain is registered correctly `docker run --network=host --rm ubercadence/cli:master --do test-domain2 domain describe`
6. Now we are ready to spin up our api and workers. In the root of this repo, run `go run main.go` this will spin up the gin api on `localhost:8090` and the workers on `localhost:8080` you will also be able to access the cadence ui for managing workflows on http://localhost:8088/
7. Once this is ready you are welcome to make curl requests to the api. I've provided a couple of samples below
```
curl -X POST http://localhost:8090/mint \
-H "Content-Type: application/json" \
-d '{
    "amount": 100.50,
    "recipient": "0xtestreceive"
}'

curl -X POST http://localhost:8090/redeem \
-H "Content-Type: application/json" \
-d '{
    "amount": 50.75,
    "recipient": "0xtestreceive"
}'
```
8. To see the api error with the workflow you can curl with a specific address: `0xdeadbeef`
```
curl -X POST http://localhost:8090/redeem \
-H "Content-Type: application/json" \
-d '{
    "amount": 50.75,
    "recipient": "0xdeadbeef"
}'
```
9. After submitting the curls you can visit http://localhost:8088/domains/test-domain2/workflows?range=last-30-days to check the status of the workflows. 

### Tests
Tests can be run by cding into each dir and running `go test`