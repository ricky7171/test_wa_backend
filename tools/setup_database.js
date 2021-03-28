//1. use database
use waku;

//2. create users collection
db.createCollection("users");

//3. create contacts collection
db.createCollection("contacts");

//4. create chats collection
db.createCollection("chats");

//5. add index to users collection
db.users.createIndex({
    "phone": 1
}, {
    unique: true
});

//6. add index to contacts collection
db.contacts.createIndex({
    "users": 1
});

//7. add index to contact_id on chats collection
db.chats.createIndex({
    "contact_id": 1
});