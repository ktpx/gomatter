package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Struct for JSON Payload
type Payload struct {
	Text    string `json:"text"`
	Channel string `json:"channel,omitempty"`
	Url     string `json:"icon_url,omitempty"`
	Emoji   string `json:"icon_emoji,omitempty"`
	User    string `json:"username,omitempty"`
	Attach  string `json:"attachment,omitempty"`
}

// Globals, change or control via env or parameters
var url = ""
var default_channel = ""
var ver = "v1.1"

// Map for known icon urls for easy reference (add your own)
var icon_urls = map[string]string{
	"mattermost": "https://mattermost.com/wp-content/uploads/2022/02/icon.png",
}

// Read from stdin
func readStdin() string {
	var s string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s = fmt.Sprintf("%s%s", s, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	return s
}
func readFromFile(f, t *string) error {

	r, err := os.ReadFile(*f)
	if err != nil {
		return err
	}
	*t = string(r)
	return nil
}

func main() {
	var stdin, verbose bool

	p := Payload{}

	flag.Usage = func() {
		fmt.Printf("Gomatter %s (https://github.com/ktpx/gomatter)\n", ver)
		fmt.Println("Usage: gomatter -c channel -m message [-n username] [-r] .. ")
		flag.PrintDefaults() // prints default usage
	}
	flag.StringVar(&p.Channel, "c", default_channel, "Specify a channel.")
	flag.StringVar(&p.Text, "m", "", "Specify text message.")
	flag.StringVar(&p.Emoji, "e", "", "Specify an icon_emoji.")
	flag.StringVar(&p.Url, "i", "", "Specify an icon url.")
	flag.StringVar(&p.User, "u", "", "Specify a Username.")
	flag.StringVar(&p.Attach, "a", "", "Specify attachments.")
	flag.StringVar(&url, "w", "", "Specify webhook url.")
	var fname = flag.String("f", "", "Read text from a file.")
	flag.BoolVar(&stdin, "r", false, "Read from stdin (presedence over -m)")
	flag.BoolVar(&verbose, "v", false, "Be move verbose.")
	appicon := flag.String("k", "", "Specify predefined app icon (if defined)")
	flag.Parse()

	if stdin {
		p.Text = readStdin()
	}
	if len(*fname) > 0 {
		err := readFromFile(fname, &p.Text)
		if err != nil {
			log.Fatalf("Error reading from file '%s'", *fname)
		}
	}
	if len(p.Text) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	if *appicon != "" {
		p.Url, _ = icon_urls[*appicon]
	}
	if len(p.Channel) == 0 {
		if v, exists := os.LookupEnv("MM_DEFAULT_CHANNEL"); exists {
			p.Channel = v
		} else {
			log.Fatal("No channel has been specified or set.")
		}
	}
	if len(url) == 0 {
		if v, exists := os.LookupEnv("MM_WEBHOOKURL"); exists {
			url = v
		} else {
			log.Fatal("No URL has been set.")
		}
	}

	json, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.Post(url, "application/json", bytes.NewBuffer(json))
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if verbose {
		fmt.Println("params: %v\n", p)
		fmt.Println("response Status:", res.Status)
		fmt.Println("response Headers:", res.Header)
	}
	body, _ := ioutil.ReadAll(res.Body)
	if verbose {
		fmt.Println("response Body:", string(body))
	}
}
