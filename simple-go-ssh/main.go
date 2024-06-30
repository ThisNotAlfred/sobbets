package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"net"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:5555")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		window := new(app.Window)
		err := run(window, &conn)
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()

	app.Main()

	conn.Close()
}

func send_command(conn *net.Conn, command string) string {
	fmt.Fprint(*conn, command)

	buffer := make([]byte, 4096)
	_, err := bufio.NewReader(*conn).Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	return string(buffer)
}

func run(window *app.Window, conn *net.Conn) error {
	theme := material.NewTheme()

	var ops op.Ops

	var editor widget.Editor
	editor.Submit = true

	var clickable widget.Clickable

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if clickable.Clicked(gtx) {
				result := send_command(conn, editor.Text())
				editor.SetText(result)
			}

			layout.Flex{
				Axis: layout.Horizontal, Alignment: layout.End, Spacing: layout.SpaceEnd,
			}.Layout(gtx, layout.Rigid(
				func(gtx layout.Context) layout.Dimensions {

					editor_style := material.Editor(theme, &editor, "editor")
					editor_style.Layout(gtx)
					editor_style.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}

					return editor_style.Layout(gtx)
				}))

			layout.Flex{
				Axis: layout.Vertical, Alignment: layout.End, Spacing: layout.SpaceStart,
			}.Layout(gtx, layout.Rigid(func(gtx layout.Context) layout.Dimensions {

				clickable_style := material.Button(theme, &clickable, "send_button")
				clickable_style.Color = color.NRGBA{R: 0, B: 0, G: 0, A: 255}
				clickable_style.Background = color.NRGBA{R: 150, G: 0, B: 100, A: 255}
				clickable_style.CornerRadius = 2

				return clickable_style.Layout(gtx)
			}))

			e.Frame(gtx.Ops)
		}
	}
}
