
create table if not exists product(
    Id int auto_increment primary key,
    UserID int not null,
    Name varchar(255) not null default "",
    Description varchar(255) not null default "default description",
    Price float not null default 0,
    Stock int not null default 0,
    Created_at timestamp default current_timestamp,
    Last_updated timestamp default current_timestamp,
    foreign key (UserID) references users(id)
)engine=innodb;