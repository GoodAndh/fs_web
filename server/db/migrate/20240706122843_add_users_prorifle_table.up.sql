create table if not exists users_profile(
ID int not null auto_increment primary key,
userID int not null,
Url varchar(255) not null default "url less",
Captions varchar(255) not null default "caption less",
foreign key (userID) references Users(id)
)engine=innodb;