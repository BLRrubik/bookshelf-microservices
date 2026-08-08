-- Seed data for Auth Service
-- Password for all users: password123 (bcrypt hash)

INSERT INTO users (id, username, email, password_hash, created_at, updated_at) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'admin', 'admin@bookshelf.dev', '$2a$10$jNrHsa4DtpbIgeAxQgZngOaNS0hRC6glNsuqpJNdzW7z0Gt5626IG', '2024-01-01 10:00:00+00', '2024-01-01 10:00:00+00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'john_reader', 'john@example.com', '$2a$10$jNrHsa4DtpbIgeAxQgZngOaNS0hRC6glNsuqpJNdzW7z0Gt5626IG', '2024-01-05 14:30:00+00', '2024-01-05 14:30:00+00'),
('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'maria_dev', 'maria@example.com', '$2a$10$jNrHsa4DtpbIgeAxQgZngOaNS0hRC6glNsuqpJNdzW7z0Gt5626IG', '2024-01-10 09:15:00+00', '2024-01-10 09:15:00+00'),
('d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'alex_coder', 'alex@example.com', '$2a$10$jNrHsa4DtpbIgeAxQgZngOaNS0hRC6glNsuqpJNdzW7z0Gt5626IG', '2024-01-12 16:45:00+00', '2024-01-12 16:45:00+00'),
('e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'elena_arch', 'elena@example.com', '$2a$10$jNrHsa4DtpbIgeAxQgZngOaNS0hRC6glNsuqpJNdzW7z0Gt5626IG', '2024-01-15 11:20:00+00', '2024-01-15 11:20:00+00');
