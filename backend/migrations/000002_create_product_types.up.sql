--Заполнение справочника
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
INSERT INTO product_types (type_id, name, for_adult) VALUES
(uuid_generate_v4(), 'Молочные продукты', false),
(uuid_generate_v4(), 'Хлебобулочные изделия', false),
(uuid_generate_v4(), 'Мясо и птица', false),
(uuid_generate_v4(), 'Рыба и морепродукты', false),
(uuid_generate_v4(), 'Овощи', false),
(uuid_generate_v4(), 'Фрукты', false),
(uuid_generate_v4(), 'Напитки', false),
(uuid_generate_v4(), 'Кондитерские изделия', false),
(uuid_generate_v4(), 'Замороженные продукты', false),
(uuid_generate_v4(), 'Готовая еда', false),
(uuid_generate_v4(), 'Табачная продукция', true),
(uuid_generate_v4(), 'Алкоголь', true);
