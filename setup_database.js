//1. use database
use waku1;

//2. create users collection
db.createCollection("users");

//3. create rooms collection
db.createCollection("rooms");

//4. add index to users collection
db.users.createIndex({
    "phone": 1
}, {
    unique: true
});