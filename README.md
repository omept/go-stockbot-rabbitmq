# Go Stock Bot Microservice with RabbitMQ Message Broker

This is a decoupled service that calls an API using a "stock_code" as a parameter (<a href="​https://stooq.com/q/l/?s=aapl.us&f=sd2t2ohlcv&h&e=csv">​https://stooq.com/q/l/?s=aapl.us&f=sd2t2ohlcv&h&e=csv</a>​, here ​aapl.us is the stock_code​) and recieves a csv file. The  bot parses the received CSV file and then sends a message to the chatroom service using a message broker (RabbitMQ). 

## Requirement
The following tools are required to run the project
<ul>
<li> Go (Golang) </li> 
<li>RabbitMQ</li>
<li> Docker(optional)</li>
</ul>

<br />

## Context
This service reads from requests from env.STKBT_RECEIVER_QUEUE queue on RabbitMQ and publishes the processed stock requests on env.STKBT_PUBLISHER_QUEUE queue on RabbitMQ.

<br />

## <u>Starting The App</u>

<br />

## Step 1 : Start RabbitMQ
RabbitMQ is the message broker between the bot service and the chat service. To run it with docker, run the folling command in your terminal: <br />
### `docker run -d --hostname rabbitmq-svc --name rbbtmq -p 15672:15672 -p 5672:5672 rabbitmq:3.11.3-management`
<br />

## Step 2 : Update Env file
 copy the `.env.example` to `.env` and update the entries. 
<br />

## Step 3 : Starting the app
### run the folling command in your terminal `go run main.go`
<br />
