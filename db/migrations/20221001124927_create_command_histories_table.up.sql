CREATE TABLE IF NOT EXISTS command_histories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(255),
    command_name VARCHAR(255),
    step INTEGER
);
