create type room_type as enum ('TALK', 'BLOG');

create table admins
(
    admin_id bigserial,
    token varchar(32) not null,
    email varchar(64) not null,
    primary key (admin_id)
);

create table employees
(
    emp_id bigserial,
    first_name varchar(32) not null,
    last_name varchar(32) not null,
    email varchar(64) not null,
    phone_number varchar(32) not null,
    joined_at bigint default unix_utc_now() not null,
    token varchar(32) not null,
    comment varchar(512),
    primary key (emp_id)
);

create table tags
(
    tag_id bigserial,
    name varchar(64) not null,
    primary key (tag_id)
);

create table rooms
(
    room_id bigserial,
    name varchar(64) not null,
    view room_type not null,
    primary key (room_id)
);

create table msg_count
(
    room_id bigint not null,
    val bigint default 0 not null,
    primary key (room_id),
    foreign key (room_id) references rooms
        on delete cascade
);

create table members
(
    emp_id bigint not null,
    room_id bigint not null,
    last_msg_read bigint default 0 not null,
    foreign key (emp_id) references employees
        on delete cascade,
    foreign key (room_id) references rooms
        on delete cascade
);

create table positions
(
    emp_id bigint not null,
    tag_id bigint not null,
    foreign key (emp_id) references employees
        on delete cascade,
    foreign key (tag_id) references tags
        on delete cascade
);

create table refresh_sessions
(
    id bigserial,
    emp_id bigint,
    refresh_token varchar(32) not null,
    expires_at bigint not null,
    created_at bigint default unix_utc_now() not null,
    foreign key (emp_id) references employees
);

create table messages
(
    room_id bigint not null,
    msg_id bigint not null,
    emp_id bigint not null,
    target_id bigint,
    body varchar(2048) not null,
    created_at bigint default unix_utc_now() not null,
    foreign key (room_id) references rooms
        on delete cascade,
    foreign key (emp_id) references employees
        on delete cascade
);

create function unix_utc_now(bigint DEFAULT 0) returns bigint
    language sql
as $$
SELECT (date_part('epoch'::text, now()))::bigint + $1
$$;