INSERT INTO public.users
(id, username, email, password_hash, first_name, family_name, phone_number, date_of_birth, created_at, updated_at)
VALUES
('01KCM0H2P6HD8J7AJ3EKGV7RRM', 'student1', 'student1@gmail.com', '$2a$10$e4Shpgy8quc6OZEwGkurFOzHBkC1aGcEtktFaocR24xK8wrnb94o6', NULL, NULL, NULL, NULL, '2025-12-07 12:25:09.992692+00', '2025-12-07 12:25:09.992692+00'),
('01KCM0HCN2T50ZTD4Q1RFSP3FZ', 'student2', 'student2@gmail.com', '$2a$10$2hXqOBV4SwR2jdGdvPgra.NsFL1mmz2x0.f03bJx/Nh0npjfreH.O', NULL, NULL, NULL, NULL, '2025-12-07 12:25:13.973051+00', '2025-12-07 12:25:13.973051+00');

INSERT INTO public.students
(id, user_id, level, created_at, updated_at)
VALUES
('01KCM0H2P6HD8J7AJ3EKGV7RRM', '01KCM0H2P6HD8J7AJ3EKGV7RRM', 'S-4-Math', '2025-12-07 12:25:09.992692+00', '2025-12-07 12:25:09.992692+00'),
('01KCM0HCN2T50ZTD4Q1RFSP3FZ', '01KCM0HCN2T50ZTD4Q1RFSP3FZ', 'S-3-Math', '2025-12-07 12:25:13.973051+00', '2025-12-07 12:25:13.973051+00');
