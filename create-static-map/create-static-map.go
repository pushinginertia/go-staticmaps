// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/flopp/go-coordsparser"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/jessevdk/go-flags"
	sm "github.com/pushinginertia/go-staticmaps"
	"go.uber.org/zap"
)

func handleTypeOption(ctx *sm.Context, parameter string) {
	tileProviders := sm.GetTileProviders()
	tp := tileProviders[parameter]
	if tp != nil {
		ctx.SetTileProvider(tp)
		return
	}

	if parameter != "list" {
		fmt.Println("Bad map type:", parameter)
	}
	fmt.Println("Possible map types (to be used with --type/-t):")
	// print sorted keys
	keys := make([]string, 0, len(tileProviders))
	for k := range tileProviders {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println(k)
	}
	os.Exit(0)
}

func handleCenterOption(ctx *sm.Context, parameter string) error {
	lat, lng, err := coordsparser.Parse(parameter)
	if err != nil {
		return err
	}
	ctx.SetCenter(s2.LatLngFromDegrees(lat, lng))
	return nil
}

func handleBboxOption(ctx *sm.Context, parameter string) error {
	pair := strings.Split(parameter, "|")
	if len(pair) != 2 {
		return fmt.Errorf("bad NW|SE coordinates pair: %s", parameter)
	}

	var err error
	var nwlat float64
	var nwlng float64
	var selat float64
	var selng float64
	nwlat, nwlng, err = coordsparser.Parse(pair[0])
	if err != nil {
		return err
	}
	selat, selng, err = coordsparser.Parse(pair[1])
	if err != nil {
		return err
	}

	var bbox *s2.Rect
	bbox, err = sm.CreateBBox(nwlat, nwlng, selat, selng)
	if err != nil {
		return err
	}

	ctx.SetBoundingBox(*bbox)
	return nil
}

func handleBackgroundOption(ctx *sm.Context, parameter string) error {
	color, err := sm.ParseColorString(parameter)
	if err != nil {
		return err
	}

	ctx.SetBackground(color)
	return nil
}

func handleMarkersOption(ctx *sm.Context, parameters []string) error {
	for _, s := range parameters {
		markers, err := sm.ParseMarkerString(s)
		if err != nil {
			return err
		}
		for _, marker := range markers {
			ctx.AddObject(marker)
		}
	}
	return nil
}

func handleImageMarkersOption(ctx *sm.Context, parameters []string) error {
	for _, s := range parameters {
		markers, err := sm.ParseImageMarkerString(s)
		if err != nil {
			return err
		}
		for _, marker := range markers {
			ctx.AddObject(marker)
		}
	}
	return nil
}

func handlePathsOption(ctx *sm.Context, parameters []string) error {
	for _, s := range parameters {
		paths, err := sm.ParsePathString(s)
		if err != nil {
			return err
		}
		for _, path := range paths {
			ctx.AddObject(path)
		}
	}
	return nil
}

func handleAreasOption(ctx *sm.Context, parameters []string) error {
	for _, s := range parameters {
		area, err := sm.ParseAreaString(s)
		if err != nil {
			return err
		}
		ctx.AddObject(area)
	}
	return nil
}

func handleCirclesOption(ctx *sm.Context, parameters []string) error {
	for _, s := range parameters {
		circles, err := sm.ParseCircleString(s)
		if err != nil {
			return err
		}
		for _, circle := range circles {
			ctx.AddObject(circle)
		}
	}
	return nil
}

type Opts struct {
	//		ClearCache bool     `long:"clear-cache" description:"Clears the tile cache"`
	Width        int      `long:"width" description:"Width of the generated static map image" value-name:"PIXELS" default:"512"`
	Height       int      `long:"height" description:"Height of the generated static map image" value-name:"PIXELS" default:"512"`
	Output       string   `short:"o" long:"output" description:"Output file name" value-name:"FILENAME" default:"map.png"`
	Type         string   `short:"t" long:"type" description:"Select the map type; list possible map types with '--type list'" value-name:"MAPTYPE"`
	Center       string   `short:"c" long:"center" description:"Center coordinates (lat,lng) of the static map" value-name:"LATLNG"`
	Zoom         int      `short:"z" long:"zoom" description:"Zoom factor" value-name:"ZOOMLEVEL"`
	BBox         string   `short:"b" long:"bbox" description:"Bounding box of the static map" value-name:"nwLATLNG|seLATLNG"`
	Background   string   `long:"background" description:"Background color" value-name:"COLOR" default:"transparent"`
	UserAgent    string   `short:"u" long:"useragent" description:"Overwrite the default HTTP user agent string" value-name:"USERAGENT"`
	Markers      []string `short:"m" long:"marker" description:"Add a marker to the static map" value-name:"MARKER"`
	ImageMarkers []string `short:"i" long:"imagemarker" description:"Add an image marker to the static map" value-name:"MARKER"`
	Paths        []string `short:"p" long:"path" description:"Add a path to the static map" value-name:"PATH"`
	Areas        []string `short:"a" long:"area" description:"Add an area to the static map" value-name:"AREA"`
	Circles      []string `short:"C" long:"circle" description:"Add a circle to the static map" value-name:"CIRCLE"`
}

func main() {

	// Fetch the default production configuration and redirect the output destination
	// from stderr to stdout, then build the logger instance from the config.
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	log, err := config.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize zap logger: %v\n", err)
		os.Exit(1)
	}

	var opts Opts
	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	parser.LongDescription = `Creates a static map`
	if _, err := parser.Parse(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse flags: %v\n", err)
		os.Exit(1)
	}

	if parser.FindOptionByLongName("help").IsSet() {
		parser.WriteHelp(os.Stdout)
		os.Exit(0)
	}

	ctx := sm.NewContext(log)

	if err := run(ctx, parser, opts); err != nil {
		fmt.Fprintf(os.Stderr, "Run failure: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx *sm.Context, parser *flags.Parser, opts Opts) error {
	if parser.FindOptionByLongName("type").IsSet() {
		handleTypeOption(ctx, opts.Type)
	}

	ctx.SetSize(opts.Width, opts.Height)

	if parser.FindOptionByLongName("zoom").IsSet() {
		ctx.SetZoom(opts.Zoom)
	}

	if parser.FindOptionByLongName("center").IsSet() {
		if err := handleCenterOption(ctx, opts.Center); err != nil {
			return err
		}
	}

	if parser.FindOptionByLongName("bbox").IsSet() {
		if err := handleBboxOption(ctx, opts.BBox); err != nil {
			return err
		}
	}

	if parser.FindOptionByLongName("background").IsSet() {
		if err := handleBackgroundOption(ctx, opts.Background); err != nil {
			return err
		}
	}

	if parser.FindOptionByLongName("useragent").IsSet() {
		ctx.SetUserAgent(opts.UserAgent)
	}

	if err := handleAreasOption(ctx, opts.Areas); err != nil {
		return err
	}
	if err := handleMarkersOption(ctx, opts.Markers); err != nil {
		return err
	}
	if err := handleImageMarkersOption(ctx, opts.ImageMarkers); err != nil {
		return err
	}
	if err := handleCirclesOption(ctx, opts.Circles); err != nil {
		return err
	}
	if err := handlePathsOption(ctx, opts.Paths); err != nil {
		return err
	}

	img, _, err := ctx.Render()
	if err != nil {
		return fmt.Errorf("failed to render image: %w", err)
	}

	if err = gg.SavePNG(opts.Output, img); err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	return nil
}
