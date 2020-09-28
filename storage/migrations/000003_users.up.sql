-- users has user information and preferences.
create table users (
  id uuid primary key default gen_random_uuid(),
  name string not null,
  email string unique not null,
  password string not null,
  created timestamp not null default current_timestamp(),
  updated timestamp not null default current_timestamp()
);

-- account_membership defines who is part of an account.
-- Note that it's a many-to-many relationship, so one user can be a part of multiple accounts,
-- and an account can (obviously) have many users.
create table account_membership (
  account_id uuid references accounts(id) on delete cascade,
  user_id uuid references users(id) on delete cascade,
  primary key (account_id, user_id)
);
