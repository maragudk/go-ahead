create extension if not exists "pgcrypto";

-- accounts have names and will be used primarily for billing.
create table accounts (
  id uuid primary key default gen_random_uuid(),
  name text not null,
  created timestamp not null default now(),
  updated timestamp not null default now()
);
