#!/bin/bash

if [ ! -f "test_data.txt" ]; then
    echo "Ошибка: файл test_data.txt не найден"
    exit 1
fi

if [ ! -f "mygrep" ]; then
    echo "Ошибка: файл mygrep не найден. Сначала выполните: make build"
    exit 1
fi

echo "Тест 1: Поиск слова 'Проверка'"
echo ""
echo "Оригинальный grep:"
grep "Проверка" test_data.txt 2>/dev/null || echo "grep не найден, пропускаем"
echo ""
echo "MyGrep:"
./mygrep Проверка --file test_data.txt --nodes localhost:9001,localhost:9002,localhost:9003 2>&1 | grep -v "^\[" || echo "Ошибка: убедитесь, что воркеры запущены"
echo ""

echo "Тест 2: Поиск с номерами строк (флаг -n)"
echo ""
echo "Оригинальный grep -n:"
grep -n "Проверка" test_data.txt 2>/dev/null || echo "grep не найден, пропускаем"
echo ""
echo "MyGrep -n:"
./mygrep Проверка --file test_data.txt -n --nodes localhost:9001,localhost:9002,localhost:9003 2>&1 | grep -E "^[0-9]+:" || echo "Ошибка: убедитесь, что воркеры запущены"
echo ""

echo "Тест 3: Поиск с регулярным выражением"
echo ""
echo "Оригинальный grep:"
grep "[0-9]" test_data.txt 2>/dev/null || echo "(нет совпадений - в файле нет цифр)"
echo ""
echo "MyGrep:"
./mygrep '[0-9]' --file test_data.txt --nodes localhost:9001,localhost:9002,localhost:9003 2>&1 | grep -v "^\[" || echo "(нет совпадений)"
echo ""
