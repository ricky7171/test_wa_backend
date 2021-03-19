
  

# WAku (Chatting app)

  

Hi sir, here I will show my app chat project.

  

## MongoDB Database Schema

![alt text](https://raw.githubusercontent.com/ricky7171/test_wa_backend/master/requirement/database/schema.png)

  

## List API

1. Login API

> url : localhost:8080/api/auth/login (POST) <br>
> body : {"phone" : "...","password" : "..."}
  

2. Register API

> url : localhost:8080/api/auth/register (POST) <br>
> body : {"phone" : "...","name" : "...","password" : "..."}
  

3. Get Contact API

> url : localhost:8080/api/contact (GET) <br>
> header : {Authorization : "Bearer ..."}
  

4. Get Chat

>url : localhost:8080/api/chat/:contactId/:lastId (GET) <br>
> header : {Authorization : "Bearer ..."} <br>
> note : This API use pagination. When load first page, set lastId to "nil". For another page, set lastId according to last id chat that get before.

  

5. New Chat

 > url : localhost:8080/api/new_chat (POST) <br>
>header : {Authorization : "Bearer ..."} <br>
>body : {"phone": "...", "message": "..."} 


6. Websocket


>url : localhost:8080/ws/:user_id <br>
>header : {Authorization : “Bearer …”} 

  
  

## Features

  

- Register, Login, & Logout

- Add other user to their contacts

- Realtime chat with other user

- See chat history

  

## Library Used

- Gin for routing (github.com/gin-gonic/gin)

- JWT Auth for authentication using token (github.com/dgrijalva/jwt-go)

- Validator for validate request (github.com/go-playground/validator/v10)

- Gorilla Websocket for communicate using websocket (github.com/gorilla/websocket)

- Godotenv for environment setting (github.com/joho/godotenv)

- Mongo Driver for connect to mongoDB (go.mongodb.org/mongo-driver)

  

## Screenshot

- Register page

![alt text](https://github.com/ricky7171/test_wa_backend/blob/master/requirement/screenshot/register.png?raw=true)

- Login page

![alt text](https://github.com/ricky7171/test_wa_backend/blob/master/requirement/screenshot/login.png?raw=true)

- Home page

![alt text](https://github.com/ricky7171/test_wa_backend/blob/master/requirement/screenshot/home.png?raw=true)

- Send new message

![alt text](https://github.com/ricky7171/test_wa_backend/blob/master/requirement/screenshot/send%20new%20message.png?raw=true)

- Incoming message

![alt text](https://github.com/ricky7171/test_wa_backend/blob/master/requirement/screenshot/first%20incoming%20message.png?raw=true)

- Continue chatting

![alt text](https://github.com/ricky7171/test_wa_backend/blob/master/requirement/screenshot/continue%20chatting.png?raw=true)

## How to Run App Locally

1. Make sure you have installed golang (https://golang.org/doc/install) & mongoDB (https://docs.mongodb.com/manual/installation/)

2. Clone this repository to your PC

3. Go to project directory and open setup_database.js file

4. Copy all script in that file

5. Now, open MongoDB terminal

&nbsp; &nbsp; &nbsp; &nbsp; 5.a. Open terminal (Because I use windows, I open command prompt) <br>
&nbsp; &nbsp; &nbsp; &nbsp; 5.b. Go to mongo directory (my mongo dir is : "C:\Program Files\MongoDB\Server\4.4\bin" <br>
&nbsp; &nbsp; &nbsp; &nbsp; 5.c. Then type "mongo"

![alt text](https://github.com/ricky7171/test_wa_backend/blob/master/requirement/screenshot/open%20mongo.png?raw=true)

6. Run all script to that mongo terminal (just paste it)

![alt text](https://github.com/ricky7171/test_wa_backend/blob/master/requirement/screenshot/run%20setup.png?raw=true)

7. Wait until done.

8. Open new terminal and go to project directory

9. Run the following command :

>go build <br>
>go run .

10. Open 2 different browser and go to localhost:8080/register to test it

11. Register with name, phone, and password

12. After register, login using that account

13. Try to send new message

![alt text](https://github.com/ricky7171/test_wa_backend/blob/master/requirement/screenshot/send%20new%20message.png?raw=true)

  

## How to Try API on Postman

1. Open project directory

2. Download and install postman

3. Open postman

4. Create new workspace

5. Import postman file from PROJECT_DIRECTORY/requirement/postman/API.postman_collection.json

  

note : postman is only used to run restful API (because until now postman still doesn't support websocket)

  

## Performance Optimization

I have do some optimization for server & database performance :
<ol>
  <li>
Create a fast performing database structure in MongoDB <br>
 I've tried several possible database structures, and the last structure I tried was fast enough. I use 3 collection : users, contacts, and chats. It faster than I just use 2 collection : users, contacts (Chat data is in the contacts collection as array)  
  </li>
  <li>
Add indexes in certain collections as needed <br>
 I have added indexes to 3 collections in certain fields so that the reading process is faster.
  </li>
  <li>
Using go routine so that the websocket process can run simultaneously <br>
 I use go routine to listen and write data to client. I also use channel as a "bridge" for websocket data communications, like example when user connect to websocket, send message, retreive message.
  </li>
  <li>
Using mutex to prevent deadlock <br>
I use mutex when process writing to websocket, because it can prevent deadlock.
  </li>
  <li>
Make queries as efficient as possible <br>
I test every query on mongo terminal and measure that time execution (using profile mode and explain() function). I also try to make 100.000 user with 5.000 contact for each user and 10.000 chat for 300 first contact. I try all query that needed in backend, and it still fine (after I adding index). I also follow some suggestion about optimization query from mongo documentation (ex : https://docs.mongodb.com/manual/reference/operator/aggregation/match/#pipeline-optimization)
  </li>
</ol>

  
  

For the future, if this application develops with more complex business processes and our servers are too busy, maybe we can do several things such as adding new servers, implementing microservices, sharding in Mongodb, etc.