create type room_type as enum ('TALK', 'BLOG');


create or replace function unix_utc_now(bigint DEFAULT 0) returns bigint
    language sql
as $$
SELECT (date_part('epoch'::text, now()))::bigint + $1
$$;

create table admins (
                        id bigserial primary key,
                        token varchar(64),
                        email varchar(64)
);

create table employees (
                           id bigserial primary key,
                           first_name varchar(64),
                           last_name varchar(64),
                           email varchar(128),
                           phone_number varchar(64),
                           joined_at bigint default unix_utc_now(),
                           password_hash varchar(64),
                           comment varchar(1024),
                           room_seq bigint[] default ARRAY[NULL::bigint]
);

create table tags (
                      id bigserial primary key,
                      name varchar(128)
);

create table rooms (
                       id bigserial primary key,
                       name varchar(128),
                       view room_type
);

create table messages (
                          id bigserial primary key,
                          room_id bigint references rooms on delete cascade,
                          emp_id bigint null references employees (id) on delete set null,
                          reply_id bigint null references messages on delete set null,
                          body varchar(2048),
                          created_at bigint default unix_utc_now(),
                          prev bigint,
                          next bigint
);

create table msg_state (
                           room_id bigint primary key references rooms on delete cascade,
                           msg_count bigint default 0,
                           last_msg_id bigint null references messages (id) on delete set null
);

create table members (
                         emp_id bigint references employees on delete cascade,
                         room_id bigint references rooms on delete cascade,
                         last_msg_read bigint null references messages,
                         notify bool default false,
                         prev_id bigint null references rooms
);

create table positions (
                           emp_id bigint references employees on delete cascade,
                           tag_id bigint references tags on delete cascade
);

create table refresh_sessions (
                                  id bigserial primary key,
                                  emp_id bigint references employees,
                                  refresh_token varchar(64),
                                  expires_at bigint,
                                  created_at bigint default unix_utc_now()
);

-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

create or replace function change_count_msg() returns trigger
    language plpgsql
as $$
BEGIN
    if tg_op = 'INSERT' then
        update msg_state
        set msg_count   = msg_count + 1,
            last_msg_id = new.id
        where msg_state.room_id = new.room_id;
        return new;
    else
        if tg_op = 'DELETE' then
            update msg_state
            set msg_count   = msg_count - 1,
                last_msg_id = (select m.id
                               from messages m
                               where m.room_id = old.room_id
                               order by m.id desc
                               limit 1)
            where msg_state.room_id = old.room_id;

            update messages m
            set next = old.next
            where m.id = old.prev;

            update messages m
            set prev = old.prev
            where m.id = old.next;

            return old;
        end if;
    end if;
    raise exception 'operation could not be detected';
end;
$$;

drop trigger if exists on_change_messages_table ON messages;
create trigger on_change_messages_table
    after insert or delete
    on messages
    for each row
execute procedure change_count_msg();

create or replace function create_count_msg_row() returns trigger
    language plpgsql
as $$
begin
    insert into msg_state (room_id) values (new.id);
    return new;
end;
$$;

drop trigger if exists on_create_room ON rooms;
create trigger on_create_room
    after insert
    on rooms
    for each row
execute procedure create_count_msg_row();

create or replace function replace_msg_prev_and_next() returns trigger
    language plpgsql
as $$
BEGIN
    new.next = null;
    new.prev = (select m.last_msg_id
                from msg_state m
                where m.room_id = new.room_id);
    update messages m
    set next = new.id
    where m.id = new.prev;
    return new;
end;
$$;

drop trigger if exists on_insert_message ON messages;
create trigger on_insert_message
    before insert
    on messages
    for each row
execute procedure replace_msg_prev_and_next();

create or replace function add_or_delete_id_in_room_seq() returns trigger
    language plpgsql
as $$
begin
    if tg_op = 'INSERT' then
        new.prev_id = (select e.room_seq[array_length(room_seq, 1)] from employees e where e.id = new.emp_id);
        update employees e set room_seq = e.room_seq || new.room_id::bigint where e.id = new.emp_id;
        --         new.prev_id = employees(new.eid).room_seq[len];
--         employees(new.eid).room_seq[len] = new.room_id;
        return new;
    else if tg_op = 'DELETE' then
        update members m set prev_id = old.prev_id where m.emp_id = old.emp_id and m.prev_id = old.room_id;
        update employees e set room_seq = array_remove(room_seq, old.room_id) where e.id = old.emp_id;
        --         Перекинуть prev_id на old.prev_id у всех строк где prev_id был равен old.room_id;
--         delete(employees(eid).room_seq, old.room_id):
--         delete from msg_state WHERE room_id = old.room_id; WHAT????
        return old;
    end if;
    end if;
    raise exception 'operation could not be detected';
end
$$;

drop trigger if exists on_insert_or_delete_member ON members;
create trigger on_insert_or_delete_member
    before insert or delete
    on members
    for each row
execute procedure add_or_delete_id_in_room_seq();

create or replace function move_room_in_the_sequence(member_emp_id bigint, movable_room_id bigint, prev_room_id bigint) returns void
    language plpgsql
as $$
declare
    origin_room_seq bigint[] := (
        select room_seq
        from employees
        where id = member_emp_id
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
    where id = member_emp_id;

end;
$$;

create or replace function load_emp_rooms(ptrs text[], empids bigint[], limits integer[], offsets integer[]) returns TABLE(ptr text, orderPos integer,  room_id bigint, name character varying, view room_type, last_msg_read bigint, last_msg_id bigint, notify bool)
    language plpgsql
as $$
declare emp_rooms text[][] = array[]::text[][]; -- {{ptrs}, {empids}, {seqs}, {room_seq}}
begin
    FOR i IN 1..array_length(ptrs, 1) LOOP
            emp_rooms[1] = array_append(emp_rooms[1]::text[], ptrs[i]::text);
            emp_rooms[2] = array_append(emp_rooms[2]::text[], empIDs[i]::text);
            emp_rooms[3] = array_append(emp_rooms[3]::text[], (
                SELECT room_seq[
                           (select coalesce(array_length(room_seq, 1) - coalesce(offsets[i], 0)+1-limits[i], 1)):
                           (select array_length(room_seq, 1) - coalesce(offsets[i], 0))]
                FROM employees WHERE id = empIDs[i]
            )::text);
            emp_rooms[4] = array_append(emp_rooms[4]::text[], (
                SELECT room_seq
                FROM employees WHERE id = empIDs[i]
            )::text);
        END LOOP;
    return query SELECT inp.ptr,
                        array_position(inp.room_seq::bigint[], r.id),
                        coalesce(r.id, 0),
                        coalesce(r.name, ''),
                        coalesce(r.view, 'TALK'),
                        m.last_msg_read,
                        c.last_msg_id,
                        m.notify
                 FROM unnest(emp_rooms[1]::text[], emp_rooms[2]::bigint[], emp_rooms[3]::text[][], emp_rooms[4]::text[][]) inp(ptr, empid, seq, room_seq)
                          LEFT JOIN members m ON m.emp_id = inp.empid
                          LEFT JOIN rooms r on r.id = m.room_id AND r.id = ANY (inp.seq::bigint[])
                          LEFT JOIN msg_state c on r.id = c.room_id;
end;
$$;