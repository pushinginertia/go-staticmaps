// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sm

import "fmt"

// TileProvider encapsulates all infos about a map tile provider service (name, url scheme, attribution, etc.)
type TileProvider struct {
	Name           string // writes to {cachedir}/{Name}
	Attribution    string
	IgnoreNotFound bool
	TileSize       int
	URLPattern     string // "%[1]s" => shard, "%[2]d" => zoom, "%[3]d" => x, "%[4]d" => y
	Shards         []string
}

func (t *TileProvider) getURL(shard string, zoom, x, y int) string {
	return fmt.Sprintf(t.URLPattern, shard, zoom, x, y)
}

// NewTileProviderOpenStreetMaps creates a TileProvider struct for OSM's tile service
func NewTileProviderOpenStreetMaps() *TileProvider {
	t := new(TileProvider)
	t.Name = "osm"
	t.Attribution = "Maps and Data (c) openstreetmap.org and contributors, ODbL"
	t.TileSize = 256
	t.URLPattern = "http://%[1]s.tile.openstreetmap.org/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"a", "b", "c"}
	return t
}

// Use ppi=72 for normal size or ppi=320 for retina output at 2x.
//
// Use mapId=newest to get the latest map tiles, or see this link for how to
// check for the current hash value:
// https://developer.here.com/documentation/map-tile/topics/map-versions.html
//
// If using a cache and mapId=newest, it should be cleared periodically or
// generated maps may not be correct when matching older cached tiles with
// newer fetched tiles.
//
// See https://developer.here.com/documentation/map-tile/topics/resource-base-maptile.html
// for valid values for map ID and scheme.
func NewTileProviderHere(mapId, scheme string, ppi int, appId, appCode string) *TileProvider {
	// URL Structure:
	// https://{1-4}.base.maps.api.here.com
	// /maptile/2.1/maptile/{map id}/{scheme}/{zoom}/{column}/{row}/{size}/{format}
	// ?app_id={YOUR_APP_ID}
	// &app_code={YOUR_APP_CODE}
	// &{param}={value}
	return &TileProvider{
		Name:        "here",
		Attribution: "here.com",
		TileSize:    512,
		URLPattern: fmt.Sprintf(
			"https://%%[1]s.base.maps.cit.api.here.com/maptile/2.1/maptile/%s/%s/%%[2]d/%%[3]d/%%[4]d/512/png?ppi=%v&app_id=%s&app_code=%s",
			mapId,
			scheme,
			ppi,
			appId,
			appCode,
		),
		Shards: []string{"1", "2", "3", "4"},
	}
}

func newTileProviderThunderforest(name string) *TileProvider {
	t := new(TileProvider)
	t.Name = fmt.Sprintf("thunderforest-%s", name)
	t.Attribution = "Maps (c) Thundeforest; Data (c) OSM and contributors, ODbL"
	t.TileSize = 256
	t.URLPattern = "https://%[1]s.tile.thunderforest.com/" + name + "/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"a", "b", "c"}
	return t
}

// NewTileProviderThunderforestLandscape creates a TileProvider struct for thundeforests's 'landscape' tile service
func NewTileProviderThunderforestLandscape() *TileProvider {
	return newTileProviderThunderforest("landscape")
}

// NewTileProviderThunderforestOutdoors creates a TileProvider struct for thundeforests's 'outdoors' tile service
func NewTileProviderThunderforestOutdoors() *TileProvider {
	return newTileProviderThunderforest("outdoors")
}

// NewTileProviderThunderforestTransport creates a TileProvider struct for thundeforests's 'transport' tile service
func NewTileProviderThunderforestTransport() *TileProvider {
	return newTileProviderThunderforest("transport")
}

// NewTileProviderStamenToner creates a TileProvider struct for stamens' 'toner' tile service
func NewTileProviderStamenToner() *TileProvider {
	t := new(TileProvider)
	t.Name = "stamen-toner"
	t.Attribution = "Maps (c) Stamen; Data (c) OSM and contributors, ODbL"
	t.TileSize = 256
	t.URLPattern = "http://%[1]s.tile.stamen.com/toner/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"a", "b", "c", "d"}
	return t
}

// NewTileProviderStamenTerrain creates a TileProvider struct for stamens' 'terrain' tile service
func NewTileProviderStamenTerrain() *TileProvider {
	t := new(TileProvider)
	t.Name = "stamen-terrain"
	t.Attribution = "Maps (c) Stamen; Data (c) OSM and contributors, ODbL"
	t.TileSize = 256
	t.URLPattern = "http://%[1]s.tile.stamen.com/terrain/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"a", "b", "c", "d"}
	return t
}

// NewTileProviderOpenTopoMap creates a TileProvider struct for opentopomap's tile service
func NewTileProviderOpenTopoMap() *TileProvider {
	t := new(TileProvider)
	t.Name = "opentopomap"
	t.Attribution = "Maps (c) OpenTopoMap [CC-BY-SA]; Data (c) OSM and contributors [ODbL]; Data (c) SRTM"
	t.TileSize = 256
	t.URLPattern = "http://%[1]s.tile.opentopomap.org/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"a", "b", "c"}
	return t
}

// NewTileProviderWikimedia creates a TileProvider struct for Wikimedia's tile service
func NewTileProviderWikimedia() *TileProvider {
	t := new(TileProvider)
	t.Name = "wikimedia"
	t.Attribution = "Map (c) Wikimedia; Data (c) OSM and contributors, ODbL."
	t.TileSize = 256
	t.URLPattern = "https://maps.wikimedia.org/osm-intl/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{}
	return t
}

// NewTileProviderOpenCycleMap creates a TileProvider struct for OpenCycleMap's tile service
func NewTileProviderOpenCycleMap() *TileProvider {
	t := new(TileProvider)
	t.Name = "cycle"
	t.Attribution = "Maps and Data (c) openstreetmaps.org and contributors, ODbL"
	t.TileSize = 256
	t.URLPattern = "http://%[1]s.tile.opencyclemap.org/cycle/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"a", "b"}
	return t
}

func newTileProviderCarto(name string) *TileProvider {
	t := new(TileProvider)
	t.Name = fmt.Sprintf("carto-%s", name)
	t.Attribution = "Map (c) Carto [CC BY 3.0] Data (c) OSM and contributors, ODbL."
	t.TileSize = 256
	t.URLPattern = "https://cartodb-basemaps-%[1]s.global.ssl.fastly.net/" + name + "_all/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"a", "b", "c", "d"}
	return t
}

// NewTileProviderCartoLight creates a TileProvider struct for Carto's tile service (light variant)
func NewTileProviderCartoLight() *TileProvider {
	return newTileProviderCarto("light")
}

// NewTileProviderCartoDark creates a TileProvider struct for Carto's tile service (dark variant)
func NewTileProviderCartoDark() *TileProvider {
	return newTileProviderCarto("dark")
}

// NewTileProviderArcgisWorldImagery creates a TileProvider struct for Arcgis' WorldImagery tiles
func NewTileProviderArcgisWorldImagery() *TileProvider {
	t := new(TileProvider)
	t.Name = "arcgis-worldimagery"
	t.Attribution = "Source: Esri, Maxar, GeoEye, Earthstar Geographics, CNES/Airbus DS, USDA, USGS, AeroGRID, IGN, and the GIS User Community"
	t.TileSize = 256
	t.URLPattern = "https://server.arcgisonline.com/arcgis/rest/services/World_Imagery/MapServer/tile/%[2]d/%[4]d/%[3]d"
	t.Shards = []string{}
	return t
}

// GetTileProviders returns a map of all available TileProviders
func GetTileProviders() map[string]*TileProvider {
	m := make(map[string]*TileProvider)

	list := []*TileProvider{
		NewTileProviderOpenStreetMaps(),
		NewTileProviderOpenCycleMap(),
		NewTileProviderThunderforestLandscape(),
		NewTileProviderThunderforestOutdoors(),
		NewTileProviderThunderforestTransport(),
		NewTileProviderStamenToner(),
		NewTileProviderStamenTerrain(),
		NewTileProviderOpenTopoMap(),
		NewTileProviderOpenStreetMaps(),
		NewTileProviderOpenCycleMap(),
		NewTileProviderCartoLight(),
		NewTileProviderCartoDark(),
		NewTileProviderArcgisWorldImagery(),
	}

	for _, tp := range list {
		m[tp.Name] = tp
	}

	return m
}
