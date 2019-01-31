CREATE DATABASE authdb;

\c authdb;

CREATE TABLE IF NOT EXISTS  access_control
(
  user_type VARCHAR(20) PRIMARY KEY,
  access_mode INTEGER NOT NULL CHECK (access_mode >= 0 AND access_mode <= 7)
);

INSERT INTO access_control VALUES ('general',4);
INSERT INTO access_control VALUES ('admin',7);
INSERT INTO access_control VALUES ('moderator',6);
INSERT INTO access_control VALUES ('manager',5);

CREATE TABLE IF NOT EXISTS  users
(
  username VARCHAR(30) PRIMARY KEY,
  first_name VARCHAR(30),
  last_name VARCHAR(30),	
  deleted BOOLEAN DEFAULT FALSE NOT NULL,
  user_type VARCHAR(20) DEFAULT 'general' REFERENCES access_control(user_type),
  e_mail VARCHAR(40) UNIQUE NOT NULL,
  password VARCHAR(70) NOT NULL
);


INSERT INTO users VALUES ('alex','Nikita','Krilov', FALSE,'admin','alex@gmail.com','4d31262e00d95e5ceb7477098a63e06fb3333cc26c0950e47f2323a791e45622'); --qwerty
INSERT INTO users VALUES ('theeska','Nikolay', 'Gladkov',FALSE, 'moderator','theeska@yandex.ru','74abcd46f765d1f49252f688541bc12a4dacf0c1422a52e77d05b0e99c27362b'); --asdfgh
INSERT INTO users VALUES ('vasya2515','Vitya','Pupkin',FALSE,'manager','kuku@mail.com','877e24b6826f0d583c2c1d9ad8873455dde43bed652edbe5c4e5d5c969dea479'); --zxcvbn
INSERT INTO users VALUES ('svetik5245','Svetlana', 'Korobkova', FALSE,'general','pistolet@mail.ru','f9accf1c65c276ea94a232fc05124b96e706212cd975530c02f28c75a7cafdf2'); --qazwsx
INSERT INTO users VALUES ('yuibn','Ekaterina','Volkova',TRUE, 'general','qweasd@gmail.com','7f01fbd5c436b72e0d048ba8fb1ff1f43fdf5131703287602af5d213b4f2a76c'); --edcrfv


