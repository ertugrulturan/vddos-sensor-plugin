package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"time"
)

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func checkForAttack(activeConnections int) {
	if activeConnections > 5 {
		fmt.Println("\nTHERE IS AN ATTACK\n")

		confFile := "/vddos/conf.d/website.conf"
		content, err := ioutil.ReadFile(confFile)
		if err != nil {
			panic(err)
		}

		newContent := regexp.MustCompile(`high\+`).ReplaceAllString(string(content), "captcha+")

		err = ioutil.WriteFile(confFile, []byte(newContent), 0644)
		if err != nil {
			panic(err)
		}

		runLwGuardian()

		countdown(1 * time.Minute)

		newContent = regexp.MustCompile(`captcha\+`).ReplaceAllString(string(content), "high+")

		err = ioutil.WriteFile(confFile, []byte(newContent), 0644)
		if err != nil {
			panic(err)
		}

		runLwGuardian()

		// Geri sayım bittikten sonra check işlemini yeniden başlat
		main()
	}
}

func countdown(duration time.Duration) {
	endTime := time.Now().Add(duration)

	for time.Now().Before(endTime) {
		clearScreen()
		fmt.Printf("Geri sayım: %v\n", endTime.Sub(time.Now()))
		time.Sleep(time.Second)
	}
}

func runLwGuardian() {
	cmd := exec.Command("/usr/bin/lwguardian", ">", "/dev/null", "2>/dev/null", "&")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	monitorFile := "/vddos/Monitor_GUI.txt"

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var lastRequests int
	var lastTime time.Time

	for {
		if _, err := os.Stat(monitorFile); os.IsNotExist(err) {
			fmt.Println("Dosya bulunamadı:", monitorFile)
			return
		}

		content, err := ioutil.ReadFile(monitorFile)
		if err != nil {
			panic(err)
		}

		urlRegex := regexp.MustCompile(`http://\d+\.\d+\.\d+\.\d+/[a-zA-Z0-9_-]*-nginx-status`)
		url := urlRegex.FindString(string(content))

		if url == "" {
			fmt.Println("Nginx status URL'i bulunamadı")
			return
		}

		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		re := regexp.MustCompile(`Active connections: (\d+)`)
		matches := re.FindStringSubmatch(string(body))
		if len(matches) < 2 {
			panic("Could not find active connections in nginx status page")
		}

		activeConnections, err := strconv.Atoi(matches[1])
		if err != nil {
			panic(err)
		}

		re = regexp.MustCompile(`\s*(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+`)
		matches = re.FindStringSubmatch(string(body))
		if len(matches) < 5 {
			panic("Could not find server statistics in nginx status page")
		}

		requests, err := strconv.Atoi(matches[3])
		if err != nil {
			panic(err)
		}

		if !lastTime.IsZero() {
			duration := time.Since(lastTime)
			requestsDiff := requests - lastRequests

			reqPerSec := float64(requestsDiff) / duration.Seconds()

			clearScreen()

			fmt.Printf("Requests per second: %.0f\n", reqPerSec)
			fmt.Printf("Active connections: %d\n", activeConnections)

			if reqPerSec > 300 {
				checkForAttack(activeConnections)
			}
		}

		lastTime = time.Now()
		lastRequests = requests

		time.Sleep(1 * time.Second)
	}
}
