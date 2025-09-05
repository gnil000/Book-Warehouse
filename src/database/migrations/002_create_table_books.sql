-- +goose Up
create table books(
    id uuid primary key,
    date_of_writing date not null,
    title text not null,
    author_id uuid not null references authors(id) on delete cascade,
    quantity int not null default 0,
    unique(date_of_writing, title, author_id)
);