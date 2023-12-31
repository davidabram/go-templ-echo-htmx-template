// Code generated by templ@v0.2.364 DO NOT EDIT.

package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func SignIn(page *Page, name string, htmxHeaders string) templ.Component {
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
		_, err = templBuffer.WriteString("<html>")
		if err != nil {
			return err
		}
		err = Head(page.Title).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<body class=\"flex items-center justify-center\"><div id=\"sign-in\"></div>")
		if err != nil {
			return err
		}
		if !page.Boosted {
			err = ClerkScript().Render(ctx, templBuffer)
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("<script type=\"module\">")
		if err != nil {
			return err
		}
		var_2 := `
      const script = document.createElement("script");
      script.async = true;

      const signInDiv = document.getElementById('sign-in');

      window.addEventListener("load", async function () {
        await window.Clerk.load();

        if (Clerk.user) {
          document.location.href="/";
        } else {
          window.Clerk.mountSignIn(signInDiv);
        }

        Clerk.addListener(({ user }) => {
          if (user) {
            document.location.href="/";
          } else {console.log('render')
            window.Clerk.mountSignIn(signInDiv);
          }
        });
      });
      document.body.appendChild(script);
    `
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script></body></html>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
