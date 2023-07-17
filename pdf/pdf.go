// 基于 github.com/chromedp/chromedp；需要确保本地有 chromium 内核可以使用
package pdf

import (
	"context"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// PrintToPDFWithHTML 根据提供的 url，生成 pdf
//
//	actions 在生成 PDF 前的等待动作，比如， chromedp.WaitVisible("#render_done")
func PrintToPDFWithURL(ctx context.Context, url string, actions ...chromedp.Action) (data []byte, err error) {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()
	if err := chromedp.Run(ctx, printToPDFWithURL(url, &data, actions...)); err != nil {
		return nil, err
	}
	return data, nil
}

func printToPDFWithURL(url string, res *[]byte, actions ...chromedp.Action) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("url"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// 如果有前置行为，先执行
			for _, action := range actions {
				action.Do(ctx)
			}
			header, footer := "", ""
			chromedp.InnerHTML("#page_header", &header, chromedp.ByQuery).Do(ctx)
			chromedp.InnerHTML("#page_footer", &footer, chromedp.ByQuery).Do(ctx)
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithDisplayHeaderFooter(true).
				WithHeaderTemplate(header).
				WithFooterTemplate(footer).
				WithMarginTop(0.8).
				// WithMarginRight(0.8).
				WithMarginBottom(0.8).
				// WithMarginLeft(0.8).
				Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}

// PrintToPDFWithHTML 根据提供的 html 内容，生成 pdf
//
//	actions 在生成 PDF 前的等待动作，比如， chromedp.WaitVisible("#render_done")
func PrintToPDFWithHTML(ctx context.Context, html string, actions ...chromedp.Action) (data []byte, err error) {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()
	if err := chromedp.Run(ctx, printToPDFWithHTML(html, &data, actions...)); err != nil {
		return nil, err
	}
	return data, nil
}

func printToPDFWithHTML(html string, res *[]byte, actions ...chromedp.Action) chromedp.Tasks {
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
			// 如果有前置行为，先执行
			for _, action := range actions {
				action.Do(ctx)
			}
			header, footer := "", ""
			chromedp.InnerHTML("#page_header", &header, chromedp.ByQuery).Do(ctx)
			chromedp.InnerHTML("#page_footer", &footer, chromedp.ByQuery).Do(ctx)
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithDisplayHeaderFooter(true).
				WithHeaderTemplate(header).
				WithFooterTemplate(footer).
				WithMarginTop(0.8).
				// WithMarginRight(0.8).
				WithMarginBottom(0.8).
				// WithMarginLeft(0.8).
				Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
