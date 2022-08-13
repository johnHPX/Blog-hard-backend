create table if not exists tb_category(
    id varchar(36) not null,
    name varchar(255) not null,
    created_at timestamp not null DEFAULT Now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_category primary key (id)
)