CREATE TABLE IF NOT EXISTS expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255),
    amount INTEGER,
    type VARCHAR(255),
    expense_date TEXT,
    username VARCHAR(255)
);
