create table if not exists tb_response_comment(
    id varchar(36) not null,
    title varchar(255) not null,
    content varchar(2024) not null,
    comment_cid varchar(36) not null,
    user_uid varchar(36) not null,
    created_at timestamp not null DEFAULT Now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_response_comment primary key (id, comment_cid, user_uid),
    constraint fk_pk_response_comment_0 foreign key (user_uid) references tb_user(id),
    constraint fk_pk_response_comment_1 foreign key (comment_cid) references tb_comment(id)
)