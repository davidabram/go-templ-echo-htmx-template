package custom_chart

// RendererProvider is a function that returns a renderer.
type CustomRendererProvider func(int, int) (CustomRenderer, error)
