package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"gopkg.in/yaml.v3"
)

type Config struct {
	FromCompany struct {
		Name      string `yaml:"name"`
		Street    string `yaml:"street"`
		City      string `yaml:"city"`
		Postcode  string `yaml:"postcode"`
		VatNumber string `yaml:"vat_number"`
	} `yaml:"from_company"`

	ToCompany struct {
		Name      string `yaml:"name"`
		Street    string `yaml:"street"`
		City      string `yaml:"city"`
		Postcode  string `yaml:"postcode"`
		VatNumber string `yaml:"vat_number"`
	} `yaml:"to_company"`

	Bank struct {
		Name          string `yaml:"name"`
		AccountName   string `yaml:"account_name"`
		SortCode      string `yaml:"sort_code"`
		AccountNumber string `yaml:"account_number"`
	} `yaml:"bank"`

	Business struct {
		VatRate         float64 `yaml:"vat_rate"`
		DailyRate       float64 `yaml:"daily_rate"`
		PaymentTermDays int     `yaml:"payment_terms_days"`
	} `yaml:"business"`

	Style struct {
		PrimaryColor string `yaml:"primary_color"`
	} `yaml:"style"`
}

type Address struct {
	Name      string
	Street    string
	City      string
	PostCode  string
	VatNumber string
}

type BankDetails struct {
	BankName      string
	AccountName   string
	SortCode      string
	AccountNumber string
}

type Invoice struct {
	Number       string
	Date         time.Time
	PONumber     string
	DueDate      time.Time
	Days         int
	DailyRate    float64
	Amount       float64
	VAT          float64
	TotalAmount  float64
	PaymentTerms int
	FromCompany  Address
	ToCompany    Address
	Bank         BankDetails
	Style        struct{ PrimaryColor string }
}

func loadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "config.yaml"
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func parseCustomDate(dateStr string) (time.Time, error) {
	return time.Parse("02-01-06", dateStr)
}

func isLastWeekOfMonth() bool {
	now := time.Now()
	lastDay := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.UTC)
	daysUntilEnd := lastDay.Sub(now).Hours() / 24

	return daysUntilEnd <= 7
}

func determineInvoiceDate() time.Time {
	now := time.Now()

	if isLastWeekOfMonth() {
		return time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.UTC)
	}

	return time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, time.UTC)
}

func generateInvoiceNumber(date time.Time) string {
	return fmt.Sprintf("INV-%d-%02d", date.Year(), date.Month())
}

func generateFileName(invoiceNumber string) string {
	safeName := strings.ReplaceAll(invoiceNumber, "/", "-")
	return fmt.Sprintf("%s.pdf", safeName)
}

func ensureOutputDir(dir string) error {
	if dir != "" && dir != "." {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func createPDF(html string, outputPath string) error {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-software-rasterizer", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	encodedHTML := base64.StdEncoding.EncodeToString([]byte(html))

	var pdfBuf []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate("data:text/html;base64,"+encodedHTML),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(1*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithMarginTop(0.2).
				WithMarginBottom(0.2).
				WithMarginLeft(0.2).
				WithMarginRight(0.2).
				WithPaperWidth(8.27).
				WithPaperHeight(11.7).
				WithPreferCSSPageSize(true).
				Do(ctx)
			if err != nil {
				return err
			}
			pdfBuf = buf
			return nil
		}),
	); err != nil {
		return fmt.Errorf("could not create PDF: %v", err)
	}

	return os.WriteFile(outputPath, pdfBuf, 0644)
}

func main() {
	// CLI flags
	days := flag.Int("days", 0, "Number of days worked")
	poNumber := flag.String("po", "", "Purchase Order reference number")
	outputDir := flag.String("outdir", ".", "Output directory for invoice files")
	configPath := flag.String("config", "", "Path to config file (default: config.yaml)")
	dateOverride := flag.String("date", "", "Override invoice date (DD-MM-YY format)")
	flag.Parse()

	// Load configuration
	config, err := loadConfig(*configPath)
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Validate required inputs
	if *days == 0 || *poNumber == "" {
		fmt.Println("Usage: invoice -days N -po NUMBER [-outdir DIR] [-config PATH] [-date DD-MM-YY]")
		fmt.Println("\nFlags:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Determine invoice date
	var date time.Time
	if *dateOverride != "" {
		date, err = parseCustomDate(*dateOverride)
		if err != nil {
			fmt.Printf("Error parsing date: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Using override date: %s\n", date.Format("January 2006"))
	} else {
		date = determineInvoiceDate()
		if isLastWeekOfMonth() {
			fmt.Printf("Using end of current month (%s) for invoice\n", date.Format("January 2006"))
		} else {
			fmt.Printf("Using end of previous month (%s) for invoice\n", date.Format("January 2006"))
		}
	}

	// Calculate amounts
	amount := float64(*days) * config.Business.DailyRate
	vat := amount * config.Business.VatRate
	total := amount + vat

	// Calculate due date
	dueDate := date.AddDate(0, 0, config.Business.PaymentTermDays)

	// Generate invoice number
	invoiceNumber := generateInvoiceNumber(date)

	// Create invoice data
	inv := Invoice{
		Number:       invoiceNumber,
		Date:         date,
		PONumber:     *poNumber,
		DueDate:      dueDate,
		Days:         *days,
		DailyRate:    config.Business.DailyRate,
		Amount:       amount,
		VAT:          vat,
		TotalAmount:  total,
		PaymentTerms: config.Business.PaymentTermDays,
		FromCompany: Address{
			Name:      config.FromCompany.Name,
			Street:    config.FromCompany.Street,
			City:      config.FromCompany.City,
			PostCode:  config.FromCompany.Postcode,
			VatNumber: config.FromCompany.VatNumber,
		},
		ToCompany: Address{
			Name:      config.ToCompany.Name,
			Street:    config.ToCompany.Street,
			City:      config.ToCompany.City,
			PostCode:  config.ToCompany.Postcode,
			VatNumber: config.ToCompany.VatNumber,
		},
		Bank: BankDetails{
			BankName:      config.Bank.Name,
			AccountName:   config.Bank.AccountName,
			SortCode:      config.Bank.SortCode,
			AccountNumber: config.Bank.AccountNumber,
		},
		Style: struct{ PrimaryColor string }{
			PrimaryColor: config.Style.PrimaryColor,
		},
	}

	// Ensure output directory exists
	if err := ensureOutputDir(*outputDir); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Generate output filename
	fileName := generateFileName(invoiceNumber)
	filePath := filepath.Join(*outputDir, fileName)

	if *outputDir == "." {
		filePath = fileName
	}

	// Create and execute template
	invoiceTmpl, err := os.ReadFile("invoice.tpl")
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New("invoice").Parse(string(invoiceTmpl))
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, inv)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		os.Exit(1)
	}

	// Generate PDF
	err = createPDF(buf.String(), filePath)
	if err != nil {
		fmt.Printf("Error generating PDF: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Invoice PDF generated successfully: %s\n", filePath)
}
