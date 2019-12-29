/* -- migrate_up -- */
create table post_likes
(
	id           bigint primary key not null auto_increment,
	post_id      bigint             not null,
	user_id      bigint             not null,
	created_date timestamp    default current_timestamp,
	foreign key (post_id) references posts (id),
	foreign key (user_id) references users (id),
	unique (post_id, user_id)
);
/* -- migrate_down -- */
drop table post_likes;