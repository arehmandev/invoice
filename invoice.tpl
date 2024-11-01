<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        :root {
            --primary-color: #ff6b35;
            --secondary-color: #f8f9fa;
            --border-color: #dee2e6;
            --text-color: #212529;
            --table-header-bg: #f5f5f5;
        }
        @page {
            margin: 0;
            size: A4;
        }
        body {
            font-family: Arial, sans-serif;
            line-height: 1.4;
            max-width: 800px;
            margin: 0;
            padding: 0;
            color: var(--text-color);
            max-height: 1090px;
            overflow: hidden;
        }
        .brand-header {
            background-color: var(--primary-color);
            color: white;
            padding: 20px 25px;
            margin-bottom: 20px;
        }
        .brand-header h1 {
            margin: 0 0 10px 0;
            font-size: 22px;
            font-weight: 600;
        }
        .company-details {
            font-size: 13px;
            line-height: 1.5;
        }
        .content {
            padding: 0 25px;
        }
        .invoice-details, .payment-details {
            margin-bottom: 20px;
            background-color: var(--secondary-color);
            border: 1px solid var(--border-color);
            border-radius: 4px;
            padding: 15px;
        }
        .invoice-details h2, .payment-details h2 {
            margin: 0 0 15px 0;
            font-size: 20px;
            color: var(--primary-color);
        }
        .invoice-details p, .payment-details p {
            margin: 0 0 15px 0;
            line-height: 1.6;
        }
        .bill-to {
            margin-bottom: 15px;
        }
        .bill-to strong {
            display: block;
            margin-bottom: 5px;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
            font-size: 13px;
        }
        th, td {
            border: 1px solid var(--border-color);
            padding: 8px 12px;
        }
        th {
            background-color: var(--table-header-bg);
            font-weight: 600;
        }
        .amount-column {
            text-align: right;
            width: 120px;
        }
        .quantity-column {
            text-align: center;
            width: 100px;
        }
        .description-column {
            width: 50%;
        }
        .totals-section td {
            font-weight: 600;
        }
        .payment-terms {
            margin-top: 15px;
            padding: 12px;
            background-color: var(--secondary-color);
            border: 1px solid var(--border-color);
            border-radius: 4px;
            font-size: 12px;
            line-height: 1.5;
        }
        .reference-line {
            margin-top: 10px !important;
            font-size: 12px !important;
            line-height: 1.2 !important;
        }
        .bank-info {
            margin: 0;
            line-height: 1.6;
            font-size: 13px;
        }
    </style>
</head>
<body>
    <div class="brand-header">
        <h1>{{.FromCompany.Name}}</h1>
        <div class="company-details">
            {{.FromCompany.Street}}<br>
            {{.FromCompany.City}}<br>
            {{.FromCompany.PostCode}}
            {{if .FromCompany.VatNumber}}
            <br>VAT Registration No.: {{.FromCompany.VatNumber}}
            {{end}}
        </div>
    </div>

    <div class="content">
        <div class="invoice-details">
            <h2>INVOICE</h2>
            <div class="bill-to">
                <strong>Bill To:</strong>
                {{.ToCompany.Name}}<br>
                {{.ToCompany.Street}}<br>
                {{.ToCompany.City}}<br>
                {{.ToCompany.PostCode}}
            </div>

            <p>
                <strong>Invoice Number:</strong> {{.Number}}<br>
                <strong>Invoice Date:</strong> {{.Date.Format "02/01/2006"}}<br>
                <strong>PO Reference:</strong> {{.PONumber}}<br>
                <strong>Due Date:</strong> {{.DueDate.Format "02/01/2006"}}
            </p>
        </div>

        <table>
            <thead>
                <tr>
                    <th class="description-column">Description</th>
                    <th class="quantity-column">Quantity</th>
                    <th class="amount-column">Rate</th>
                    <th class="amount-column">Amount</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>Consultancy fees for services provided in the month of {{.Date.Format "January 2006"}}</td>
                    <td class="quantity-column">{{.Days}}</td>
                    <td class="amount-column">£{{printf "%.2f" .DailyRate}}</td>
                    <td class="amount-column">£{{printf "%.2f" .Amount}}</td>
                </tr>
                <tr class="totals-section">
                    <td colspan="3" style="text-align: right;">Subtotal</td>
                    <td class="amount-column">£{{printf "%.2f" .Amount}}</td>
                </tr>
                <tr class="totals-section">
                    <td colspan="3" style="text-align: right;">VAT @ 20%</td>
                    <td class="amount-column">£{{printf "%.2f" .VAT}}</td>
                </tr>
                <tr class="totals-section">
                    <td colspan="3" style="text-align: right;">Total</td>
                    <td class="amount-column">£{{printf "%.2f" .TotalAmount}}</td>
                </tr>
            </tbody>
        </table>

        <div class="payment-details">
            <h2>Payment Details</h2>
            <div class="bank-info">
                <strong>Bank:</strong> {{.Bank.BankName}}<br>
                <strong>Account Name:</strong> {{.Bank.AccountName}}<br>
                <strong>Sort Code:</strong> {{.Bank.SortCode}}<br>
                <strong>Account Number:</strong> {{.Bank.AccountNumber}}
            </div>

            <p class="reference-line">
                Please quote invoice number {{.Number}} on all payments and correspondence.
            </p>
        </div>

        <div class="payment-terms">
            <strong>Payment Terms:</strong> In consideration of the provision of the Services, the Client shall pay each invoice submitted by the Consultant Company in accordance with clause 4.1, within {{.PaymentTerms}} days of receipt.
        </div>
    </div>
</body>
</html>