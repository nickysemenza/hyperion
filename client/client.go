package client

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	pb "github.com/nickysemenza/hyperion/api/proto"
	"github.com/nickysemenza/hyperion/core/config"
	"google.golang.org/grpc"

	ui "github.com/gizak/termui"
)

//Run runs the client
func Run(ctx context.Context) {

	go data(ctx)
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	p := ui.NewPar(":PRESS q TO QUIT DEMO")
	p.Height = 3
	p.Width = 50
	p.TextFgColor = ui.ColorWhite
	p.BorderLabel = "Text Box"
	p.BorderFg = ui.ColorCyan

	g := ui.NewGauge()
	g.Percent = 50
	g.Width = 50
	g.Height = 3
	g.Y = 11
	g.BorderLabel = "Gauge"
	g.BarColor = ui.ColorRed
	g.BorderFg = ui.ColorWhite
	g.BorderLabelFg = ui.ColorCyan

	ll := ui.NewPar("light list")
	ll.Height = 10
	ll.Width = 50
	ll.TextFgColor = ui.ColorWhite
	ll.BorderLabel = "light list"
	ll.BorderFg = ui.ColorRGB(200, 100, 23)

	rows2 := [][]string{
		[]string{"header1", "header2", "header3"},
		[]string{"Foundations", "Go-lang is so cool", "Im working on Ruby"},
		[]string{"2016", "11", "11"},
	}

	table2 := ui.NewTable()
	table2.Rows = rows2
	table2.FgColor = ui.ColorWhite
	table2.BgColor = ui.ColorDefault
	table2.TextAlign = ui.AlignCenter
	table2.Separator = false
	table2.Analysis()
	table2.SetSize()
	table2.BgColors[2] = ui.ColorRGB(200, 100, 23)
	table2.Y = 10
	table2.X = 0
	table2.Border = true

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, p),
			ui.NewCol(6, 0, g)),
		ui.NewRow(
			ui.NewCol(6, 0, ll),
			ui.NewCol(6, 0, table2),
			// 	ui.NewCol(3, 0, widget30, widget31, widget32),
			// 	ui.NewCol(6, 0, widget4)
		),
	)

	// calculate layout
	ui.Body.Align()

	ui.SetOutputMode(ui.Output256)

	ui.AddColorMap("aa", ui.ColorRGB(255, 90, 123))

	draw := func(data string) {
		ll.Text = data
		// g.Percent = count % 101
		// l.Items = listData[count%9:]
		// sls.Lines[0].Data = sparklineData[:30+count%50]
		// sls.Lines[1].Data = sparklineData[:35+count%50]
		// lc.Data["default"] = sinData[count/2%220:]
		// lc2.Data["default"] = sinData[2*count%220:]
		// bc.Data = barchartData[count/2%10:]

		ui.Render(ui.Body)
	}

	ui.Handle("q", func(ui.Event) {
		// press q to quit
		ui.StopLoop()
	})

	drawTicker := time.NewTicker(time.Second)
	drawTickerCount := 1
	go func() {
		for {
			draw(lightListPretty)

			drawTickerCount++
			<-drawTicker.C
		}
	}()

	ui.Loop()
}

var lightListPretty string

func data(ctx context.Context) {

	config := config.GetClientConfig(ctx)
	conn, cerr := grpc.Dial(config.ServerAddress, grpc.WithInsecure())
	if cerr != nil {
		log.Println(cerr)
	}
	defer conn.Close()

	client := pb.NewAPIClient(conn)

	lights := make(map[string]*pb.Light)
	stream, err := client.StreamGetLights(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatal(client, err)
	}

	for {
		received, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(client, err)
		}

		// spew.Dump(received)
		for _, l := range received.Lights {
			lights[l.Name] = l
		}

		var keys []string
		for k := range lights {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var b strings.Builder
		for _, k := range keys {
			rgb := lights[k].CurrentColor
			colorBlock := fmt.Sprintf("[███](bg-aa) %v", rgb.R) // rgbterm.Bytes([]byte("███"), uint8(rgb.R), uint8(rgb.G), uint8(rgb.B), 0, 0, 0)
			fmt.Fprintf(&b, "%s, %v\n", k, lights[k])
			fmt.Fprintf(&b, "%s %s\n", colorBlock, k)

		}
		lightListPretty = b.String()
	}
}
