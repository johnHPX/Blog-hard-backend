create table if not exists tb_user(
    id varchar(36) not null, 
    nick varchar(255) not null,
    email varchar(255) not null,
    secret varchar(255) not null,
    kind varchar(10) not null,
    created_at timestamp not null DEFAULT Now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_user primary key (id)
)