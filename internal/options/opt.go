package options

// Opts represent struct for all program options
type Opts struct {
	SavePath  string   `short:"s" long:"save-path" env:"SAVE_PATH" default:"./result" description:"where app will download tiles"`
	URL       string   `short:"p" long:"provider-url" env:"PROVIDER_URL" required:"true" description:"url where map-tile-provider serving"`
	Zoom      int      `short:"z" long:"zoom" env:"ZOOM" default:"20" description:"zoom for downloading tiles"`
	Side      int      `long:"side" env:"SIDE" default:"3" description:"square side for each tile"`
	SetMax    bool     `long:"set-max" env:"SET_MAX" description:"if provider max zoom is min when specified in zoom app will dowload tile in max zoom for provider"`
	Providers []string `long:"providers" env:"PROVIDERS" description:"providers to download if empty - download all"`
	Points    string   `long:"points" env:"POINTS" required:"true" description:"what points download"`
}
