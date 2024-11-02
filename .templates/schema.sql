create table users (
  id integer primary key,
  email varchar not null,
  pass varchar not null
);

create table user_items (
  id integer primary key,
  user_id integer not null,
  name varchar not null
);

create table categories (
  id integer primary key,
  title varchar not null,
  slug varchar not null
);

create table categories_users (
  user_id integer not null,
  category_id integer not null
);
