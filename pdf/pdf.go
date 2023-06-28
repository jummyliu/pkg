// 基于 github.com/chromedp/chromedp；需要确保本地有 chromium 内核可以使用
package pdf

import (
	"context"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// PrintToPDFWithHTML 根据提供的 url，生成 pdf
func PrintToPDFWithURL(ctx context.Context, url string) (data []byte, err error) {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()
	if err := chromedp.Run(ctx, printToPDFWithURL(url, &data)); err != nil {
		return nil, err
	}
	return data, nil
}

func printToPDFWithURL(url string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("url"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}

// PrintToPDFWithHTML 根据提供的 html 内容，生成 pdf
func PrintToPDFWithHTML(ctx context.Context, html string) (data []byte, err error) {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()
	if err := chromedp.Run(ctx, printToPDFWithHTML(html, &data)); err != nil {
		return nil, err
	}
	return data, nil
}

func printToPDFWithHTML(html string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		// 空白页
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return nil
			}
			// 渲染 html 内容
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
