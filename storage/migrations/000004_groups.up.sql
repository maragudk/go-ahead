-- groups of users are named and have a secret that's used to join the group.
create table groups (
  id uuid primary key default gen_random_uuid(),
  name text not null,
  invitation_secret bytea not null,
  created timestamp not null default now(),
  updated timestamp not null default now()
);

-- group_membership defines who's part of a group.
create table group_membership (
  group_id uuid references groups(id) on delete cascade,
  user_id uuid references users(id) on delete cascade,
  primary key (group_id, user_id)
);

-- group_ownership defines who owns a group.
-- Note that a user cannot be deleted if the user is a group owner.
create table group_ownership (
   group_id uuid references groups(id) on delete cascade,
   user_id uuid references users(id),
   primary key (group_id, user_id)
);
