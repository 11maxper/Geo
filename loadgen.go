package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	baseURL      = "http://localhost:8080"
	username     = "loadtest_user1"
	password     = "loadtest_password"
	workers      = 3  // Количество горутин (воркеров)
	duration     = 60 // Длительность теста в секундах
	reqPerSec    = 10 // Запросов в секунду на воркер (3 воркера × 10 = 30 запросов/сек)
	numAddresses = 10 // Количество адресов
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Error string `json:"error,omitempty"`
}

type AddressSearchRequest struct {
	Query string `json:"query"`
}

type GeocodeRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

var (
	addresses = []string{
		"Москва, Ленина, 1",
		"Санкт-Петербург, Невский проспект, 2",
		"Казань, Кремлевская, 3",
		"Екатеринбург, Ленина, 4",
		"Новосибирск, Красный проспект, 5",
		"Краснодар, Красная, 6",
		"Сочи, Курортный проспект, 7",
		"Ростов-на-Дону, Большая Садовая, 8",
		"Уфа, Ленина, 9",
		"Воронеж, проспект Революции, 10",
	}

	coordinates = []struct {
		Lat string
		Lng string
	}{
		{"55.7558", "37.6173"}, // Москва
		{"59.9343", "30.3351"}, // СПб
		{"55.8304", "49.0661"}, // Казань
		{"56.8431", "60.6454"}, // Екатеринбург
		{"55.0084", "82.9357"}, // Новосибирск
		{"45.0355", "38.9753"}, // Краснодар
		{"43.6028", "39.7342"}, // Сочи
		{"47.2357", "39.7015"}, // Ростов
		{"54.7431", "55.9678"}, // Уфа
		{"51.6720", "39.1843"}, // Воронеж
	}
)

func registerUser() error {
	creds := Credentials{
		Username: username,
		Password: password,
	}

	jsonData, _ := json.Marshal(creds)
	resp, err := http.Post(baseURL+"/api/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil // Игнорируем ошибку, если пользователь уже существует
}

func login() (string, error) {
	creds := Credentials{
		Username: username,
		Password: password,
	}

	jsonData, _ := json.Marshal(creds)
	resp, err := http.Post(baseURL+"/api/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", err
	}

	if loginResp.Error != "" {
		return "", fmt.Errorf("ошибка логина: %s", loginResp.Error)
	}

	return loginResp.Token, nil
}

func addressSearch(client *http.Client, query, token string) error {
	reqBody := AddressSearchRequest{Query: query}
	jsonData, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", baseURL+"/api/address/search", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка запроса: статус %d", resp.StatusCode)
	}

	return nil
}

func geocodeSearch(client *http.Client, lat, lng, token string) error {
	reqBody := GeocodeRequest{Lat: lat, Lng: lng}
	jsonData, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", baseURL+"/api/address/geocode", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка запроса: статус %d", resp.StatusCode)
	}

	return nil
}

func worker(id int, token string, wg *sync.WaitGroup, stopChan chan struct{}) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	requestInterval := time.Second / reqPerSec
	ticker := time.NewTicker(requestInterval)
	defer ticker.Stop()

	requestCount := 0
	errorCount := 0

	for {
		select {
		case <-stopChan:
			log.Printf("Воркер %d завершен. Запросов: %d, Ошибок: %d", id, requestCount, errorCount)
			return
		case <-ticker.C:
			// Чередуем запросы поиска адресов и геокодирования
			if requestCount%2 == 0 {
				query := addresses[requestCount%len(addresses)]
				if err := addressSearch(client, query, token); err != nil {
					errorCount++
					if errorCount%10 == 0 {
						log.Printf("Воркер %d: ошибка поиска адреса (всего ошибок: %d): %v", id, errorCount, err)
					}
				}
			} else {
				coord := coordinates[requestCount%len(coordinates)]
				if err := geocodeSearch(client, coord.Lat, coord.Lng, token); err != nil {
					errorCount++
					if errorCount%10 == 0 {
						log.Printf("Воркер %d: ошибка геокодирования (всего ошибок: %d): %v", id, errorCount, err)
					}
				}
			}
			requestCount++
		}
	}
}

func main() {
	log.Println("=== Нагрузочное тестирование ===")
	log.Printf("URL: %s", baseURL)
	log.Printf("Воркеров: %d", workers)
	log.Printf("Запросов в секунду на воркер: %d", reqPerSec)
	log.Printf("Общая нагрузка: %d запросов/сек", workers*reqPerSec)
	log.Printf("Длительность: %d секунд", duration)
	log.Printf("Адресов: %d", len(addresses))
	log.Printf("Координат: %d", len(coordinates))

	// Регистрация пользователя
	log.Println("\n1. Регистрация пользователя...")
	if err := registerUser(); err != nil {
		log.Printf("Ошибка регистрации: %v", err)
	}

	// Логин
	log.Println("2. Авторизация...")
	token, err := login()
	if err != nil {
		log.Fatalf("Ошибка логина: %v", err)
	}
	log.Println("Токен получен")

	// Запускаем воркеры
	var wg sync.WaitGroup
	stopChan := make(chan struct{})

	log.Println("\n3. Запуск воркеров...")
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(i+1, token, &wg, stopChan)
	}

	// Ждем указанное время
	log.Println("4. Генерация нагрузки...")
	time.Sleep(time.Duration(duration) * time.Second)

	// Останавливаем воркеров
	log.Println("\n5. Остановка воркеров...")
	close(stopChan)
	wg.Wait()

	log.Println("\n=== Тестирование завершено ===")
}
