CREATE TABLE IF NOT EXISTS expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255),
    price INTEGER,
    type VARCHAR(255),
    expense_date TEXT,
    created_at TEXT,
    updated_at TEXT
);
