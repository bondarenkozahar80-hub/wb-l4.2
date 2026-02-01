# MyGrep

## Описание

Распределенная CLI-утилита для поиска строк в файлах. Файл делится на фрагменты, каждый сервер обрабатывает свою часть через HTTP. Результат формируется при достижении кворума (N/2+1 ответов).

**Основные возможности:**
- Распределенная обработка файлов и потоков с разбиением на фрагменты
- Кворум для обеспечения надежности (N/2+1)
- Параллельная обработка через горутины и каналы
- Поддержка регулярных выражений
- Чтение из файла или стандартного ввода (stdin)
- Сравнение результатов с оригинальной утилитой grep
- HTTP-транспорт для обмена данными между серверами

## CLI интерфейс

Утилита работает в двух режимах:

### Режим воркера (сервера)

Воркер запускается с флагом `--server` и обрабатывает запросы от координатора.

### Режим координатора (клиента)

Координатор принимает шаблон поиска и источник данных (файл или stdin), делит данные на части, рассылает задачи на воркеры параллельно и собирает результаты.

___
## Линтер

Проект использует **golangci-lint** для проверки качества кода.

### Запуск линтера

```bash
make linter
```

## Установка и запуск проекта

### 1. Клонирование репозитория

```bash
git clone https://github.com/kstsm/wb-l4.2
cd wb-l4.2
```

### 2. Сборка проекта

**С помощью Make:**
```bash
make build
```

**Или вручную:**
```bash
go build -o mygrep
```

### 3. Запуск воркеров

Запустите несколько воркеров в разных терминалах:

**С помощью Make:**
```bash
# Терминал 1
make worker1

# Терминал 2
make worker2

# Терминал 3
make worker3
```

**Или вручную:**
```bash
# Терминал 1
./mygrep --server --addr :9001

# Терминал 2
./mygrep --server --addr :9002

# Терминал 3
./mygrep --server --addr :9003
```

### 4. Запуск координатора
**Вручную:**
```bash
./mygrep Проверка --file test_data.txt --nodes localhost:9001,localhost:9002,localhost:9003
```

**Чтение из stdin:**
```bash
cat test_data.txt | ./mygrep Проверка --nodes localhost:9001,localhost:9002,localhost:9003
```

Координатор разделит данные на части, отправит их воркерам параллельно и соберет результаты. Работа завершится при достижении кворума (N/2 + 1 ответов).
___

# CLI команды

## Запуск воркера

**Команда:** `./mygrep --server --addr :PORT`

**Параметры:**

- `--server` (обязательно) - режим сервера (воркер)
- `--addr` (опционально) - адрес сервера (по умолчанию `:9001`)

**Пример:**

```bash
./mygrep --server --addr :9001
```

**Ожидаемый вывод:**

```
[2025/12/30T01:37:54.634] [application] [INFO] [worker.go:17,Run] Worker started on :9001
```

Воркер будет слушать HTTP-запросы на указанном порту и обрабатывать задачи поиска.

---

## Запуск координатора

**Команда:** `./mygrep PATTERN [--file FILE] --nodes NODES`

**Параметры:**

- `PATTERN` (обязательно) - шаблон для поиска (регулярное выражение)
- `--file` (опционально) - путь к файлу для обработки. Если не указан, данные читаются из stdin
- `--nodes` (обязательно) - список серверов через запятую (например: `localhost:9001,localhost:9002,localhost:9003`)
- `-n` (опционально) - показывать номера строк в формате `номер:содержимое`

**Пример:**

```bash
./mygrep Проверка --file test_data.txt --nodes localhost:9001,localhost:9002,localhost:9003
```

**Ожидаемый вывод:**

```
[2025/12/30T01:37:54.634] [application] [INFO] [coordinator.go:40,RunWithOptions] Starting search: pattern=Проверка, source=test_data.txt, servers=3, quorum=2
[2025/12/30T01:37:54.640] [application] [INFO] [coordinator.go:91,func1] Server localhost:9001 responded successfully with 4 matches
[2025/12/30T01:37:54.642] [application] [INFO] [coordinator.go:91,func1] Server localhost:9002 responded successfully with 4 matches
[2025/12/30T01:37:54.644] [application] [INFO] [coordinator.go:91,func1] Server localhost:9003 responded successfully with 4 matches
[2025/12/30T01:37:54.645] [application] [INFO] [coordinator.go:121,RunWithOptions] Quorum reached: 3/3 servers responded successfully
Проверка mygrep
Проверка проекта
Проверка для тестирования
[2025/12/30T01:37:54.646] [application] [INFO] [coordinator.go:134,RunWithOptions] Found matches: 4
```

### Ошибки:

**Не указан шаблон:**

```
[application] [ERROR] [main.go:14,main] application error: parsing arguments: pattern not specified
```

**Чтение из stdin:**

```bash
cat test_data.txt | ./mygrep PATTERN --nodes localhost:9001,localhost:9002,localhost:9003
```

**Не указаны серверы:**

```
[application] [ERROR] [main.go:14,main] application error: parsing arguments: servers not specified
```

**Ошибка чтения файла:**

```
[application] [ERROR] [main.go:14,main] application error: coordinator: reading input: open nonexistent.txt: no such file or directory
```

**Кворум не достигнут:**

```
[application] [ERROR] [main.go:14,main] application error: coordinator: quorum not reached: got 1 successful responses out of 3, required 2 (failed: 2)
```

**Ошибка от сервера:**

```
[2025/12/30T01:37:54.640] [application] [WARN] [coordinator.go:83,func1] Server localhost:9002 failed: sending request to http://localhost:9002/search: dial tcp [::1]:9002: connect: connection refused
```

---

## Примеры использования

### Пример 1: Поиск простого слова

**Файл `test_data.txt`:**
```
Проверка mygrep
Проверка проекта
Проверка для тестирования
Без нее
```

**Команда:**
```bash
./mygrep Проверка --file test_data.txt --nodes localhost:9001,localhost:9002,localhost:9003
```

**Результат:**
```
Проверка mygrep
Проверка проекта
Проверка для тестирования
```

**Примечание:** Для вывода номеров строк используйте флаг `-n`:
```bash
./mygrep Проверка --file test_data.txt -n --nodes localhost:9001,localhost:9002,localhost:9003
```

---

### Пример 2: Поиск с регулярным выражением

**Команда:**
```bash
./mygrep '[0-9]' --file test_data.txt --nodes localhost:9001,localhost:9002,localhost:9003
```

**Результат:**
```
[2025/12/30T19:21:30.731] [application] [INFO] [coordinator.go:40,RunWithOptions] Starting search: pattern=[0-9], source=test_data.txt, servers=3, quorum=2  
[2025/12/30T19:21:30.734] [application] [INFO] [coordinator.go:91,func1] Server localhost:9001 responded successfully with 0 matches  
[2025/12/30T19:21:30.735] [application] [INFO] [coordinator.go:91,func1] Server localhost:9003 responded successfully with 2 matches  
[2025/12/30T19:21:30.735] [application] [INFO] [coordinator.go:91,func1] Server localhost:9002 responded successfully with 1 matches  
[2025/12/30T19:21:30.735] [application] [INFO] [coordinator.go:121,RunWithOptions] Quorum reached: 3/3 servers responded successfully  
1
2
3
[2025/12/30T19:21:30.735] [application] [INFO] [coordinator.go:134,RunWithOptions] Found matches: 3  
```

---

### Пример 3: Поиск с несколькими серверами

При обработке файлов координатор автоматически разделит файл на части и отправит их разным воркерам для параллельной обработки. Чем больше серверов, тем больше параллелизма.

**Команда:**
```bash
./mygrep Проверка --file test_data.txt --nodes localhost:9001,localhost:9002,localhost:9003,localhost:9004
```

Координатор разделит `test_data.txt` на 4 части и отправит их параллельно на 4 воркера. Результат будет собран при достижении кворума (3 из 4 ответов).

---

### Пример 4: Чтение из stdin

Утилита поддерживает чтение данных из стандартного ввода:

**Команда:**
```bash
cat test_data.txt | ./mygrep Проверка --nodes localhost:9001,localhost:9002,localhost:9003
```

или

```bash
echo -e "Проверка mygrep\nПроверка проекта\nПроверка для тестирования" | ./mygrep Проверка --nodes localhost:9001,localhost:9002,localhost:9003
```

**Результат:**
```
[2025/12/30T01:37:54.634] [application] [INFO] [coordinator.go:40,RunWithOptions] Starting search: pattern=Проверка, source=stdin, servers=3, quorum=2
[2025/12/30T01:37:54.640] [application] [INFO] [coordinator.go:91,func1] Server localhost:9001 responded successfully with 3 matches
[2025/12/30T01:37:54.642] [application] [INFO] [coordinator.go:91,func1] Server localhost:9002 responded successfully with 3 matches
[2025/12/30T01:37:54.644] [application] [INFO] [coordinator.go:91,func1] Server localhost:9003 responded successfully with 3 matches
[2025/12/30T01:37:54.645] [application] [INFO] [coordinator.go:121,RunWithOptions] Quorum reached: 3/3 servers responded successfully
Проверка mygrep
Проверка проекта
Проверка для тестирования
[2025/12/30T01:37:54.646] [application] [INFO] [coordinator.go:134,RunWithOptions] Found matches: 3
```

---

## Сравнительный тест с оригинальным grep

Для автоматического сравнения результатов с оригинальной утилитой grep используйте скрипт `test_comparison.sh`.

**Подготовка:**

1. Соберите проект:
```bash
make build
```

2. Запустите воркеры в разных терминалах:
```bash
# Терминал 1
./mygrep --server --addr :9001

# Терминал 2
./mygrep --server --addr :9002

# Терминал 3
./mygrep --server --addr :9003
```

3. Запустите скрипт сравнения:
```bash
chmod +x test_comparison.sh
./test_comparison.sh
```

Скрипт автоматически выполнит несколько тестов:
- Тест 1: Поиск слова "Проверка" (сравнение вывода без номеров строк)
- Тест 2: Поиск с флагом `-n` (сравнение вывода с номерами строк)
- Тест 3: Поиск с регулярным выражением

**Ручное сравнение:**

```bash
# Оригинальный grep
grep "Проверка" test_data.txt

# MyGrep (без номеров строк)
./mygrep Проверка --file test_data.txt --nodes localhost:9001,localhost:9002,localhost:9003

# С номерами строк
grep -n "Проверка" test_data.txt
./mygrep Проверка --file test_data.txt -n --nodes localhost:9001,localhost:9002,localhost:9003
