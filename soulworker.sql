CREATE DATABASE soulworker;

USE soulworker;

CREATE TABLE users (id INTEGER PRIMARY KEY AUTO_INCREMENT, username VARCHAR(10), password VARCHAR(64));

# austin:coolman83, placeholder info
INSERT INTO users (username, password) VALUES ('austin', '0afa580693c931eb17b465b761d35fc1eeb00aa70125256d1cc41319033337da');

create table characters
(
    id int primary key auto_increment,
    accountId int null,
    `index` int null,
    name varchar(20) null,
    class int null,
    appearance bigint null,
    level int null,
    constraint characters_users_id_fk
        foreign key (accountId) references users (id)
);

