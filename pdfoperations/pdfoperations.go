package pdfoperations

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/Abhi-singh-karuna/my_Liberary/baselogger"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/sighandler"
)

// Pdf requestpdf struct.
type RequestPdf struct {
	log *baselogger.BaseLogger
}

// Pdf signature struct.
type Signature struct {
	Name           string
	Reason         string
	Date           time.Time
	DateFormat     string
	Rect           []float64
	FontSize       float64
	MakeString     string
	OnPageNum      int
	PrivateKey     *rsa.PrivateKey
	Certificate    *x509.Certificate
	SignatureLines []SignatureLine
}

// Pdf signatureline struct.
type SignatureLine struct {
	Desc string
	Text string
}

// New request to pdf function.
func NewRequestPdf(log *baselogger.BaseLogger, unidocMeteredKey string) *RequestPdf {
	meteredState, _ := license.GetMeteredState()
	if !meteredState.OK {
		err := license.SetMeteredKey(unidocMeteredKey)
		if err != nil {
			panic(err)
		}
	}
	return &RequestPdf{
		log: log,
	}
}

// Generate pdf function.
func (r *RequestPdf) GenerateHTMLtoPDF(htmlString string, htmlPath string, pdfPath string, PageSizeA4 string, dpi uint) (bool, error) {
	//Write whole the body.
	err1 := ioutil.WriteFile(htmlPath, []byte(htmlString), 0644)
	if err1 != nil {
		r.log.Fatal(err1)
		return false, err1
	}
	r.log.Debug("After pdfoperations.GeneratePDF.WriteFile")
	file, err := os.Open(filepath.Clean(htmlPath))
	if file != nil {
		defer func() {
			_ = file.Close()
		}()
	}
	if err != nil {
		r.log.Fatal(err)
		return false, err
	}
	r.log.Debug("After pdfoperations.GeneratePDF.OpenFile")
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		os.Remove(htmlPath)
		r.log.Fatal(err)
		return false, err
	}
	pages := wkhtmltopdf.NewPageReader(file)
	pages.FooterCenter.Set("[page]")
	pages.FooterFontSize.Set(10)
	pdfg.AddPage(pages)

	pdfg.PageSize.Set(PageSizeA4)
	pdfg.MarginLeft.Set(uint(25))  // became 1 inch margin
	pdfg.MarginRight.Set(uint(25)) // became 1 inch margin

	pdfg.Dpi.Set(dpi)

	err = pdfg.Create()
	if err != nil {
		r.log.Fatal(err)
		return false, err
	}
	r.log.Debug("After pdfoperations.GeneratePDF.PDFCreate")
	err = pdfg.WriteFile(pdfPath)
	r.log.Debug("After pdfoperations.GeneratePDF.WritePDFFile")
	if err != nil {
		r.log.Fatal(err)
		return false, err
	}
	os.Remove(htmlPath)
	if err != nil {
		r.log.Infof("Error Removeing html file : %s , Error :%s", htmlPath, err)
		return true, nil
	}
	return true, nil
}

// Generate private key and certificate.
func (r *RequestPdf) GenerateKeyCertificate(name string, organizations []string) (*rsa.PrivateKey, *x509.Certificate, error) {
	//Generate private key.
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	var now = time.Now()
	//Initialize X509 certificate template.
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   name,
			Organization: organizations,
		},
		NotBefore:          now.Add(-time.Hour).UTC(),
		NotAfter:           now.Add(time.Hour * 24 * 365).UTC(),
		PublicKeyAlgorithm: x509.RSA,
		Issuer: pkix.Name{
			CommonName:   name,
			Organization: organizations,
		},
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	//Generate X509 certificate.
	certData, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv)
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(certData)
	if err != nil {
		return nil, nil, err
	}

	return priv, cert, nil
}

// Add Digital Signature to pdf file.
func (r *RequestPdf) AddDigitalSignature(pdfFilePath string, sig Signature, pageNum int, priv *rsa.PrivateKey, cert *x509.Certificate) error {
	//Create reader.
	file, err := os.Open(filepath.Clean(pdfFilePath))
	if err != nil {
		r.log.Fatal(err)
		return err
	}
	if file != nil {
		defer func() {
			_ = file.Close()
		}()
	}
	r.log.Debug("After pdfoperations.AddDigitalSignature.OpenFile")
	err = os.Remove(pdfFilePath)
	if err != nil {
		r.log.Fatal(err)
		return err
	}
	r.log.Debug("After pdfoperations.AddDigitalSignature.Remove")
	if file != nil {
		defer func() {
			_ = file.Close()
		}()
	}
	reader, err := model.NewPdfReader(file)
	if err != nil {
		r.log.Fatal(err)
		return err
	}
	r.log.Debug("After pdfoperations.AddDigitalSignature.NewPdfReader")
	//Create appender.
	appender, err := model.NewPdfAppender(reader)
	if err != nil {
		r.log.Fatal(err)
		return err
	}
	r.log.Debug("After pdfoperations.AddDigitalSignature.NewPdfAppender")
	//Create signature handler.
	handler, err := sighandler.NewAdobePKCS7Detached(priv, cert)
	if err != nil {
		r.log.Fatal(err)
		return err
	}
	r.log.Debug("After pdfoperations.AddDigitalSignature.NewAdobePKCS7Detached")
	//Create signature.
	signature := model.NewPdfSignature(handler)
	signature.SetName(sig.Name)
	signature.SetReason(sig.Reason)
	signature.SetDate(time.Now(), "")

	if err := signature.Initialize(); err != nil {
		r.log.Fatal(err)
		return err
	}
	r.log.Debug("After pdfoperations.AddDigitalSignature.NewPdfSignature")
	//Create signature field and appearance.
	opts := annotator.NewSignatureFieldOpts()
	opts.FontSize = sig.FontSize
	opts.Rect = sig.Rect
	opts.AutoSize = true
	var signatureLines []*annotator.SignatureLine
	for _, sLine := range sig.SignatureLines {
		signatureLines = append(signatureLines, annotator.NewSignatureLine(sLine.Desc, sLine.Text))
	}
	field, err := annotator.NewSignatureField(
		signature,
		signatureLines,
		opts,
	)
	field.T = core.MakeString(sig.MakeString)

	if err = appender.Sign(int(pageNum), field); err != nil {
		r.log.Fatal(err)
		return err
	}
	r.log.Debug("After pdfoperations.AddDigitalSignature.Sign")
	//Write output PDF file.
	err = appender.WriteToFile(pdfFilePath)
	if err != nil {
		r.log.Fatal(err)
		return err
	}
	r.log.Debug("After pdfoperations.AddDigitalSignature.WriteToFile")
	return nil
}

// Add Digital Signature to pdf file.
func (r *RequestPdf) AddMultiDigitalSignature(pdfFilePath string, signatureAry []Signature) error {
	//Create reader.
	file, err := os.Open(filepath.Clean(pdfFilePath))
	if err != nil {
		r.log.Fatal(err)
		return err
	}
	if file != nil {
		defer func() {
			_ = file.Close()
		}()
	}
	r.log.Debug("After pdfoperations.AddDigitalSignature.OpenFile")
	err = os.Remove(pdfFilePath)
	if err != nil {
		r.log.Fatal(err)
		return err
	}
	r.log.Debug("After pdfoperations.AddDigitalSignature.Remove")
	reader, err := model.NewPdfReader(file)
	if err != nil {
		r.log.Fatal(err)
		return err
	}
	r.log.Debug("After pdfoperations.AddDigitalSignature.NewPdfReader")
	//Create appender.
	appender, err := model.NewPdfAppender(reader)
	if err != nil {
		r.log.Fatal(err)
		return err
	}
	r.log.Debug("After pdfoperations.AddDigitalSignature.NewPdfAppender")

	for _, sig := range signatureAry {

		//Create signature handler.
		handler, err := sighandler.NewAdobePKCS7Detached(sig.PrivateKey, sig.Certificate)
		if err != nil {
			r.log.Fatal(err)
			return err
		}
		r.log.Debug("After pdfoperations.AddDigitalSignature.NewAdobePKCS7Detached")
		//Create signature.
		signature := model.NewPdfSignature(handler)
		signature.SetName(sig.Name)
		signature.SetReason(sig.Reason)
		signature.SetDate(time.Now(), "")

		if err := signature.Initialize(); err != nil {
			r.log.Fatal(err)
			return err
		}
		r.log.Debug("After pdfoperations.AddDigitalSignature.NewPdfSignature")
		//Create signature field and appearance.
		opts := annotator.NewSignatureFieldOpts()
		opts.FontSize = sig.FontSize
		opts.Rect = sig.Rect
		opts.AutoSize = true
		var signatureLines []*annotator.SignatureLine
		for _, sLine := range sig.SignatureLines {
			signatureLines = append(signatureLines, annotator.NewSignatureLine(sLine.Desc, sLine.Text))
		}
		field, err := annotator.NewSignatureField(
			signature,
			signatureLines,
			opts,
		)
		field.T = core.MakeString(sig.MakeString)

		if err = appender.Sign(int(sig.OnPageNum), field); err != nil {
			r.log.Fatal(err)
			return err
		}
		r.log.Debug("After pdfoperations.AddDigitalSignature.Sign")
	}
	//Write output PDF file.
	err = appender.WriteToFile(pdfFilePath)
	if err != nil {
		r.log.Fatal(err)
		return err
	}
	r.log.Debug("After pdfoperations.AddDigitalSignature.WriteToFile")
	return nil
}
