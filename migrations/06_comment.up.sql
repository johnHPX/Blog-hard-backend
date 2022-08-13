create table if not exists tb_comment(
    id varchar(36) not null,
    title varchar(255) not null,
    content varchar(2024) not null,
    user_uid varchar(36) not null,
    post_pid varchar(36) not null,
    created_at timestamp not null DEFAULT Now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_comment primary key (id, user_uid, post_pid),
    constraint fk_pk_comment_0 foreign key (user_uid) references tb_user(id),
    constraint fk_pk_comment_1 foreign key (post_pid) references tb_post(id)
)