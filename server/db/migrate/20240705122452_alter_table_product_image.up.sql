alter table product_image
add column captions varchar(255) not null default "empty captions" after url;