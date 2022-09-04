create table if not exists tb_code_recovery(
    code varchar(6) unique not null,
    user_uid varchar(36) not null,
    expired_at timestamp not null,
    created_at timestamp not null DEFAULT Now(),
    constraint pk_code_recovery primary key (code),
    constraint fk_pk_code_recovery_0 foreign key (user_uid) references tb_user(id)
)