CREATE TABLE authors (
  author_id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  bio TEXT
);

CREATE TABLE categories (
  category_id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL
);

CREATE TABLE books (
  book_id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  publication_year INTEGER,
  author_id INTEGER,
  category_id INTEGER,
  price REAL,
  FOREIGN KEY (author_id) REFERENCES authors(author_id),
  FOREIGN KEY (category_id) REFERENCES categories(category_id)
);
