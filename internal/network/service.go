package network

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/DlnKot/arc/internal/config"
	"github.com/DlnKot/arc/internal/domain"
	"github.com/DlnKot/arc/internal/logging"
)

type Service struct {
	logger logging.Logger
}

func New(logger logging.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) Geo() (map[string]any, error) {
	if s.logger != nil {
		s.logger.Infof("network geo lookup started")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://ipwho.is/", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geo service returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	if success, _ := payload["success"].(bool); !success {
		return nil, errors.New("geo service error")
	}

	connection, _ := payload["connection"].(map[string]any)
	return map[string]any{
		"status":      "success",
		"country":     asString(payload["country"]),
		"countryCode": asString(payload["country_code"]),
		"regionName":  asString(payload["region"]),
		"city":        asString(payload["city"]),
		"isp":         asString(connection["isp"]),
		"org":         asString(connection["org"]),
		"query":       asString(payload["ip"]),
	}, nil
}

func (s *Service) Ping(host string, packets int) (map[string]any, error) {
	if s.logger != nil {
		s.logger.Infof("network ping started for host=%s packets=%d", host, packets)
	}
	metrics, err := runPing(strings.TrimSpace(host), packets)
	if err != nil {
		metrics.Error = err.Error()
	}
	return map[string]any{
		"ping":       metrics,
		"evaluation": evaluatePing(metrics),
	}, nil
}

func runPing(host string, packets int) (domain.PingMetrics, error) {
	if host == "" {
		return domain.PingMetrics{}, errors.New("host is required")
	}
	if packets <= 0 {
		packets = 4
	}

	args := []string{"-c", fmt.Sprintf("%d", packets), host}
	if runtime.GOOS == "windows" {
		args = []string{"-n", fmt.Sprintf("%d", packets), host}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var raw []byte
	var err error

	if runtime.GOOS == "windows" {
		cmd := exec.CommandContext(ctx, "cmd", append([]string{"/c", "chcp", "65001", ">", "NUL", "&", "ping"}, args...)...)
		applyPlatformPingAttrs(cmd)
		raw, err = cmd.CombinedOutput()
	} else {
		cmd := exec.CommandContext(ctx, "ping", args...)
		applyPlatformPingAttrs(cmd)
		raw, err = cmd.CombinedOutput()
	}

	output := string(raw)
	metrics := domain.PingMetrics{Raw: output}
	parsePingOutput(&metrics, output)
	if err != nil {
		return metrics, err
	}
	return metrics, nil
}

func parsePingOutput(metrics *domain.PingMetrics, output string) {
	lossRe := regexp.MustCompile(`([0-9]+(?:\.[0-9]+)?)%\s+packet loss`)
	statsRe := regexp.MustCompile(`=\s*([0-9]+(?:\.[0-9]+)?)/([0-9]+(?:\.[0-9]+)?)/([0-9]+(?:\.[0-9]+)?)/`)
	winLossRe := regexp.MustCompile(`\(([0-9]+)% loss\)`)
	winAvgRe := regexp.MustCompile(`Average = ([0-9]+)ms`)
	winMinRe := regexp.MustCompile(`Minimum = ([0-9]+)ms`)
	winMaxRe := regexp.MustCompile(`Maximum = ([0-9]+)ms`)

	if match := lossRe.FindStringSubmatch(output); len(match) == 2 {
		fmt.Sscanf(match[1], "%f", &metrics.LostPercent)
	}
	if match := statsRe.FindStringSubmatch(output); len(match) == 4 {
		fmt.Sscanf(match[1], "%f", &metrics.MinMs)
		fmt.Sscanf(match[2], "%f", &metrics.AvgMs)
		fmt.Sscanf(match[3], "%f", &metrics.MaxMs)
	}
	if match := winLossRe.FindStringSubmatch(output); len(match) == 2 {
		fmt.Sscanf(match[1], "%f", &metrics.LostPercent)
	}
	if match := winAvgRe.FindStringSubmatch(output); len(match) == 2 {
		fmt.Sscanf(match[1], "%f", &metrics.AvgMs)
	}
	if match := winMinRe.FindStringSubmatch(output); len(match) == 2 {
		fmt.Sscanf(match[1], "%f", &metrics.MinMs)
	}
	if match := winMaxRe.FindStringSubmatch(output); len(match) == 2 {
		fmt.Sscanf(match[1], "%f", &metrics.MaxMs)
	}
}

func evaluatePing(metrics domain.PingMetrics) domain.PingEvaluation {
	if metrics.Error != "" || metrics.LostPercent >= 100 {
		return domain.PingEvaluation{Status: "down", Label: "Недоступен", Recommendation: "Проверьте интернет-соединение и VPN, если он требуется."}
	}
	if metrics.LostPercent > 0 {
		return domain.PingEvaluation{Status: "loss", Label: "Потери пакетов", Recommendation: "Соединение нестабильно: возможны обрывы или деградация качества."}
	}
	if metrics.AvgMs > config.DefaultPingWarn {
		return domain.PingEvaluation{Status: "high_latency", Label: "Высокая задержка", Recommendation: "Соединение доступно, но может работать медленно."}
	}
	return domain.PingEvaluation{Status: "ok", Label: "Доступен", Recommendation: "Сервис отвечает стабильно."}
}

func asString(value any) string {
	if s, ok := value.(string); ok {
		return s
	}
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}
