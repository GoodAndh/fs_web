

create table if not exists orders_status(
    ID int not null auto_increment primary key,
    ProductID int not null,
    UserID int not null,
    Status enum("paid","wait","cancel") not null default "wait",
    TotalPrice float not null default 0,
    foreign key (UserID) references users(id),
    foreign key (ProductID) references product(id)
)engine=innodb;