

create table if not exists cart(
    Id int auto_increment primary key,
    UserID int not null,
    ProductID int not null, 
    /* Status enum () */
    Total int not null,
    Price float not null,
    ProductName varchar(255) not null,
    Created_at timestamp default current_timestamp,
    Last_updated timestamp default current_timestamp,
    foreign key (UserID) references users(id),
    foreign key (ProductID) references product(id)
)engine=innodb;