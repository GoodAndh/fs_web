

create table if not exists orders_items(
    ID int not null auto_increment primary key,
    OrderID int not null,
    UserID int not null,
    ProductID int not null,
    Total int not null,
    Price float not null default 0,
    foreign key (UserID) references users(id),
    foreign key (ProductID) references product(id),
    foreign key (OrderID) references orders_status(id)
)engine=innodb;