### HUNT Prototype ###
HUNT (Heads-up node tracker) is a system designed to complement command center real-time awareness
visualization tools. HUNT users use AR glasses to view updates to locations they're tracking and can select / de-select targets. HUNT users also present their location to the command center via gps updates.

### Setup ###
1. To run:
    a. cd into /hunt
    b. Run the run.sh script
    c. To view the map: 
        i. Open localhost:8080
        ii. In the top app bar, click the 'Connect' button
    d. To toggle the use of in-memory app state vs. mongoDB, change the useState variable
        in /hunt/socket/handler to your chosen setting
2. To test:
    a. Open a new terminal
    b. cd into /hunt_test
    c. Run the app (go run .)
    d. To view metrics:
        i. Open localhost:9090
    e. To increase the number of clients:
        i. Alter the 'ClientCount' constant in /hunt_test/main.go to a number below 256
        ii. To increase the maximum number of clients, go to /hunt/state/state.go and alter the 
            constant value "maxClients"
        iii. To increase the number of location updates sent per second, alter the 'UpdatePeriod' 
                constant in /hunt_test/main.go


### Components ###
1. /socket contains the websocket handler, handler functions, and events
2. /collections contains structs and methods that interact with mongoDB
3. /db contains connection logic and database creation for mongoDB
4. /logging contains a zerolog logger
5. /constants contains app-wide constants
6. /static contains files to serve a web application

### Messaging ###
The application uses JSON to send/receive messages

### State ###
1. Global
    a. There are a handful of exported state variables / constants that can be used throughout the codebase
    b. Collections are global and exported from their package for use in various areas of the codebase
    c. State objects are unexported and kept within their package
2. In-memory vs DB
    a. There is an experimental in-memory data store that uses a data-oriented design style of 'struct of arrays' 
    b. The mongoDB database can also be used to store state such as users and their locations, command nodes, and units
3. Socket connections
    a. There is a connMap variable in /hunt/socket/handler.go that stores each connection's uid and socket information

### TODO ###
1. Use protobuf over json
2. Encapsulate state variables within socket manager
3. Implement caching when using DB for things that update less frequently (command nodes, users, units)