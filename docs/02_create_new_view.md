# Create a new view

Assuming that you have already [setup the controller](./01_create_new_controller.md), we know that the goliveview controller
exposes a [Handler](https://pkg.go.dev/github.com/goliveview/controller#Controller) api which accepts a type which satisfies 
the [View](https://pkg.go.dev/github.com/goliveview/controller#View) interface.

```go
glvc := glv.Websocket("goliveview-starter", glv.DevelopmentMode(mode))
r := chi.NewRouter()
...
r.NotFound(glvc.Handler(&views.NotfoundView{}))
```

The `View` interface:

```go
type View interface {
	Content() string
	Layout() string
	OnMount(w http.ResponseWriter, r *http.Request) (Status, M)
	OnEvent(ctx Context) error
	LayoutContentName() string
	Partials() []string
	Extensions() []string
	FuncMap() template.FuncMap
}
```

To keep the boilerplate to the minimum, the `controller` package exposes a [DefaultView](https://pkg.go.dev/github.com/goliveview/controller#DefaultView)
The `DefaultView` implements the `View` interface using sane defaults. A new view can satisfy the `View` interface by
simply embedding the `DefaultView`.

```go
package views

import (
	glv "github.com/goliveview/controller"
)

type NotfoundView struct {
	glv.DefaultView
}
```

When the above view is rendered by `r.NotFound(glvc.Handler(&views.NotfoundView{}))`, the default layout and content
are used.

```go
func (d DefaultView) Content() string {
	return "./templates/index.html"
}

func (d DefaultView) Layout() string {
	return "./templates/layouts/index.html"
}
```

Here we want to show a custom 404 page, so we should override the `Content` and `Layout` methods.

```go
package views

import (
	glv "github.com/goliveview/controller"
)

type NotfoundView struct {
	glv.DefaultView
}

func (n *NotfoundView) Content() string {
	return "./templates/404.html"
}

func (n *NotfoundView) Layout() string {
	return "./templates/layouts/error.html"
}
```

Now when the view is rendered by `r.NotFound(glvc.Handler(&views.NotfoundView{}))`, it displays `./templates/404.html`
within the layout`./templates/layouts/error.html`.