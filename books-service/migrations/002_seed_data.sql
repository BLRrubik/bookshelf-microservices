-- Seed data for Books Service

INSERT INTO books (id, title, author, description, isbn, published_year, created_by) VALUES
('f5eebc99-9c0b-4ef8-bb6d-6bb9bd380b01', 'Чистая архитектура', 'Роберт Мартин', 'Практическое руководство по программной архитектуре.', '978-5-4461-0772-8', 2018, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
('f6eebc99-9c0b-4ef8-bb6d-6bb9bd380b02', 'Чистый код', 'Роберт Мартин', 'Справочник для создания читаемого кода.', '978-5-4461-0960-9', 2008, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
('fa2ebc99-9c0b-4ef8-bb6d-6bb9bd380b08', 'Высоконагруженные приложения', 'Мартин Клеппман', 'DDIA. Принципы проектирования надёжных систем.', '978-5-4461-0512-0', 2017, 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11');

INSERT INTO reviews (book_id, user_id, rating, title, content) VALUES
('f5eebc99-9c0b-4ef8-bb6d-6bb9bd380b01', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 5, 'Фундаментальная книга', 'Обязательна к прочтению.'),
('fa2ebc99-9c0b-4ef8-bb6d-6bb9bd380b08', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 5, 'Шедевр', 'Лучшая книга по распределённым системам.');
