drop table admins;

drop table msg_state;

drop table members;

drop table positions;

drop table tags;

drop table refresh_sessions;

drop table messages;

drop table employees;

drop table rooms;

drop type room_type;

drop function unix_utc_now(bigint);

drop function change_count_msg();

drop function create_or_delete_count_msg_row();

drop function replace_msg_id();

drop function add_or_delete_id_in_room_seq();

drop function move_room_in_the_sequence(bigint, bigint, bigint);

drop function load_emp_rooms(text[], bigint[], integer[], integer[]);