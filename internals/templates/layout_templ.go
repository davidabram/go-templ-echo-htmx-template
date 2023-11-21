// Code generated by templ@v0.2.364 DO NOT EDIT.

package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

type Page struct {
	Title   string
	Boosted bool
}

func ClerkScript() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_1 := templ.GetChildren(ctx)
		if var_1 == nil {
			var_1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<script>")
		if err != nil {
			return err
		}
		var_2 := `
		const clerkPublishableKey = 'pk_test_b25lLXJheS04LmNsZXJrLmFjY291bnRzLmRldiQ';
		const frontendApi = 'one-ray-8.clerk.accounts.dev';
		const version = '@latest';

		const script = document.createElement('script');
		script.setAttribute('data-clerk-frontend-api', frontendApi);
		script.setAttribute('data-clerk-publishable-key', clerkPublishableKey);
		script.async = true;
		script.src = ` + "`" + `https://${frontendApi}/npm/@clerk/clerk-js${version}/dist/clerk.browser.js` + "`" + `;

		script.addEventListener('load', async function () {
			await window.Clerk.load({

			});
		});
		document.body.appendChild(script);
	`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func Layout(page *Page) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_3 := templ.GetChildren(ctx)
		if var_3 == nil {
			var_3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<html>")
		if err != nil {
			return err
		}
		err = Head(page.Title).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		if !page.Boosted {
			err = Navigation(false).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("<body>")
		if err != nil {
			return err
		}
		var_4 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			templBuffer, templIsBuffer := w.(*bytes.Buffer)
			if !templIsBuffer {
				templBuffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templBuffer)
			}
			err = var_3.Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			if !templIsBuffer {
				_, err = io.Copy(w, templBuffer)
			}
			return err
		})
		err = Content().Render(templ.WithChildren(ctx, var_4), templBuffer)
		if err != nil {
			return err
		}
		if !page.Boosted {
			err = ClerkScript().Render(ctx, templBuffer)
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</body></html>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func Head(title string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_5 := templ.GetChildren(ctx)
		if var_5 == nil {
			var_5 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<head><meta charset=\"UTF-8\"><title>")
		if err != nil {
			return err
		}
		var var_6 string = title
		_, err = templBuffer.WriteString(templ.EscapeString(var_6))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</title><script src=\"https://unpkg.com/htmx.org@1.9.6/dist/htmx.min.js\">")
		if err != nil {
			return err
		}
		var_7 := ``
		_, err = templBuffer.WriteString(var_7)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script><script src=\"https://unpkg.com/htmx.org@1.9.6/dist/ext/head-support.js\">")
		if err != nil {
			return err
		}
		var_8 := ``
		_, err = templBuffer.WriteString(var_8)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script><link href=\"/style.css\" rel=\"stylesheet\"><script>")
		if err != nil {
			return err
		}
		var_9 := `
		`
		_, err = templBuffer.WriteString(var_9)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script></head>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func Content() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_10 := templ.GetChildren(ctx)
		if var_10 == nil {
			var_10 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<main class=\"bg-red-500\">")
		if err != nil {
			return err
		}
		err = var_10.Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</main>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func Navigation(isSignedIn bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_11 := templ.GetChildren(ctx)
		if var_11 == nil {
			var_11 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<nav hx-boost=\"true\" hx-target=\"main\" hx-swap=\"outerHTML show:unset\" class=\"bg-blue-500 p-4\"><a href=\"/\" class=\"text-white hover:text-blue-200\">")
		if err != nil {
			return err
		}
		var_12 := `Home`
		_, err = templBuffer.WriteString(var_12)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a><a href=\"/about\" class=\"text-white hover:text-blue-200\">")
		if err != nil {
			return err
		}
		var_13 := `About`
		_, err = templBuffer.WriteString(var_13)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a><a href=\"/books\" class=\"text-white hover:text-blue-200\">")
		if err != nil {
			return err
		}
		var_14 := `Books`
		_, err = templBuffer.WriteString(var_14)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a><a href=\"/charts\" class=\"text-white hover:text-blue-200\">")
		if err != nil {
			return err
		}
		var_15 := `Charts`
		_, err = templBuffer.WriteString(var_15)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a><a href=\"/contact\" class=\"text-white hover:text-blue-200\">")
		if err != nil {
			return err
		}
		var_16 := `Contact`
		_, err = templBuffer.WriteString(var_16)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a><div id=\"user-button\"></div></nav><script type=\"module\">")
		if err != nil {
			return err
		}
		var_17 := `
      const script = document.createElement("script");
      script.async = true;

	  const userButton = document.getElementById('user-button');

      window.addEventListener("load", async function () {
        await window.Clerk.load();

		document.addEventListener("htmx:removingHeadElement", function (event) {
			if(event.detail.headElement.getAttribute('data-emotion') === 'cl-internal' && event.detail.headElement.hasAttributes('data-s')) {
				event.preventDefault();
			}
		});

        if (Clerk.user) {
			window.Clerk.mountUserButton(userButton);
        } else {
			document.location.href="/signin";
		}

        Clerk.addListener(async ({ user }) => {
          if (user) {
			window.Clerk.mountUserButton(userButton);
          } else {
			document.location.href="/signin";
		}
        });
      });
      document.body.appendChild(script);
    `
		_, err = templBuffer.WriteString(var_17)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script><hr class=\"my-4\">")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
