#!/bin/bash

# Скрипт для запуска нагрузочного теста и снятия профилей

BASE_URL="http://localhost:8080"
BASE_USER="loadtest_user"
PASSWORD="loadtest_password"
# Используем первого пользователя для получения токена для профилирования
USERNAME="${BASE_USER}1"

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Нагрузочное тестирование с профилированием ===${NC}\n"

# Получаем токен
echo "Получение JWT токена..."
TOKEN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")

TOKEN=$(echo $TOKEN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "Ошибка: не удалось получить токен. Попробуем зарегистрировать пользователя..."
  curl -s -X POST "$BASE_URL/api/register" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}" > /dev/null
  
  TOKEN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")
  
  TOKEN=$(echo $TOKEN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
fi

if [ -z "$TOKEN" ]; then
  echo -e "${YELLOW}Ошибка: не удалось получить токен. Убедитесь, что сервер запущен.${NC}"
  exit 1
fi

echo -e "${GREEN}Токен получен${NC}\n"

# Создаем директорию для профилей
mkdir -p profiles

echo -e "${GREEN}=== Запуск нагрузочного теста ===${NC}"
echo "Запускаем Go-скрипт нагрузочного теста..."

# Запускаем нагрузочный тест в фоне
go run loadgen.go &
LOAD_TEST_PID=$!

# Ждем немного, чтобы нагрузка началась
sleep 5

echo -e "\n${GREEN}=== Снятие профилей (30 секунд) ===${NC}"

# CPU профиль (30 секунд)
echo "Снимаем CPU профиль (30 секунд)..."
curl -s -H "Authorization: Bearer $TOKEN" \
  "$BASE_URL/mycustompath/pprof/profile?seconds=30" > profiles/profile
echo "✓ CPU профиль сохранен: profiles/profile"

# Трейс (30 секунд)
echo "Снимаем трейс (30 секунд)..."
curl -s -H "Authorization: Bearer $TOKEN" \
  "$BASE_URL/mycustompath/pprof/trace?seconds=30" > profiles/trace
echo "✓ Трейс сохранен: profiles/trace"

# Ждем завершения нагрузочного теста
echo -e "\n${YELLOW}Ожидание завершения нагрузочного теста...${NC}"
wait $LOAD_TEST_PID

echo -e "\n${GREEN}=== Готово! ===${NC}"
echo "Профили сохранены в директории profiles/"
echo ""
echo "Для анализа профилей используйте:"
echo "  CPU профиль:"
echo "    go tool pprof profiles/profile"
echo "    или с веб-интерфейсом:"
echo "    go tool pprof -http=:8081 profiles/profile"
echo "    (откроется веб-интерфейс в браузере на http://localhost:8081)"
echo ""
echo "Трейс:"
echo "go tool trace -http=localhost:8081 profiles/trace" 
echo "    (откроется веб-интерфейс в браузере на http://localhost:8081)"
