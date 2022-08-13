create table if not exists tb_access(
    token varchar(255) not null,
    user_uid varchar(36) not null,
    expired_at timestamp not null,
    is_blocked boolean not null DEFAULT FALSE,
    created_at timestamp not null DEFAULT Now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_access primary key (token),
    constraint fk_access_0 foreign key (user_uid) references tb_user(id)
)