create table if not exists ulasan_room_chat_product_message(
    ID int not null primary key auto_increment,
    roomChatID int not null,
    message text,
    sendAt datetime not null default current_timestamp,
    isDeleted boolean default false,
    foreign key (roomChatID) references ulasan_room_chat_product(id)
)engine=innodb;