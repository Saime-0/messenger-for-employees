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
    token varchar(64) not null,
    comment varchar(512),
    room_seq bigint[] default ARRAY[]::bigint[] not null,
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

create table msg_state
(
    room_id bigint not null,
    msg_count bigint default 0 not null,
    last_msg_id bigint default 0 not null,
    constraint msg_count_pkey
        primary key (room_id),
    constraint msg_count_room_id_fkey
        foreign key (room_id) references rooms
            on delete cascade
);

create table members
(
    emp_id bigint not null,
    room_id bigint not null,
    last_msg_read bigint default 0 not null,
    prev_id bigint,
    foreign key (emp_id) references employees
        on delete cascade,
    foreign key (room_id) references rooms
        on delete cascade,
    foreign key (prev_id) references rooms
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

create function change_count_msg() returns trigger
    language plpgsql
as $$
BEGIN
    if tg_op = 'INSERT' then
        update msg_state
        set
            msg_count = msg_count + 1,
            last_msg_id = new.msg_id
        where msg_state.room_id = new.room_id;
        return new;
    else if tg_op = 'DELETE' then
        update msg_state
        set
            msg_count = msg_count - 1,
            last_msg_id = (
                select m.msg_id
                from messages m
                where m.room_id = old.room_id
                order by m.msg_id desc
                limit 1
            )
        where msg_state.room_id = old.room_id;
        return old;
    end if;
    end if;
    raise exception 'operation could not be detected';
end;
$$;

create trigger on_change_messages_table
    after insert or delete
    on messages
    for each row
execute procedure change_count_msg();

create function create_or_delete_count_msg_row() returns trigger
    language plpgsql
as $$
begin
    if tg_op = 'INSERT' then
        insert into msg_state (room_id) values (new.room_id);
        return new;
    else if tg_op = 'DELETE' then
        delete from msg_state WHERE room_id = old.room_id;
        return old;
    end if;
    end if;
    raise exception 'operation could not be detected';
end;
$$;

create trigger on_create_or_delete_room
    after insert or delete
    on rooms
    for each row
execute procedure create_or_delete_count_msg_row();

create function replace_msg_id() returns trigger
    language plpgsql
as $$
BEGIN
    new.msg_id = (
        SELECT m.last_msg_id+1
        FROM msg_state m
        WHERE m.room_id = new.room_id
    );
    return new;
end;
$$;

create trigger on_insert_message
    before insert
    on messages
    for each row
execute procedure replace_msg_id();

create function add_or_delete_id_in_room_seq() returns trigger
    language plpgsql
as $$
begin
    if tg_op = 'INSERT' then
        new.prev_id = (select e.room_seq[array_length(room_seq, 1)] from employees e where e.emp_id = new.emp_id);
        update employees e set room_seq = e.room_seq || new.room_id::bigint where e.emp_id = 8;
        --         new.prev_id = employees(new.eid).room_seq[len];
--         employees(new.eid).room_seq[len] = new.room_id;
        return new;
    else if tg_op = 'DELETE' then
        update members m set prev_id = old.prev_id where m.emp_id = old.emp_id and m.prev_id = old.room_id;
        update employees e set room_seq = array_remove(room_seq, old.room_id) where e.emp_id = old.emp_id;
        --         Перекинуть prev_id на old.prev_id у всех строк где prev_id был равен old.room_id;
--         delete(employees(eid).room_seq, old.room_id):
        delete from msg_state WHERE room_id = old.room_id;
        return old;
    end if;
    end if;
    raise exception 'operation could not be detected';
end;
$$;

create trigger on_insert_or_delete_member
    after insert or delete
    on members
    for each row
execute procedure add_or_delete_id_in_room_seq();

create function move_room_in_the_sequence(member_emp_id bigint, movable_room_id bigint, prev_room_id bigint) returns void
    language plpgsql
as $$
declare
    origin_room_seq bigint[] := (
        select room_seq
        from employees
        where emp_id = member_emp_id
    );
    cleared_room_seq bigint[] := array_remove(origin_room_seq::bigint[], movable_room_id::bigint);
    prev_room_pos integer := array_position(cleared_room_seq::bigint[], prev_room_id::bigint);
BEGIN
    if movable_room_id is null then
        raise exception 'movable_room_id cannot be null';
    else if movable_room_id = prev_room_id then
        raise exception 'movable_room_id and prev_room_id are the same room';
    else if array_position(origin_room_seq, movable_room_id) is null then
        raise exception 'room_seq does not contain movable_room_id';
    else if prev_room_pos is null then
        raise exception 'room_seq does not contain prev_room_id';
    end if;
    end if;
    end if;
    end if;
    --    [ in members ]
-- тот который стоит перед movable (prev = movable)
    update members
    set prev_id = (
        select prev_id
        from members
        where emp_id = member_emp_id and room_id = movable_room_id
    )
    where emp_id = member_emp_id and prev_id = movable_room_id;

-- тот который стоит перед prev
    update members
    set prev_id = movable_room_id
    where emp_id = member_emp_id and (prev_id = prev_room_id or prev_room_id is null and prev_id is null);

    -- у самого movable
    update members set prev_id = prev_room_id
    where emp_id = member_emp_id and room_id = movable_room_id;

--    [ in employees ]
    update employees e
    set room_seq =
                    cleared_room_seq[1:prev_room_pos]
                    || movable_room_id::bigint ||
                    cleared_room_seq[prev_room_pos+1:]
    where emp_id = member_emp_id;

end;
$$;