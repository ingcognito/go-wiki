CREATE DATABASE bot;
\c bot
CREATE table pages (
    title varchar(255),
    extract text,
    link text
)