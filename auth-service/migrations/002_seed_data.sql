-- Seed data for Auth Service (password: password123)

INSERT INTO users (id, username, email, password_hash) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'admin', 'admin@bookshelf.dev', '$2a$10$jNrHsa4DtpbIgeAxQgZngOaNS0hRC6glNsuqpJNdzW7z0Gt5626IG'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'john_reader', 'john@example.com', '$2a$10$jNrHsa4DtpbIgeAxQgZngOaNS0hRC6glNsuqpJNdzW7z0Gt5626IG'),
('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'maria_dev', 'maria@example.com', '$2a$10$jNrHsa4DtpbIgeAxQgZngOaNS0hRC6glNsuqpJNdzW7z0Gt5626IG');
