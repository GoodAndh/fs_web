
create table if not exists Users(
id int not null auto_increment primary key,
username varchar(255) not null unique,  
email varchar(255) not null unique,
password varchar(255) not null,
created_at timestamp default current_timestamp,
last_updated timestamp default current_timestamp
)engine=innodb;