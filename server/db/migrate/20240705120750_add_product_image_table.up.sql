create table if not exists product_image(
ID int not null auto_increment primary key,
ProductID int not null,
Url varchar(255) not null,
foreign key (ProductID) references product(id)
)engine=innodb;