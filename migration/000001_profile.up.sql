create table if not exists profiles(
    id bigserial primary key,
    phone_number varchar(20) unique not null,
    balance float8 not null default 0
);

create index phone_number_idx on profiles(phone_number);