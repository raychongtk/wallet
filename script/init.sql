create table app_user
(
    id           uuid primary key,
    email        varchar(255) not null,
    phone_number varchar(100) not null,
    password     varchar(255) not null,
    first_name   varchar(36)  not null,
    last_name    varchar(36)  not null,
    created_at   timestamp default CURRENT_TIMESTAMP,
    updated_at   timestamp
);

create
    unique index app_user_email_uindex
    on app_user (email);

create
    unique index app_user_phone_number_uindex
    on app_user (phone_number);

insert into app_user
    (id, email, phone_number, password, first_name, last_name)
values ('2d988f4a-a037-4ce9-a350-f13445793e88', 'test@gmail.com', '+85212345678',
        '40bd001563085fc35165329ea1ff5c5ecbdbbeef', 'John', 'Doe');
insert
into app_user
    (id, email, phone_number, password, first_name, last_name)
values ('c6e97817-0254-43ad-8610-7ac9d3f7af92', 'test2@gmail.com', '+85212345679',
        '7c4a8d09ca3762af61e59520943dc26494f8941b', 'Ray', 'Doe');

create table account
(
    id           uuid primary key,
    user_id      uuid null,
    account_type varchar(255) not null,
    created_at   timestamp default current_timestamp,
    updated_at   timestamp
);

insert into account
    (id, user_id, account_type)
values ('81407970-675d-4c59-b49f-01beec7e9280',
        '2d988f4a-a037-4ce9-a350-f13445793e88', 'CUSTOMER');

insert into account
    (id, user_id, account_type)
values ('87100f76-87e8-49b5-96f8-e741365260a1',
        'c6e97817-0254-43ad-8610-7ac9d3f7af92', 'CUSTOMER');

insert into account
    (id, user_id, account_type)
values ('cf94e219-ed3c-4cb3-a13d-3d8d756192c8',
        'c19f00b4-c457-43f5-9e30-d10ada02a94f', 'CHART');

insert into account
    (id, user_id, account_type)
values ('105b261c-600d-4a60-8b99-63d60d05e82b',
        'd3b07384-d9a0-4f3b-8a2b-6c9e5b8b8f3c', 'CHART');

create table wallet
(
    id            uuid primary key,
    account_id    uuid        not null,
    currency      varchar(3)  not null,
    decimal_place int         not null,
    wallet_status varchar(30) not null,
    created_at    timestamp default current_timestamp,
    updated_at    timestamp
);

insert into wallet
    (id, account_id, currency, decimal_place, wallet_status)
values ('1cc535a5-bc57-4731-a64b-041b7ff41c30',
        '81407970-675d-4c59-b49f-01beec7e9280', 'USD', 2, 'ACTIVE');

insert into wallet
    (id, account_id, currency, decimal_place, wallet_status)
values ('c7d90b83-e080-423a-ab1b-f48094d7533e',
        '87100f76-87e8-49b5-96f8-e741365260a1', 'USD', 2, 'ACTIVE');

insert into wallet
    (id, account_id, currency, decimal_place, wallet_status)
values ('338b3f97-e428-4bff-9775-f759b5fccc4d',
        'cf94e219-ed3c-4cb3-a13d-3d8d756192c8', 'USD', 2, 'ACTIVE');

insert into wallet
    (id, account_id, currency, decimal_place, wallet_status)
values ('141e3fd8-c350-4b44-a2d5-2e2602aca72a',
        '105b261c-600d-4a60-8b99-63d60d05e82b', 'USD', 2, 'ACTIVE');

create table balance
(
    id           uuid primary key,
    wallet_id    uuid        not null,
    balance_type varchar(50) not null,
    balance      int         not null,
    created_at   timestamp default current_timestamp,
    updated_at   timestamp
);

insert into balance
    (id, wallet_id, balance_type, balance)
values ('e80c2fd1-e470-4c21-b88e-0459cc52164c',
        '1cc535a5-bc57-4731-a64b-041b7ff41c30', 'COMMITTED', 0);

insert into balance
    (id, wallet_id, balance_type, balance)
values ('2707b0fa-ce0d-4748-bb5d-98ce7160c83a',
        'c7d90b83-e080-423a-ab1b-f48094d7533e', 'COMMITTED', 0);

insert into balance
    (id, wallet_id, balance_type, balance)
values ('14b13349-e309-4d2d-8ad6-81355d8bfbbc',
        '338b3f97-e428-4bff-9775-f759b5fccc4d', 'COMMITTED', 0);

insert into balance
    (id, wallet_id, balance_type, balance)
values ('f69815d2-0c50-46a6-9aa9-c5ef1ebddc24',
        '141e3fd8-c350-4b44-a2d5-2e2602aca72a', 'COMMITTED', 0);

create table movement
(
    id               uuid primary key,
    debit_wallet_id  uuid        not null,
    credit_wallet_id uuid        not null,
    debit_balance    int         not null,
    credit_balance   int         not null,
    movement_status  varchar(50) not null,
    created_at       timestamp default current_timestamp,
    updated_at       timestamp
);

create table transaction
(
    id           uuid primary key,
    movement_id  uuid        not null,
    wallet_id    uuid        not null,
    balance_type varchar(50) not null,
    balance      int         not null,
    created_at   timestamp default current_timestamp,
    updated_at   timestamp
);