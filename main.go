package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/net/proxy"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Links []string `yaml:"links"`
}

type Job struct {
	Index int
	URL   string
}

type Result struct {
	Index int
	URL   string
	Err   error
}

func dosyaNumarasiBul() int {
	i := 1
	for dosyaVarMi(fmt.Sprintf("screenshot_%d.png", i)) {
		i++
	}
	return i
}

func dosyaVarMi(dosyaAdi string) bool {
	_, err := os.Stat(dosyaAdi)
	return err == nil
}

func checkIp(browserCtx context.Context) {
	visitUrl := "https://check.torproject.org/"
	var torMesaj, ipText string

	err := runInNewTab(browserCtx, 90*time.Second,
		chromedp.Navigate(visitUrl),
		chromedp.Sleep(5*time.Second),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.TextContent("h1", &torMesaj, chromedp.ByQuery),
	)
	if err != nil {
		fmt.Println("[ERR] Tor kontrolünde hata:", err)
		return
	}

	_ = runInNewTab(browserCtx, 15*time.Second,
		chromedp.Navigate(visitUrl),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.TextContent("p strong", &ipText, chromedp.ByQuery),
	)

	fmt.Println("[OK] Tor mesajı:", torMesaj)
	if ipText != "" {
		fmt.Println("[OK] Görülen IP:", ipText)
	}
}

func scanReport(total int) (*os.File, error) {
	f, err := os.Create("scan_report.log")
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(f, "DARK WEB SCAN REPORT")
	fmt.Fprintf(f, "Date        : %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintln(f, "Tor Port    : 9150")
	fmt.Fprintf(f, "Total Target: %d\n\n", total)
	return f, nil
}

func writeTargetEntry(f *os.File, index int, url string) {
	fmt.Fprintf(f, "[%d] Target\n", index)
	fmt.Fprintf(f, "URL    : %s\n", url)
}

func writeScanResult(f *os.File, status string) {
	fmt.Fprintf(f, "Result : %s\n", status)
	fmt.Fprintln(f, "--------------------------------")
}

func runInNewTab(parent context.Context, d time.Duration, actions ...chromedp.Action) error {
	tabCtx, tabCancel := chromedp.NewContext(parent)
	defer tabCancel()

	runCtx, runCancel := context.WithTimeout(tabCtx, d)
	defer runCancel()

	return chromedp.Run(runCtx, actions...)
}

func ekranGoruntusuCek(browserCtx context.Context, url, dosyaAdi string) {
	var resimVerisi []byte

	err := runInNewTab(browserCtx, 120*time.Second,
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(3*time.Second),
		chromedp.FullScreenshot(&resimVerisi, 80),
	)
	if err != nil {
		fmt.Println("[ERR] Ekran görüntüsü alınamadı:", err)
		return
	}

	if err := os.WriteFile(dosyaAdi, resimVerisi, 0644); err != nil {
		fmt.Println("[ERR] Dosyaya yazılamadı:", err)
		return
	}
	fmt.Println("[OK] Ekran görüntüsü kaydedildi:", dosyaAdi)
}

func saveHTMLWithTimestamp(html []byte) (string, error) {
	filename := fmt.Sprintf(
		"%s.html",
		time.Now().Format("20060102_150405"),
	)

	if err := os.WriteFile(filename, html, 0644); err != nil {
		return "", err
	}
	return filename, nil
}

func checkTorConnection() error {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9150", nil, nil)
	if err != nil {
		return fmt.Errorf("Tor ağına bağlanılamadı :()")
	}

	conn, err := dialer.Dial("tcp", "example.com:80")
	if err != nil {
		return fmt.Errorf("Tor ağına bağlanılamadı :()")
	}
	_ = conn.Close()
	return nil
}

func main() {
	fmt.Print("Hedef dosya adı (örn: targets.yaml): ")
	var hedefDosya string
	fmt.Scanln(&hedefDosya)

	hedefDosya = strings.TrimSpace(hedefDosya)
	if hedefDosya == "" {
		fmt.Println("[WARNING] Hedef dosya adı boş olamaz.")
		return
	} else if !dosyaVarMi(hedefDosya) {
		fmt.Println("[WARNING] Hedef dosya bulunamadı.")
		return
	}

	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("proxy-server", "socks5://127.0.0.1:9150"), // 9150 tor browser
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"), //https://www.zenrows.com/blog/golang-net-http-user-agent#what-is-it
		chromedp.Flag("headless", true),
	)

	fmt.Println("Lütfen bekleyin, Tor ağına bağlanılıyor...(Port: 9150)")

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	ctx, ctxCancel := chromedp.NewContext(allocCtx)
	defer ctxCancel()

	if err := checkTorConnection(); err != nil {
		fmt.Println("[ERR]", err)
		return
	}
	fmt.Println("[OK] Tor ağıyla bağlantı kuruldu!")

	dialer, _ := proxy.SOCKS5("tcp", "127.0.0.1:9150", nil, proxy.Direct) //https://forum.golangbridge.org/t/go-http-client-and-transport/37458/3
	httpClient := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, n, a string) (net.Conn, error) {
				return dialer.Dial(n, a)
			},
		},
	}

	checkIp(ctx)

	b, err := os.ReadFile(hedefDosya)
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil { //https://pkg.go.dev/gopkg.in/yaml.v2#Marshal
		log.Fatal(err)
	}
	fmt.Printf("[INFO] YAML'dan %d link okundu.\n", len(cfg.Links))

	reportFile, err := scanReport(len(cfg.Links))
	if err != nil {
		log.Fatal(err)
	}
	defer reportFile.Close()

	jobs := make(chan Job)
	n := 2
	if len(cfg.Links) < n {
		n = len(cfg.Links)
	}

	var wg sync.WaitGroup
	for w := 0; w < n; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				dosya := fmt.Sprintf("screenshot_%d.png", j.Index)
				fmt.Printf("[INFO] (%d) %s -> %s\n", j.Index, j.URL, dosya)

				ekranGoruntusuCek(ctx, j.URL, dosya)
				status := "SUCCESS"

				resp, err := httpClient.Get(j.URL)
				if err != nil {
					fmt.Println("[ERR] HTML alınamadı:", j.URL, err)
					status = "FAILED_HTTP"
					writeScanResult(reportFile, status)
					continue
				}

				body, err := io.ReadAll(resp.Body)
				_ = resp.Body.Close()
				if err != nil {
					fmt.Println("[ERR] HTML okunamadı:", j.URL, err)
					status = "FAILED_READ"
					writeScanResult(reportFile, status)
					continue
				}

				file, err := saveHTMLWithTimestamp(body)
				if err != nil {
					fmt.Println("[ERR] HTML yazılamadı:", err)
					status = "FAILED_WRITE"
					writeScanResult(reportFile, status)
					continue
				}

				fmt.Println("[OK] HTML kaydedildi:", file)

				writeScanResult(reportFile, status)
			}

		}()
	}

	for i, url := range cfg.Links {
		url = strings.TrimSpace(url)
		if url == "" {
			continue
		}
		writeTargetEntry(reportFile, i+1, url)
		jobs <- Job{Index: i + 1, URL: url}
	}
	close(jobs)

	wg.Wait()
	fmt.Println("[OK] " + hedefDosya + " okundu. " + strconv.Itoa(len(cfg.Links)) + " forum analiz edildi.")
}
