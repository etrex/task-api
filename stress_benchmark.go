package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogolook/task-api/model"
	"github.com/gogolook/task-api/storage"
)

// 壓力測試程式

// 監控工具函數
func getCPUUsage() float64 {
	cmd := exec.Command("ps", "-p", strconv.Itoa(os.Getpid()), "-o", "pcpu")
	output, err := cmd.Output()
	if err != nil {
		return 0.0
	}
	
	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return 0.0
	}
	
	cpu, err := strconv.ParseFloat(strings.TrimSpace(lines[1]), 64)
	if err != nil {
		return 0.0
	}
	return cpu
}

func getConnectionCount() int {
	cmd := exec.Command("netstat", "-an")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}
	
	count := 0
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, ":8080") && strings.Contains(line, "ESTABLISHED") {
			count++
		}
	}
	return count
}

func getDetailedMemStats() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return fmt.Sprintf("Alloc=%dKB, Sys=%dKB, GC=%d, Goroutines=%d", 
		m.Alloc/1024, m.Sys/1024, m.NumGC, runtime.NumGoroutine())
}

// 可調整的並發參數
var (
	storageGoroutines  = 100000
	httpConcurrency    = 1000
	mixedReaders       = 70
	mixedWriters       = 20
	mixedListers       = 10
	longRunningWorkers = 50
	initialTaskCount   = 1000
	mixedInitialTasks  = 1000
)

type TestResult struct {
	TotalRequests    int64
	SuccessRequests  int64
	FailedRequests   int64
	TimeoutRequests  int64
	AverageTime      time.Duration
	MaxTime          time.Duration
	MinTime          time.Duration
}

func main() {
	fmt.Println("=== Task API 壓力測試 ===\n")
	
	// 1. 直接測試 Storage 層
	fmt.Println("1. Storage 層壓力測試")
	testStorageStress()
	
	// 2. 測試 HTTP API
	fmt.Println("\n2. HTTP API 壓力測試")
	testHTTPStress()
	
	// 2.1 HTTP API 漸進式壓力測試
	fmt.Println("\n2.1 HTTP API 漸進式壓力測試")
	testHTTPProgressiveStress()
	
	// 3. 混合讀寫測試
	fmt.Println("\n3. 混合讀寫壓力測試")
	testMixedStress()
	
	// 4. 長時間壓力測試 (已停用)
	// fmt.Println("\n4. 長時間壓力測試")
	// testLongRunningStress()
}

func testStorageStress() {
	fmt.Printf("測試場景：%d 個 goroutine 同時對 Storage 進行讀寫\n", storageGoroutines)
	
	memStorage := storage.NewMemoryStorage()
	
	// 準備一些初始資料
	taskIDs := make([]string, initialTaskCount)
	for i := 0; i < initialTaskCount; i++ {
		task := &model.Task{
			Name:   fmt.Sprintf("task-%d", i),
			Status: i % 2,
		}
		memStorage.Create(task)
		taskIDs[i] = task.ID
	}
	
	var result TestResult
	var wg sync.WaitGroup
	
	startTime := time.Now()
	
	// 啟動 Storage 測試的 goroutine
	for i := 0; i < storageGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			operationStart := time.Now()
			
			// 隨機進行不同操作
			switch id % 5 {
			case 0, 1, 2: // 60% 讀取操作
				_, err := memStorage.Get(taskIDs[id%len(taskIDs)])
				if err != nil {
					atomic.AddInt64(&result.FailedRequests, 1)
				} else {
					atomic.AddInt64(&result.SuccessRequests, 1)
				}
				
			case 3: // 20% 創建操作
				task := &model.Task{
					Name:   fmt.Sprintf("new-task-%d", id),
					Status: 0,
				}
				err := memStorage.Create(task)
				if err != nil {
					atomic.AddInt64(&result.FailedRequests, 1)
				} else {
					atomic.AddInt64(&result.SuccessRequests, 1)
				}
				
			case 4: // 20% 列表操作
				tasks := memStorage.List()
				if len(tasks) > 0 {
					atomic.AddInt64(&result.SuccessRequests, 1)
				} else {
					atomic.AddInt64(&result.FailedRequests, 1)
				}
			}
			
			operationTime := time.Since(operationStart)
			atomic.AddInt64(&result.TotalRequests, 1)
			
			// 更新時間統計（簡化版）
			if operationTime > result.MaxTime {
				result.MaxTime = operationTime
			}
			if result.MinTime == 0 || operationTime < result.MinTime {
				result.MinTime = operationTime
			}
		}(i)
	}
	
	// 設定超時檢測
	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()
	
	select {
	case <-done:
		totalTime := time.Since(startTime)
		result.AverageTime = totalTime / time.Duration(result.TotalRequests)
		
		fmt.Printf("✅ Storage 測試完成:\n")
		fmt.Printf("   總請求: %d\n", result.TotalRequests)
		fmt.Printf("   成功: %d\n", result.SuccessRequests)
		fmt.Printf("   失敗: %d\n", result.FailedRequests)
		fmt.Printf("   總時間: %v\n", totalTime)
		fmt.Printf("   平均時間: %v\n", result.AverageTime)
		fmt.Printf("   最大時間: %v\n", result.MaxTime)
		fmt.Printf("   最小時間: %v\n", result.MinTime)
		
	case <-time.After(30 * time.Second):
		fmt.Printf("❌ Storage 測試超時 (30秒)!\n")
		fmt.Printf("   這表示沒有超時保護可能導致卡死\n")
		fmt.Printf("   已完成請求: %d/%d\n", result.SuccessRequests+result.FailedRequests, result.TotalRequests)
		result.TimeoutRequests = result.TotalRequests - result.SuccessRequests - result.FailedRequests
	}
}

func testHTTPStress() {
	fmt.Printf("測試場景：%d 個並發 HTTP 請求\n", httpConcurrency)
	
	// 檢查 API 是否運行
	baseURL := "http://localhost:8080"
	testClient := &http.Client{Timeout: 1 * time.Second}
	_, err := testClient.Get(baseURL + "/health")
	if err != nil {
		fmt.Printf("❌ API 服務器未啟動 (%s)，跳過 HTTP 測試\n", err.Error())
		return
	}
	
	var result TestResult
	var wg sync.WaitGroup
	
	// 設定 HTTP 客戶端with連接池
	transport := &http.Transport{
		MaxIdleConns:        100,               // 最大空閒連接數
		MaxIdleConnsPerHost: 100,               // 每個主機最大空閒連接數
		MaxConnsPerHost:     200,               // 每個主機最大連接數
		IdleConnTimeout:     30 * time.Second,  // 空閒連接超時
	}
	client := &http.Client{
		Timeout:   10 * time.Second, // 10秒超時
		Transport: transport,
	}
	
	startTime := time.Now()
	
	// 啟動 HTTP 並發請求
	for i := 0; i < httpConcurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			requestStart := time.Now()
			atomic.AddInt64(&result.TotalRequests, 1)
			
			// 隨機請求類型
			var resp *http.Response
			var err error
			
			switch id % 4 {
			case 0, 1: // 50% GET 請求
				resp, err = client.Get(baseURL + "/tasks")
				
			case 2: // 25% POST 請求
				task := model.Task{
					Name:   fmt.Sprintf("http-task-%d", id),
					Status: 0,
				}
				jsonData, _ := json.Marshal(task)
				resp, err = client.Post(baseURL+"/tasks", "application/json", bytes.NewBuffer(jsonData))
				
			case 3: // 25% GET 單個任務
				resp, err = client.Get(baseURL + "/tasks/test-id")
			}
			
			if err != nil {
				atomic.AddInt64(&result.FailedRequests, 1)
				if err.Error() == "timeout" {
					atomic.AddInt64(&result.TimeoutRequests, 1)
				}
				// 記錄錯誤類型以便調試
				if id < 10 { // 只記錄前10個錯誤避免spam
					fmt.Printf("   錯誤 %d: %s\n", id, err.Error())
				}
			} else {
				resp.Body.Close()
				if resp.StatusCode < 500 {
					atomic.AddInt64(&result.SuccessRequests, 1)
				} else {
					atomic.AddInt64(&result.FailedRequests, 1)
					// 記錄HTTP錯誤狀態碼
					if id < 10 {
						fmt.Printf("   HTTP錯誤 %d: 狀態碼 %d\n", id, resp.StatusCode)
					}
				}
			}
			
			requestTime := time.Since(requestStart)
			if requestTime > result.MaxTime {
				result.MaxTime = requestTime
			}
			if result.MinTime == 0 || requestTime < result.MinTime {
				result.MinTime = requestTime
			}
		}(i)
	}
	
	// 等待所有請求完成或超時
	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()
	
	select {
	case <-done:
		totalTime := time.Since(startTime)
		result.AverageTime = totalTime / time.Duration(result.TotalRequests)
		
		fmt.Printf("✅ HTTP 測試完成:\n")
		fmt.Printf("   總請求: %d\n", result.TotalRequests)
		fmt.Printf("   成功: %d\n", result.SuccessRequests)
		fmt.Printf("   失敗: %d\n", result.FailedRequests)
		fmt.Printf("   超時: %d\n", result.TimeoutRequests)
		fmt.Printf("   總時間: %v\n", totalTime)
		fmt.Printf("   平均時間: %v\n", result.AverageTime)
		fmt.Printf("   最大時間: %v\n", result.MaxTime)
		fmt.Printf("   最小時間: %v\n", result.MinTime)
		
	case <-time.After(60 * time.Second):
		fmt.Printf("❌ HTTP 測試超時 (60秒)!\n")
		fmt.Printf("   這表示 API 可能在高並發下卡死\n")
		fmt.Printf("   已完成請求: %d/%d\n", result.SuccessRequests+result.FailedRequests, result.TotalRequests)
	}
}

func testHTTPProgressiveStress() {
	fmt.Println("測試場景：逐步增加並發數直到通過率低於 95%")
	
	// 檢查 API 是否運行
	baseURL := "http://localhost:8080"
	testClient := &http.Client{Timeout: 1 * time.Second}
	_, err := testClient.Get(baseURL + "/health")
	if err != nil {
		fmt.Printf("❌ API 服務器未啟動 (%s)，跳過漸進式測試\n", err.Error())
		return
	}
	
	// 顯示系統資源信息
	fmt.Printf("系統資源檢查：\n")
	fmt.Printf("  GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("  NumGoroutine: %d\n", runtime.NumGoroutine())
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("  Memory: Alloc=%d KB, Sys=%d KB\n", m.Alloc/1024, m.Sys/1024)
	
	// 設定 HTTP 客戶端
	transport := &http.Transport{
		MaxIdleConns:        500,
		MaxIdleConnsPerHost: 500,
		MaxConnsPerHost:     1000,
		IdleConnTimeout:     30 * time.Second,
	}
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}
	
	fmt.Println("開始漸進式測試...")
	fmt.Printf("%-8s %-8s %-8s %-8s %-10s %-12s\n", "並發數", "總請求", "成功", "失敗", "成功率", "平均時間")
	fmt.Println("---------------------------------------------------------------")
	
	// 從 100 開始，每次增加 100，直到成功率低於 95%
	for concurrency := 100; concurrency <= 5000; concurrency += 100 {
		result := runHTTPTest(client, baseURL, concurrency)
		
		successRate := float64(result.SuccessRequests) / float64(result.TotalRequests) * 100
		avgTime := result.AverageTime
		
		// 顯示詳細資源使用情況
		cpuUsage := getCPUUsage()
		connCount := getConnectionCount()
		memStats := getDetailedMemStats()
		
		resourceInfo := fmt.Sprintf(" [CPU: %.1f%%, Conn: %d, %s]", 
			cpuUsage, connCount, memStats)
		
		fmt.Printf("%-8d %-8d %-8d %-8d %-10.2f%% %-12v%s\n", 
			concurrency, result.TotalRequests, result.SuccessRequests, 
			result.FailedRequests, successRate, avgTime, resourceInfo)
		
		// 在關鍵區域進行更詳細的監控
		fmt.Printf("  詳細監控: ")
		
		// 多次取樣 CPU 使用率
		cpuSamples := make([]float64, 3)
		for i := 0; i < 3; i++ {
			cpuSamples[i] = getCPUUsage()
			time.Sleep(100 * time.Millisecond)
		}
		
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		runtime.GC() // 手動觸發 GC 看看記憶體回收情況
		runtime.ReadMemStats(&m)
		
		fmt.Printf("CPU採樣: %.1f/%.1f/%.1f%%, ", cpuSamples[0], cpuSamples[1], cpuSamples[2])
		fmt.Printf("GC: %d次, PauseTotalNs: %d, HeapInuse: %dKB\n", 
			m.NumGC, m.PauseTotalNs/1000000, m.HeapInuse/1024)
		
		// 如果成功率低於 95%，停止測試
		if successRate < 95.0 {
			fmt.Printf("\n🔥 發現性能瓶頸！\n")
			fmt.Printf("   最大穩定並發數: %d\n", concurrency-10)
			fmt.Printf("   在 %d 並發時成功率降至 %.2f%%\n", concurrency, successRate)
			break
		}
		
		// 如果成功率 100%，可以加大增長步長
		if successRate == 100.0 && concurrency >= 3100 {
			concurrency += 50 // 加大步長，快速找到瓶頸
		}
		
		// 短暫休息讓服務器恢復
		time.Sleep(100 * time.Millisecond)
	}
}

func runHTTPTest(client *http.Client, baseURL string, concurrency int) TestResult {
	var result TestResult
	var wg sync.WaitGroup
	
	startTime := time.Now()
	
	// 啟動指定數量的並發請求
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			requestStart := time.Now()
			atomic.AddInt64(&result.TotalRequests, 1)
			
			// 隨機請求類型
			var resp *http.Response
			var err error
			
			switch id % 3 {
			case 0: // 50% GET 請求
				resp, err = client.Get(baseURL + "/tasks")
			case 1: // 33% POST 請求
				task := model.Task{
					Name:   fmt.Sprintf("test-task-%d", id),
					Status: 0,
				}
				jsonData, _ := json.Marshal(task)
				resp, err = client.Post(baseURL+"/tasks", "application/json", bytes.NewBuffer(jsonData))
			case 2: // 17% GET 單個任務
				resp, err = client.Get(baseURL + "/tasks/test-id")
			}
			
			if err != nil {
				atomic.AddInt64(&result.FailedRequests, 1)
				if err.Error() == "timeout" {
					atomic.AddInt64(&result.TimeoutRequests, 1)
				}
				// 在高並發時記錄錯誤
				if concurrency >= 200 && atomic.LoadInt64(&result.FailedRequests) <= 5 {
					fmt.Printf("   [並發%d] 錯誤: %s\n", concurrency, err.Error())
				}
			} else {
				resp.Body.Close()
				if resp.StatusCode < 500 {
					atomic.AddInt64(&result.SuccessRequests, 1)
				} else {
					atomic.AddInt64(&result.FailedRequests, 1)
					// 記錄 HTTP 錯誤狀態碼
					if concurrency >= 200 && atomic.LoadInt64(&result.FailedRequests) <= 5 {
						fmt.Printf("   [並發%d] HTTP錯誤: 狀態碼 %d\n", concurrency, resp.StatusCode)
					}
				}
			}
			
			requestTime := time.Since(requestStart)
			if requestTime > result.MaxTime {
				result.MaxTime = requestTime
			}
			if result.MinTime == 0 || requestTime < result.MinTime {
				result.MinTime = requestTime
			}
		}(i)
	}
	
	// 等待所有請求完成
	wg.Wait()
	
	totalTime := time.Since(startTime)
	result.AverageTime = totalTime / time.Duration(result.TotalRequests)
	
	return result
}

func testMixedStress() {
	fmt.Println("測試場景：混合讀寫操作，模擬真實使用情境")
	
	memStorage := storage.NewMemoryStorage()
	
	// 初始化一些資料
	taskIDs := make([]string, mixedInitialTasks)
	for i := 0; i < mixedInitialTasks; i++ {
		task := &model.Task{
			Name:   fmt.Sprintf("initial-task-%d", i),
			Status: 0,
		}
		memStorage.Create(task)
		taskIDs[i] = task.ID
	}
	
	var operations int64
	var errors int64
	
	// 啟動多個類型的工作者
	var wg sync.WaitGroup
	
	// 讀取工作者
	for i := 0; i < mixedReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_, err := memStorage.Get(taskIDs[j%len(taskIDs)])
				atomic.AddInt64(&operations, 1)
				if err != nil {
					atomic.AddInt64(&errors, 1)
				}
				time.Sleep(time.Millisecond) // 模擬真實間隔
			}
		}()
	}
	
	// 寫入工作者
	for i := 0; i < mixedWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				task := &model.Task{
					Name:   fmt.Sprintf("writer-%d-task-%d", id, j),
					Status: 0,
				}
				err := memStorage.Create(task)
				atomic.AddInt64(&operations, 1)
				if err != nil {
					atomic.AddInt64(&errors, 1)
				}
				time.Sleep(2 * time.Millisecond) // 寫入稍慢
			}
		}(i)
	}
	
	// 列表工作者
	for i := 0; i < mixedListers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				tasks := memStorage.List()
				atomic.AddInt64(&operations, 1)
				if len(tasks) == 0 {
					atomic.AddInt64(&errors, 1)
				}
				time.Sleep(5 * time.Millisecond) // 列表操作更慢
			}
		}()
	}
	
	// 監控進度
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				ops := atomic.LoadInt64(&operations)
				errs := atomic.LoadInt64(&errors)
				fmt.Printf("   進度: %d 操作完成, %d 錯誤\n", ops, errs)
			}
		}
	}()
	
	startTime := time.Now()
	
	// 等待完成或超時
	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()
	
	select {
	case <-done:
		totalTime := time.Since(startTime)
		finalOps := atomic.LoadInt64(&operations)
		finalErrs := atomic.LoadInt64(&errors)
		
		fmt.Printf("✅ 混合測試完成:\n")
		fmt.Printf("   總操作: %d\n", finalOps)
		fmt.Printf("   錯誤: %d\n", finalErrs)
		fmt.Printf("   成功率: %.2f%%\n", float64(finalOps-finalErrs)/float64(finalOps)*100)
		fmt.Printf("   總時間: %v\n", totalTime)
		fmt.Printf("   吞吐量: %.2f ops/sec\n", float64(finalOps)/totalTime.Seconds())
		
	case <-time.After(45 * time.Second):
		finalOps := atomic.LoadInt64(&operations)
		finalErrs := atomic.LoadInt64(&errors)
		
		fmt.Printf("❌ 混合測試超時 (45秒)!\n")
		fmt.Printf("   完成操作: %d\n", finalOps)
		fmt.Printf("   錯誤: %d\n", finalErrs)
		fmt.Printf("   這表示在混合負載下可能出現性能問題\n")
	}
}

func testLongRunningStress() {
	fmt.Println("測試場景：長時間運行測試 (2分鐘)")
	
	memStorage := storage.NewMemoryStorage()
	
	var totalOps int64
	var totalErrs int64
	
	// 持續運行的工作者
	stopChan := make(chan bool)
	var wg sync.WaitGroup
	
	// 啟動長時間測試的 goroutine
	for i := 0; i < longRunningWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			for {
				select {
				case <-stopChan:
					return
				default:
					// 執行操作
					task := &model.Task{
						Name:   fmt.Sprintf("long-task-%d-%d", id, time.Now().UnixNano()),
						Status: 0,
					}
					
					err := memStorage.Create(task)
					atomic.AddInt64(&totalOps, 1)
					if err != nil {
						atomic.AddInt64(&totalErrs, 1)
					}
					
					// 隨機延遲
					time.Sleep(time.Duration(id%10) * time.Millisecond)
				}
			}
		}(i)
	}
	
	// 監控統計
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				ops := atomic.LoadInt64(&totalOps)
				errs := atomic.LoadInt64(&totalErrs)
				fmt.Printf("   長時間測試進度: %d 操作, %d 錯誤\n", ops, errs)
			case <-stopChan:
				return
			}
		}
	}()
	
	// 運行 2 分鐘
	time.Sleep(2 * time.Minute)
	
	// 停止所有工作者
	close(stopChan)
	wg.Wait()
	
	finalOps := atomic.LoadInt64(&totalOps)
	finalErrs := atomic.LoadInt64(&totalErrs)
	
	fmt.Printf("✅ 長時間測試完成:\n")
	fmt.Printf("   總操作: %d\n", finalOps)
	fmt.Printf("   錯誤: %d\n", finalErrs)
	fmt.Printf("   成功率: %.2f%%\n", float64(finalOps-finalErrs)/float64(finalOps)*100)
	fmt.Printf("   平均 QPS: %.2f\n", float64(finalOps)/120.0) // 2分鐘 = 120秒
	
	if finalErrs > 0 {
		fmt.Printf("⚠️  發現 %d 個錯誤，建議檢查並發安全性\n", finalErrs)
	}
}