package pdfgenerator

import (
	"context"
	"nearbyassist/internal/config"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func GeneratePDF(url string) ([]byte, error) {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("/usr/bin/google-chrome"),
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("ignore-certificate-errors", true),
	)

	alloc, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(alloc)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(config.Instance.PDF_GENERATION_TIMEOUT))
	defer cancel()

	var pdfbuf []byte
	err := chromedp.Run(
		ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Evaluate(`Promise.all(Array.from(document.images).map(img => {
            if (img.complete) return Promise.resolve();
            return new Promise(resolve => img.onload = img.onerror = resolve);
        }))`, nil),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfbuf, _, err = page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, err
	}

	return pdfbuf, nil
}
