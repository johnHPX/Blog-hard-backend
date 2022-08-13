create table if not exists tb_configs(
    id serial not null,
    collors varchar[] not null,
    links varchar[] not null,
    menuAs varchar[] not null,
    bannerURL varchar not null,
    created_at timestamp not null DEFAULT Now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_configs primary key (id)
)