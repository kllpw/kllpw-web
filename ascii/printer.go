package ascii

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/mbndr/figlet4go"
	"google.golang.org/api/option"
)

var ascii, options = initASCII()

func initASCII() (*figlet4go.AsciiRender, *figlet4go.RenderOptions) {
	ascii := figlet4go.NewAsciiRender()
	// Adding the colours to RenderOptions
	options := figlet4go.NewRenderOptions()
	options.FontName = "bloody"
	ascii.LoadFont("./ascii/fonts/")
	return ascii, options
}

// RenderString returns the string in the format defined in initASCII
func RenderString(text string) string {
	renderStr, _ := ascii.RenderOpts(text, options)
	return renderStr
}
