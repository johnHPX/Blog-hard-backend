create table if not exists tb_person(
    id varchar(36) not null,
    user_uid varchar(36) not null,
    name varchar(255) not null,
    telephone varchar(13) not null,
    created_at timestamp not null DEFAULT Now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_person primary key (id, user_uid),
    constraint fk_pk_person foreign key (user_uid) references tb_user(id)
)