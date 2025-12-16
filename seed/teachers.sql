INSERT INTO public.users
(id, username, email, password_hash, first_name, family_name, phone_number, date_of_birth, created_at, updated_at)
VALUES
('01KCM0PGS19BQGM18CPZW49B08', 'teacher1', 'teacher1@gmail.com', '$2a$10$e4Shpgy8quc6OZEwGkurFOzHBkC1aGcEtktFaocR24xK8wrnb94o6', NULL, NULL, NULL, NULL, '2025-12-07 12:25:09.992692+00', '2025-12-07 12:25:09.992692+00'),
('01KCM0Q9X5H5MJ9BTGXNXSBV9V', 'teacher2', 'teacher2@gmail.com', '$2a$10$2hXqOBV4SwR2jdGdvPgra.NsFL1mmz2x0.f03bJx/Nh0npjfreH.O', NULL, NULL, NULL, NULL, '2025-12-07 12:25:13.973051+00', '2025-12-07 12:25:13.973051+00');

INSERT INTO public.teachers
(id, user_id, created_at, updated_at)
VALUES
('01KCM0PGS19BQGM18CPZW49B08', '01KCM0PGS19BQGM18CPZW49B08', '2025-12-07 12:25:09.992692+00', '2025-12-07 12:25:09.992692+00'),
('01KCM0Q9X5H5MJ9BTGXNXSBV9V', '01KCM0Q9X5H5MJ9BTGXNXSBV9V', '2025-12-07 12:25:13.973051+00', '2025-12-07 12:25:13.973051+00');
