-- +goose Up
create table authors(
    id uuid primary key,
    date_of_birth date not null,
    first_name text not null,
    second_name text not null,
    surname text not null,
    unique(date_of_birth, first_name, second_name, surname)
);