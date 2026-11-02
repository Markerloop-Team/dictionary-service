-- Добавление тестовых данных: Казахстан и Алматы
-- Выполните этот скрипт в базе данных dictionary

-- 1. Добавляем страну Казахстан
INSERT INTO
    countries (code, created_at, updated_at)
VALUES ('KZ', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- 2. Добавляем переводы для страны Казахстан
INSERT INTO
    translations (
        code,
        language,
        value,
        app,
        created_at,
        updated_at
    )
VALUES (
        'country.KZ',
        'ru',
        'Казахстан',
        NULL,
        NOW(),
        NOW()
    ),
    (
        'country.KZ',
        'kk',
        'Қазақстан',
        NULL,
        NOW(),
        NOW()
    ),
    (
        'country.KZ',
        'en',
        'Kazakhstan',
        NULL,
        NOW(),
        NOW()
    )
ON CONFLICT (code, language, app) DO NOTHING;

-- 3. Добавляем город Алматы
INSERT INTO
    cities (
        code,
        country_code,
        created_at,
        updated_at
    )
VALUES ('ALA', 'KZ', NOW(), NOW())
ON CONFLICT (code, country_code) DO NOTHING;

-- 4. Добавляем переводы для города Алматы
INSERT INTO
    translations (
        code,
        language,
        value,
        app,
        created_at,
        updated_at
    )
VALUES (
        'city.ALA.KZ',
        'ru',
        'Алматы',
        NULL,
        NOW(),
        NOW()
    ),
    (
        'city.ALA.KZ',
        'kk',
        'Алматы',
        NULL,
        NOW(),
        NOW()
    ),
    (
        'city.ALA.KZ',
        'en',
        'Almaty',
        NULL,
        NOW(),
        NOW()
    )
ON CONFLICT (code, language, app) DO NOTHING;







