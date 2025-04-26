create table user
(
    id         bigint primary key auto_increment,
    email      varchar(255),
    username   varchar(64),
    password   varchar(255),
    created_at datetime default current_timestamp,
    updated_at datetime default current_timestamp on update current_timestamp,
    deleted_at datetime,
    UNIQUE KEY `username` (`username`)
);

create table token_record
(
    id         bigint primary key auto_increment,
    user_id     bigint       not null,
    token      varchar(255) not null,
    created_at datetime default current_timestamp,
    deleted_at datetime,
    key (user_id),
    UNIQUE KEY `token` (`token`)
)