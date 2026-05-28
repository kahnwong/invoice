package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/signintech/gopdf"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed "Inter/Inter Variable/Inter.ttf"
var interFont []byte

//go:embed "Inter/Inter Hinted for Windows/Desktop/Inter-Bold.ttf"
var interBoldFont []byte

type Invoice struct {
	Id    string `json:"id" yaml:"id"`
	Title string `json:"title" yaml:"title"`

	Logo  string `json:"logo" yaml:"logo"`
	From  string `json:"from" yaml:"from"`
	To    string `json:"to" yaml:"to"`
	Date  string `json:"date" yaml:"date"`
	Due   string `json:"due" yaml:"due"`
	Valid string `json:"valid" yaml:"valid"`

	Items      []string  `json:"items" yaml:"items"`
	Quantities []int     `json:"quantities" yaml:"quantities"`
	Rates      []float64 `json:"rates" yaml:"rates"`

	Tax      float64 `json:"tax" yaml:"tax"`
	Discount float64 `json:"discount" yaml:"discount"`
	Currency string  `json:"currency" yaml:"currency"`

	Note string `json:"note" yaml:"note"`
}

func DefaultInvoice() Invoice {
	return Invoice{
		Id:         time.Now().Format("20060102"),
		Title:      "INVOICE",
		Rates:      []float64{25},
		Quantities: []int{2},
		Items:      []string{"Paper Cranes"},
		From:       "Project Folded, Inc.",
		To:         "Untitled Corporation, Inc.",
		Date:       time.Now().Format("Jan 02, 2006"),
		Due:        time.Now().AddDate(0, 0, 14).Format("Jan 02, 2006"),
		Tax:        0,
		Discount:   0,
		Currency:   "USD",
	}
}

func DefaultQuotation() Invoice {
	return Invoice{
		Id:         time.Now().Format("20060102"),
		Title:      "QUOTATION",
		Rates:      []float64{25},
		Quantities: []int{2},
		Items:      []string{"Paper Cranes"},
		From:       "Project Folded, Inc.",
		To:         "Untitled Corporation, Inc.",
		Date:       time.Now().Format("Jan 02, 2006"),
		Valid:      time.Now().AddDate(0, 0, 14).Format("Jan 02, 2006"),
		Tax:        0,
		Discount:   0,
		Currency:   "USD",
	}
}

var (
	importPath     string
	output         string
	file           = Invoice{}
	defaultInvoice = DefaultInvoice()

	quotationImportPath string
	quotationOutput     string
	quotationFile       = Invoice{}
	defaultQuotation    = DefaultQuotation()
)

func init() {
	viper.AutomaticEnv()

	// Invoice flags
	generateCmd.Flags().StringVar(&importPath, "import", "", "Imported file (.json/.yaml)")
	generateCmd.Flags().StringVar(&file.Id, "id", time.Now().Format("20060102"), "ID")
	generateCmd.Flags().StringVar(&file.Title, "title", "INVOICE", "Title")

	generateCmd.Flags().Float64SliceVarP(&file.Rates, "rate", "r", defaultInvoice.Rates, "Rates")
	generateCmd.Flags().IntSliceVarP(&file.Quantities, "quantity", "q", defaultInvoice.Quantities, "Quantities")
	generateCmd.Flags().StringSliceVarP(&file.Items, "item", "i", defaultInvoice.Items, "Items")

	generateCmd.Flags().StringVarP(&file.Logo, "logo", "l", defaultInvoice.Logo, "Company logo")
	generateCmd.Flags().StringVarP(&file.From, "from", "f", defaultInvoice.From, "Issuing company")
	generateCmd.Flags().StringVarP(&file.To, "to", "t", defaultInvoice.To, "Recipient company")
	generateCmd.Flags().StringVar(&file.Date, "date", defaultInvoice.Date, "Date")
	generateCmd.Flags().StringVar(&file.Due, "due", defaultInvoice.Due, "Payment due date")

	generateCmd.Flags().Float64Var(&file.Tax, "tax", defaultInvoice.Tax, "Tax")
	generateCmd.Flags().Float64VarP(&file.Discount, "discount", "d", defaultInvoice.Discount, "Discount")
	generateCmd.Flags().StringVarP(&file.Currency, "currency", "c", defaultInvoice.Currency, "Currency")

	generateCmd.Flags().StringVarP(&file.Note, "note", "n", "", "Note")
	generateCmd.Flags().StringVarP(&output, "output", "o", "invoice.pdf", "Output file (.pdf)")

	// Quotation flags
	generateQuotationCmd.Flags().StringVar(&quotationImportPath, "import", "", "Imported file (.json/.yaml)")
	generateQuotationCmd.Flags().StringVar(&quotationFile.Id, "id", time.Now().Format("20060102"), "ID")
	generateQuotationCmd.Flags().StringVar(&quotationFile.Title, "title", "QUOTATION", "Title")

	generateQuotationCmd.Flags().Float64SliceVarP(&quotationFile.Rates, "rate", "r", defaultQuotation.Rates, "Rates")
	generateQuotationCmd.Flags().IntSliceVarP(&quotationFile.Quantities, "quantity", "q", defaultQuotation.Quantities, "Quantities")
	generateQuotationCmd.Flags().StringSliceVarP(&quotationFile.Items, "item", "i", defaultQuotation.Items, "Items")

	generateQuotationCmd.Flags().StringVarP(&quotationFile.Logo, "logo", "l", defaultQuotation.Logo, "Company logo")
	generateQuotationCmd.Flags().StringVarP(&quotationFile.From, "from", "f", defaultQuotation.From, "Issuing company")
	generateQuotationCmd.Flags().StringVarP(&quotationFile.To, "to", "t", defaultQuotation.To, "Recipient company")
	generateQuotationCmd.Flags().StringVar(&quotationFile.Date, "date", defaultQuotation.Date, "Date")
	generateQuotationCmd.Flags().StringVar(&quotationFile.Valid, "valid", defaultQuotation.Valid, "Validity date")

	generateQuotationCmd.Flags().Float64Var(&quotationFile.Tax, "tax", defaultQuotation.Tax, "Tax")
	generateQuotationCmd.Flags().Float64VarP(&quotationFile.Discount, "discount", "d", defaultQuotation.Discount, "Discount")
	generateQuotationCmd.Flags().StringVarP(&quotationFile.Currency, "currency", "c", defaultQuotation.Currency, "Currency")

	generateQuotationCmd.Flags().StringVarP(&quotationFile.Note, "note", "n", "", "Note")
	generateQuotationCmd.Flags().StringVarP(&quotationOutput, "output", "o", "quotation.pdf", "Output file (.pdf)")

	flag.Parse()
}

var rootCmd = &cobra.Command{
	Use:   "invoice",
	Short: "Invoice generates invoices from the command line.",
	Long:  `Invoice generates invoices from the command line.`,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate an invoice",
	Long:  `Generate an invoice`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if importPath != "" {
			err := importData(importPath, &file, cmd.Flags())
			if err != nil {
				return err
			}
		}

		return generateDocument(&file, output, false)
	},
}

var generateQuotationCmd = &cobra.Command{
	Use:   "generate-quotation",
	Short: "Generate a quotation",
	Long:  `Generate a quotation`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if quotationImportPath != "" {
			err := importData(quotationImportPath, &quotationFile, cmd.Flags())
			if err != nil {
				return err
			}
		}

		return generateDocument(&quotationFile, quotationOutput, true)
	},
}

func generateDocument(doc *Invoice, outputPath string, isQuotation bool) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
	})
	pdf.SetMargins(40, 40, 40, 40)
	pdf.AddPage()
	err := pdf.AddTTFFontData("Inter", interFont)
	if err != nil {
		return err
	}

	err = pdf.AddTTFFontData("Inter-Bold", interBoldFont)
	if err != nil {
		return err
	}

	writeLogo(&pdf, doc.Logo, doc.From)
	writeTitle(&pdf, doc.Title, doc.Id, doc.Date)
	if isQuotation {
		writeBillTo(&pdf, "QUOTED TO", doc.To)
	} else {
		writeBillTo(&pdf, "BILL TO", doc.To)
	}
	writeHeaderRow(&pdf)
	subtotal := 0.0
	for i := range doc.Items {
		q := 1
		if len(doc.Quantities) > i {
			q = doc.Quantities[i]
		}

		r := 0.0
		if len(doc.Rates) > i {
			r = doc.Rates[i]
		}

		writeRow(&pdf, doc.Items[i], q, r, doc.Currency)
		subtotal += float64(q) * r
	}
	if doc.Note != "" {
		writeNotes(&pdf, doc.Note)
	}
	writeTotals(&pdf, subtotal, subtotal*doc.Tax, subtotal*doc.Discount, doc.Currency)

	var dueOrValid string
	if isQuotation {
		if doc.Valid != "" {
			dueOrValid = doc.Valid
		} else {
			dueOrValid = doc.Due
		}
	} else {
		dueOrValid = doc.Due
	}
	if dueOrValid != "" {
		if isQuotation {
			writeDueDate(&pdf, "Valid Until", dueOrValid)
		} else {
			writeDueDate(&pdf, "Due Date", dueOrValid)
		}
	}

	writeFooter(&pdf, doc.Id)
	outputPath = strings.TrimSuffix(outputPath, ".pdf") + ".pdf"
	err = pdf.WritePdf(outputPath)
	if err != nil {
		return err
	}

	fmt.Printf("Generated %s\n", outputPath)

	return nil
}

func main() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(generateQuotationCmd)
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
