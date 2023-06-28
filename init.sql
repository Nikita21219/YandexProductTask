create table "order"
(
    id            serial primary key,
    courier_id    integer references courier,
    weight        integer not null,
    region        integer not null,
    delivery_time varchar(255),
    complete_time timestamp,
    price         integer not null
);

create table courier
(
    id            serial primary key,
    courier_type  varchar(50) not null,
    regions       integer[]   not null,
    working_hours text[]      not null
);
