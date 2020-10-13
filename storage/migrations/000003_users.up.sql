-- users has user information and preferences.
create table users (
  id uuid primary key default gen_random_uuid(),
  name text not null,
  email text unique not null,
  password text not null,
  created timestamp not null default now(),
  updated timestamp not null default now()
);

-- account_membership defines who is part of an account.
-- Note that it's a many-to-many relationship, so one user can be a part of multiple accounts,
-- and an account can (obviously) have many users.
create table account_membership (
  account_id uuid references accounts(id) on delete cascade,
  user_id uuid references users(id) on delete cascade,
  primary key (account_id, user_id)
);
