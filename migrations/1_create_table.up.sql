create table "users"(
    "id" serial not null primary key,
    "first_name" varchar(50) not null,
    "last_name" varchar(50) not null,
    "username" varchar(50) unique,
    "phone_number" varchar(30) unique,
    "email" varchar(100) not null unique,
    "password" varchar not null,
    "image_url" varchar,
    "type" VARCHAR(20) CHECK ("type" IN('superadmin', 'user')) NOT NULL,
    "created_at" timestamp default current_timestamp,
    "updated_at" timestamp,
    "deleted_at" timestamp 
);

INSERT INTO users(first_name, last_name, email, password, type) 
VALUES('zohid', 'saidov', 'zohidsaidov17@gmail.com', '$2a$10$nbFlsRUTs.8//V3s5TAczuETJClXKWx9zM95Kami9n9V1ODL/tQti', 'superadmin');

