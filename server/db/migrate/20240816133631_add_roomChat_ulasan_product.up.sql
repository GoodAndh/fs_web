create table if not exists ulasan_room_chat_product(
ID int not null  primary key auto_increment,
roomID varchar(25) not null unique,
userID int not null,
ProductID int not null,
username varchar(255) not null,
foreign key (userID) references users(id),
foreign key (ProductID) references product(id)
)engine=innodb;