# blackout-bot
[Read In English](README.md)

Бот для сповіщення про фактичні відключення електроенергії та повідомлення про планові відключення

## Можливості:
- Надсилає повідомлення про фактичні відключення
- Надсилає повідомлення про планові відключення
- Відображає час, протягом якого електроенергія була присутня або відсутня
- Можливість відключення бота
- Можливість відключити планові повідомлення про відключення
- Можливість оновлення окремого повідомлення про стан електроенергії в реальному часі

### Адмін команди:
| Команда          | Опис                                  |
|------------------|---------------------------------------|
| `turn_bot`       | Включає/виключає бота                 |
| `turn_emergency` | Включає/виключає планові повідомлення |

### Команди:
| Команда          | Опис                                                      |
|------------------|-----------------------------------------------------------|
| `get_id`         | Повертає ChatID для поточного чату                        |
| `get_message_id` | Повертає MessageID для повідомлення на яке була відповідь |
| `curr_status`    | Повертає статус мережі онлайн чи оффлайн                  |

### Як бот працює:
1. Ви налаштовуєте конфігураційний файл
2. Ви додаєте бота в канал або чат і надаєте йому необхідні права
3. Якщо протягом 15 секунд від встановленого ендпоінту не буде відповіді (або 3 спроби ping), бот надішле інформацію про те, що електроенергію відключено
4. Кожного вечора о 19:00 бот буде надсилати графік відключень на наступний день для вказаної групи
5. Також, в окремому повідомленні (якщо вказати це в `.env`), то бот буде оновлювати інформацію в реальному часі

## Залежності:
- Docker `(або golang)`
- Docker-Compose `(або golang)`
- Makefile

## Установка:
1. Склонуйте репозиторій `git clone https://github.com/sergikoio/blackout-bot.git`
2. Створіть .env файл `cp  .env.example .env`
3. Створіть messages.json файл `cp messages.json.example messages.json`
4. Змініть `.env` файл та `messages.json`, якщо треба
5. Запустіть бота:
   - Через Docker: `make build`
   - Або локально через Go:
      - `go mod download`
      - `go run main.go`
        - запустіть від імені `root` користувача - якщо потребуєте використання PING `REQUEST_TYPE`

## ЧаПи
- Що таке ендпоінт?
   - Ви повинні запустити вдома сервер, який буде відповідати на будь-який запит. Або скористайтеся своїм роутером (дивіться щось на кшталт: доступ до роутера глобально). Але для цього у вас повинна бути або постійна IP-адреса, або DDNS на роутері
- Для якого міста складено графік планових відключень?
   - Київ
- Де я можу запустити бота?
   - VPS, Qovery, Heroku, інші...

## Зробіть внесок:
Будь-який внесок вітається, також, якщо в процесі є проблеми - створіть issue