alter table entries add column changed_at timestamp with time zone;
update entries set changed_at = published_at;
alter table entries alter column changed_at set not null;

create table api_keys (
    id serial not null,
    user_id int not null references users(id) on delete cascade,
    token text not null unique,
    description text not null,
    last_used_at timestamp with time zone,
    created_at timestamp with time zone default now(),
    primary key(id),
    unique (user_id, description)
);


alter table entries add column share_code text not null default '';
create unique index entries_share_code_idx on entries using btree(share_code) where share_code <> '';
