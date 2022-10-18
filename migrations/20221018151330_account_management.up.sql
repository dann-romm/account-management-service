create table users
(
    id         serial primary key,
    username   varchar(255) not null,
    password   varchar(255) not null,
    created_at timestamp    not null default now()
);

create table accounts
(
    id         serial primary key,
    balance    int       not null,
    created_at timestamp not null default now()
);

create table operation_types
(
    id   serial primary key,
    name varchar(255) not null
);

create table operations
(
    id                serial primary key,
    user_id           int       not null,
    account_id        int       not null,
    amount            int       not null,
    operation_type_id int       not null,
    description       text               default null,
    created_at        timestamp not null default now(),
    foreign key (user_id) references users (id),
    foreign key (account_id) references accounts (id),
    foreign key (operation_type_id) references operation_types (id)
);
