create table if not exists tb_post_category (
    id varchar(36) not null,
    post_pid varchar(36) not null,
    category_cid varchar(36) not null,
    created_at timestamp not null DEFAULT Now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_post_category primary key (id, post_pid, category_cid),
    constraint fk_pk_post_category_0 foreign key (post_pid) references tb_post(id),
    constraint fk_pk_post_category_1 foreign key (category_cid) references tb_category(id)
)