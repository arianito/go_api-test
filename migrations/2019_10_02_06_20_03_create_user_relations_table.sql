/* -- migrate_up -- */
create table user_relations
(
	id           bigint primary key not null auto_increment,
	user_a      bigint             not null,
	post_b      bigint             not null,
	created_date timestamp default current_timestamp,
	foreign key (user_a) references users (id),
	foreign key (post_b) references users (id)
);
/* -- migrate_down -- */
drop table user_relations;