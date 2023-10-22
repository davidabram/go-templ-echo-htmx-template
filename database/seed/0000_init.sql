INSERT INTO authors (name, bio) VALUES
('J.K. Rowling', 'Author of the Harry Potter series'),
('George Orwell', 'English novelist and essayist'),
('Jane Austen', 'English novelist known for her romantic fiction');

INSERT INTO categories (name) VALUES
('Fiction'),
('Non-Fiction'),
('Romance');

INSERT INTO books (title, publication_year, author_id, category_id, price) VALUES
('Harry Potter and the Philosopher''s Stone', 1997, 1, 1, 19.99),
('1984', 1949, 2, 1, 9.99),
('Pride and Prejudice', 1813, 3, 3, 4.99);
