package ascii

import "github.com/mbndr/figlet4go"

var ascii, options = initASCII()

func initASCII() (*figlet4go.AsciiRender, *figlet4go.RenderOptions) {
	ascii := figlet4go.NewAsciiRender()
	// Adding the colors to RenderOptions
	options := figlet4go.NewRenderOptions()
	options.FontName = "bloody"
	options.FontColor = []figlet4go.Color{
		// Colors can be given by default ansi color codes...
		figlet4go.ColorCyan,
	}
	ascii.LoadFont("./static/fonts/")
	return ascii, options
}

// RenderString returns the string in the format defined in initASCII
func RenderString(text string) string {
	renderStr, _ := ascii.RenderOpts(text, options)
	return renderStr
}

// RenderStringHTML returns the string in the format defined in initASCII
func RenderStringHTML(text string) string {
	options.FontColor = []figlet4go.Color{}
	renderStr, _ := ascii.RenderOpts(text, options)
	return "<pre>" + renderStr + "</pre>"
}
