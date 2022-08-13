create table if not exists tb_post(
    id varchar(36) not null,
    title varchar(255) not null,
    content text not null,
    created_at timestamp not null DEFAULT Now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_post primary key (id)
)