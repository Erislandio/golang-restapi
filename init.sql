CREATE DATABASE IF NOT EXISTS `account`;

use `account`;

CREATE TABLE IF NOT EXISTS users (
    id int not null auto_increment,
    name varchar(100) not null,
    email varchar(100) not null unique,
    phone varchar(15),
    primary key (id)
);

INSERT INTO 
`users` (name, email, phone) 
values 
('Erislandio', 'erislandiosoares21@gmail.com', '11942676399'),
('maria', 'mariasoares21@gmail.com', '352362457'),
('edu', 'edusoares21@gmail.com', '1345234523464'),
('jessica', 'jessicasoares21@gmail.com', '2224123135');