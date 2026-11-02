-- Переводы для Telegram бота (app = 'telegram-bot')
-- Ключи × 3 языка (ru, en, kk)

INSERT INTO translations (code, language, value, app, created_at, updated_at) VALUES

-- =====================================================
-- Общие
-- =====================================================

('telegram.bot.start.welcome', 'ru', 'Добро пожаловать в MarkerLoop Bot!', 'telegram-bot', NOW(), NOW()),
('telegram.bot.start.welcome', 'en', 'Welcome to MarkerLoop Bot!', 'telegram-bot', NOW(), NOW()),
('telegram.bot.start.welcome', 'kk', 'MarkerLoop Bot-қа қош келдіңіз!', 'telegram-bot', NOW(), NOW()),

('telegram.bot.start.message', 'ru', 'Используйте меню ниже для навигации.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.start.message', 'en', 'Use the menu below to navigate.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.start.message', 'kk', 'Навигация үшін төмендегі менюді пайдаланыңыз.', 'telegram-bot', NOW(), NOW()),

('telegram.bot.start.registration_success', 'ru', 'Вы успешно зарегистрированы!', 'telegram-bot', NOW(), NOW()),
('telegram.bot.start.registration_success', 'en', 'You have been successfully registered!', 'telegram-bot', NOW(), NOW()),
('telegram.bot.start.registration_success', 'kk', 'Сіз сәтті тіркелдіңіз!', 'telegram-bot', NOW(), NOW()),

('telegram.bot.command.unknown', 'ru', 'Неизвестная команда. Используйте /menu для навигации.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.command.unknown', 'en', 'Unknown command. Use /menu to navigate.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.command.unknown', 'kk', 'Белгісіз команда. Навигация үшін /menu пайдаланыңыз.', 'telegram-bot', NOW(), NOW()),

('telegram.bot.error.occurred', 'ru', 'Произошла ошибка. Попробуйте позже.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.error.occurred', 'en', 'An error occurred. Please try again later.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.error.occurred', 'kk', 'Қате орын алды. Кейінірек қайталап көріңіз.', 'telegram-bot', NOW(), NOW()),

('telegram.bot.notification.title_format', 'ru', 'Уведомление: %s', 'telegram-bot', NOW(), NOW()),
('telegram.bot.notification.title_format', 'en', 'Notification: %s', 'telegram-bot', NOW(), NOW()),
('telegram.bot.notification.title_format', 'kk', 'Хабарландыру: %s', 'telegram-bot', NOW(), NOW()),

-- =====================================================
-- Авторизация
-- =====================================================

('telegram.bot.button.auth', 'ru', 'Авторизоваться', 'telegram-bot', NOW(), NOW()),
('telegram.bot.button.auth', 'en', 'Log in', 'telegram-bot', NOW(), NOW()),
('telegram.bot.button.auth', 'kk', 'Кіру', 'telegram-bot', NOW(), NOW()),

('telegram.bot.auth.success', 'ru', 'Аккаунт успешно привязан!', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.success', 'en', 'Account linked successfully!', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.success', 'kk', 'Аккаунт сәтті байланыстырылды!', 'telegram-bot', NOW(), NOW()),

('telegram.bot.auth.error', 'ru', 'Произошла ошибка. Попробуйте позже.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.error', 'en', 'An error occurred. Please try again later.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.error', 'kk', 'Қате орын алды. Кейінірек қайталап көріңіз.', 'telegram-bot', NOW(), NOW()),

('telegram.bot.auth.already_linked', 'ru', 'Ваш Telegram аккаунт уже привязан.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.already_linked', 'en', 'Your Telegram account is already linked.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.already_linked', 'kk', 'Сіздің Telegram аккаунтыңыз қазірдің өзінде байланыстырылған.', 'telegram-bot', NOW(), NOW()),

('telegram.bot.auth.token_invalid', 'ru', 'Неверный или истекший токен.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.token_invalid', 'en', 'Invalid or expired token.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.token_invalid', 'kk', 'Жарамсыз немесе мерзімі өткен токен.', 'telegram-bot', NOW(), NOW()),

('telegram.bot.auth.token_expired', 'ru', 'Токен истёк. Запросите новый.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.token_expired', 'en', 'Token expired. Please request a new one.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.token_expired', 'kk', 'Токен мерзімі аяқталды. Жаңасын сұраңыз.', 'telegram-bot', NOW(), NOW()),

('telegram.bot.auth.instructions', 'ru', 'Для авторизации перейдите по ссылке или введите данные ниже.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.instructions', 'en', 'To authorize, follow the link or enter your credentials below.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.instructions', 'kk', 'Авторизация үшін сілтемені басыңыз немесе деректеріңізді төменде енгізіңіз.', 'telegram-bot', NOW(), NOW()),

('telegram.bot.auth.generate_new_link', 'ru', 'Сгенерируйте новую ссылку на сайте.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.generate_new_link', 'en', 'Generate a new link on the website.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.generate_new_link', 'kk', 'Сайтта жаңа сілтеме жасаңыз.', 'telegram-bot', NOW(), NOW()),

('telegram.bot.auth.enter_username', 'ru', 'Введите ваш username:', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.enter_username', 'en', 'Enter your username:', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.enter_username', 'kk', 'Username-ді енгізіңіз:', 'telegram-bot', NOW(), NOW()),

('telegram.bot.auth.username_required', 'ru', 'Username обязателен для заполнения.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.username_required', 'en', 'Username is required.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.username_required', 'kk', 'Username міндетті.', 'telegram-bot', NOW(), NOW()),

('telegram.bot.auth.enter_password', 'ru', 'Введите ваш пароль:', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.enter_password', 'en', 'Enter your password:', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.enter_password', 'kk', 'Құпия сөзді енгізіңіз:', 'telegram-bot', NOW(), NOW()),

('telegram.bot.auth.password_required', 'ru', 'Пароль обязателен для заполнения.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.password_required', 'en', 'Password is required.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.password_required', 'kk', 'Құпия сөз міндетті.', 'telegram-bot', NOW(), NOW()),

('telegram.bot.auth.user_not_found_enter_email', 'ru', 'Пользователь не найден. Введите email для регистрации:', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.user_not_found_enter_email', 'en', 'User not found. Enter your email to register:', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.user_not_found_enter_email', 'kk', 'Пайдаланушы табылмады. Тіркелу үшін email енгізіңіз:', 'telegram-bot', NOW(), NOW()),

('telegram.bot.auth.login_error', 'ru', 'Ошибка входа. Проверьте данные и попробуйте снова.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.login_error', 'en', 'Login error. Check your credentials and try again.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.login_error', 'kk', 'Кіру қатесі. Деректеріңізді тексеріп, қайталап көріңіз.', 'telegram-bot', NOW(), NOW()),

('telegram.bot.auth.user_exists_use_login', 'ru', 'Пользователь уже существует. Используйте вход.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.user_exists_use_login', 'en', 'User already exists. Please log in.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.user_exists_use_login', 'kk', 'Пайдаланушы қазірдің өзінде бар. Кіруді пайдаланыңыз.', 'telegram-bot', NOW(), NOW()),

('telegram.bot.auth.registration_error', 'ru', 'Ошибка регистрации. Попробуйте позже.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.registration_error', 'en', 'Registration error. Please try again later.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.registration_error', 'kk', 'Тіркелу қатесі. Кейінірек қайталап көріңіз.', 'telegram-bot', NOW(), NOW()),

('telegram.bot.auth.required', 'ru', 'Для продолжения необходимо авторизоваться.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.required', 'en', 'Please authorize to continue.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.required', 'kk', 'Жалғастыру үшін авторизациядан өтіңіз.', 'telegram-bot', NOW(), NOW()),

('telegram.bot.auth.web_login_success', 'ru', 'Вы успешно авторизованы на сайте! Вернитесь в браузер.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.web_login_success', 'en', 'You are successfully logged in! Return to the browser.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.auth.web_login_success', 'kk', 'Сіз сайтта сәтті авторизацияланыңыз! Браузерге оралыңыз.', 'telegram-bot', NOW(), NOW()),

-- =====================================================
-- Меню — кнопки
-- =====================================================

('telegram.bot.menu.main', 'ru', 'Главное меню', 'telegram-bot', NOW(), NOW()),
('telegram.bot.menu.main', 'en', 'Main menu', 'telegram-bot', NOW(), NOW()),
('telegram.bot.menu.main', 'kk', 'Басты мәзір', 'telegram-bot', NOW(), NOW()),

('telegram.bot.button.profile', 'ru', 'Мой профиль', 'telegram-bot', NOW(), NOW()),
('telegram.bot.button.profile', 'en', 'My profile', 'telegram-bot', NOW(), NOW()),
('telegram.bot.button.profile', 'kk', 'Менің профилім', 'telegram-bot', NOW(), NOW()),

('telegram.bot.button.settings', 'ru', 'Настройки', 'telegram-bot', NOW(), NOW()),
('telegram.bot.button.settings', 'en', 'Settings', 'telegram-bot', NOW(), NOW()),
('telegram.bot.button.settings', 'kk', 'Баптаулар', 'telegram-bot', NOW(), NOW()),

('telegram.bot.button.logout', 'ru', 'Выход', 'telegram-bot', NOW(), NOW()),
('telegram.bot.button.logout', 'en', 'Log out', 'telegram-bot', NOW(), NOW()),
('telegram.bot.button.logout', 'kk', 'Шығу', 'telegram-bot', NOW(), NOW()),

('telegram.bot.button.users', 'ru', 'Пользователи', 'telegram-bot', NOW(), NOW()),
('telegram.bot.button.users', 'en', 'Users', 'telegram-bot', NOW(), NOW()),
('telegram.bot.button.users', 'kk', 'Пайдаланушылар', 'telegram-bot', NOW(), NOW()),

('telegram.bot.button.stats', 'ru', 'Статистика', 'telegram-bot', NOW(), NOW()),
('telegram.bot.button.stats', 'en', 'Statistics', 'telegram-bot', NOW(), NOW()),
('telegram.bot.button.stats', 'kk', 'Статистика', 'telegram-bot', NOW(), NOW()),

('telegram.bot.button.notifications', 'ru', 'Уведомления', 'telegram-bot', NOW(), NOW()),
('telegram.bot.button.notifications', 'en', 'Notifications', 'telegram-bot', NOW(), NOW()),
('telegram.bot.button.notifications', 'kk', 'Хабарландырулар', 'telegram-bot', NOW(), NOW()),

('telegram.bot.button.open_site', 'ru', 'Открыть сайт', 'telegram-bot', NOW(), NOW()),
('telegram.bot.button.open_site', 'en', 'Open site', 'telegram-bot', NOW(), NOW()),
('telegram.bot.button.open_site', 'kk', 'Сайтты ашу', 'telegram-bot', NOW(), NOW()),

('telegram.bot.button.back', 'ru', '← Назад', 'telegram-bot', NOW(), NOW()),
('telegram.bot.button.back', 'en', '← Back', 'telegram-bot', NOW(), NOW()),
('telegram.bot.button.back', 'kk', '← Артқа', 'telegram-bot', NOW(), NOW()),

-- =====================================================
-- Профиль
-- =====================================================

('telegram.bot.profile.title', 'ru', '👤 Мой профиль', 'telegram-bot', NOW(), NOW()),
('telegram.bot.profile.title', 'en', '👤 My profile', 'telegram-bot', NOW(), NOW()),
('telegram.bot.profile.title', 'kk', '👤 Менің профилім', 'telegram-bot', NOW(), NOW()),

('telegram.bot.profile.role', 'ru', 'Роль', 'telegram-bot', NOW(), NOW()),
('telegram.bot.profile.role', 'en', 'Role', 'telegram-bot', NOW(), NOW()),
('telegram.bot.profile.role', 'kk', 'Рөл', 'telegram-bot', NOW(), NOW()),

('telegram.bot.profile.registered', 'ru', 'Регистрация', 'telegram-bot', NOW(), NOW()),
('telegram.bot.profile.registered', 'en', 'Registered', 'telegram-bot', NOW(), NOW()),
('telegram.bot.profile.registered', 'kk', 'Тіркелу', 'telegram-bot', NOW(), NOW()),

-- =====================================================
-- Настройки / Язык
-- =====================================================

('telegram.bot.language.choose', 'ru', 'Выберите язык:', 'telegram-bot', NOW(), NOW()),
('telegram.bot.language.choose', 'en', 'Choose language:', 'telegram-bot', NOW(), NOW()),
('telegram.bot.language.choose', 'kk', 'Тілді таңдаңыз:', 'telegram-bot', NOW(), NOW()),

('telegram.bot.language.changed', 'ru', 'Язык успешно изменён.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.language.changed', 'en', 'Language changed successfully.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.language.changed', 'kk', 'Тіл сәтті өзгертілді.', 'telegram-bot', NOW(), NOW()),

-- =====================================================
-- Выход
-- =====================================================

('telegram.bot.logout.success', 'ru', 'Вы вышли из аккаунта.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.logout.success', 'en', 'You have been logged out.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.logout.success', 'kk', 'Сіз аккаунттан шықтыңыз.', 'telegram-bot', NOW(), NOW()),

-- =====================================================
-- Пользователи (superadmin)
-- =====================================================

('telegram.bot.users.title', 'ru', '👥 Пользователи', 'telegram-bot', NOW(), NOW()),
('telegram.bot.users.title', 'en', '👥 Users', 'telegram-bot', NOW(), NOW()),
('telegram.bot.users.title', 'kk', '👥 Пайдаланушылар', 'telegram-bot', NOW(), NOW()),

('telegram.bot.users.total', 'ru', 'Всего', 'telegram-bot', NOW(), NOW()),
('telegram.bot.users.total', 'en', 'Total', 'telegram-bot', NOW(), NOW()),
('telegram.bot.users.total', 'kk', 'Барлығы', 'telegram-bot', NOW(), NOW()),

-- =====================================================
-- Статистика (superadmin)
-- =====================================================

('telegram.bot.stats.title', 'ru', '📊 Статистика', 'telegram-bot', NOW(), NOW()),
('telegram.bot.stats.title', 'en', '📊 Statistics', 'telegram-bot', NOW(), NOW()),
('telegram.bot.stats.title', 'kk', '📊 Статистика', 'telegram-bot', NOW(), NOW()),

('telegram.bot.stats.total_users', 'ru', 'Пользователей', 'telegram-bot', NOW(), NOW()),
('telegram.bot.stats.total_users', 'en', 'Users', 'telegram-bot', NOW(), NOW()),
('telegram.bot.stats.total_users', 'kk', 'Пайдаланушылар', 'telegram-bot', NOW(), NOW()),

('telegram.bot.stats.total_roles', 'ru', 'Ролей', 'telegram-bot', NOW(), NOW()),
('telegram.bot.stats.total_roles', 'en', 'Roles', 'telegram-bot', NOW(), NOW()),
('telegram.bot.stats.total_roles', 'kk', 'Рөлдер', 'telegram-bot', NOW(), NOW()),

('telegram.bot.stats.deploy_subscribers', 'ru', 'Подписчиков деплоя', 'telegram-bot', NOW(), NOW()),
('telegram.bot.stats.deploy_subscribers', 'en', 'Deploy subscribers', 'telegram-bot', NOW(), NOW()),
('telegram.bot.stats.deploy_subscribers', 'kk', 'Деплой жазылушылары', 'telegram-bot', NOW(), NOW()),

-- =====================================================
-- Уведомления о деплоях (superadmin)
-- =====================================================

('telegram.bot.notifications.title', 'ru', '🔔 Уведомления о деплоях', 'telegram-bot', NOW(), NOW()),
('telegram.bot.notifications.title', 'en', '🔔 Deploy notifications', 'telegram-bot', NOW(), NOW()),
('telegram.bot.notifications.title', 'kk', '🔔 Деплой хабарландырулары', 'telegram-bot', NOW(), NOW()),

('telegram.bot.notifications.enabled', 'ru', 'Статус: Включено ✅', 'telegram-bot', NOW(), NOW()),
('telegram.bot.notifications.enabled', 'en', 'Status: Enabled ✅', 'telegram-bot', NOW(), NOW()),
('telegram.bot.notifications.enabled', 'kk', 'Күй: Қосылған ✅', 'telegram-bot', NOW(), NOW()),

('telegram.bot.notifications.disabled', 'ru', 'Статус: Отключено ❌', 'telegram-bot', NOW(), NOW()),
('telegram.bot.notifications.disabled', 'en', 'Status: Disabled ❌', 'telegram-bot', NOW(), NOW()),
('telegram.bot.notifications.disabled', 'kk', 'Күй: Өшірілген ❌', 'telegram-bot', NOW(), NOW()),

('telegram.bot.notifications.description', 'ru', 'Получать уведомления при обновлении сервисов в ArgoCD.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.notifications.description', 'en', 'Receive notifications when services are updated in ArgoCD.', 'telegram-bot', NOW(), NOW()),
('telegram.bot.notifications.description', 'kk', 'ArgoCD-де сервистер жаңартылғанда хабарландыру алу.', 'telegram-bot', NOW(), NOW()),

('telegram.bot.notifications.enable', 'ru', 'Включить', 'telegram-bot', NOW(), NOW()),
('telegram.bot.notifications.enable', 'en', 'Enable', 'telegram-bot', NOW(), NOW()),
('telegram.bot.notifications.enable', 'kk', 'Қосу', 'telegram-bot', NOW(), NOW()),

('telegram.bot.notifications.disable', 'ru', 'Отключить', 'telegram-bot', NOW(), NOW()),
('telegram.bot.notifications.disable', 'en', 'Disable', 'telegram-bot', NOW(), NOW()),
('telegram.bot.notifications.disable', 'kk', 'Өшіру', 'telegram-bot', NOW(), NOW())

ON CONFLICT (code, language, app) DO UPDATE SET
    value = EXCLUDED.value,
    updated_at = NOW();
