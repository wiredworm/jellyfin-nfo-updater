package main

import (
	"math"
	"sync"
	"fmt"
	"flag"
	"log"
	"time"
	"os"
	"strings"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"encoding/xml"
	"github.com/fsnotify/fsnotify"
)

const (
	// Header is a generic XML header suitable for use with the output of Marshal.
	// This is not automatically added to any output of this package,
	// it is provided as a convenience.
	Header = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
)

type Tvshow struct {
	XMLName       xml.Name `xml:"tvshow"`
	Plot          string   `xml:"plot"`
	Outline       string   `xml:"outline"`
	Lockdata      string   `xml:"lockdata"`
	Dateadded     string   `xml:"dateadded"`
	Title         string   `xml:"title"`
	Originaltitle string   `xml:"originaltitle"`
	Trailer       string   `xml:"trailer"`
	Rating        string   `xml:"rating"`
	Year          string   `xml:"year"`
	Mpaa          string   `xml:"mpaa"`
	ImdbID        string   `xml:"imdb_id"`
	Tmdbid        string   `xml:"tmdbid"`
	Premiered     string   `xml:"premiered"`
	Releasedate   string   `xml:"releasedate"`
	Enddate       string   `xml:"enddate"`
	Runtime       string   `xml:"runtime"`
	Genre         []string `xml:"genre"`
	Studio        string   `xml:"studio"`
	Tag           []string `xml:"tag"`
	Tvrageid      string   `xml:"tvrageid"`
	Tvdbid        string   `xml:"tvdbid"`
	Art           struct {
		Poster string `xml:"poster"`
		Fanart string `xml:"fanart"`
	} `xml:"art"`
	Actor []struct {
		Name      string `xml:"name"`
		Role      string `xml:"role"`
		Type      string `xml:"type"`
		Sortorder string `xml:"sortorder"`
		Thumb     string `xml:"thumb"`
	} `xml:"actor"`
	ID           string `xml:"id"`
	Episodeguide struct {
		URL  struct {
			Cache string `xml:"cache,attr"`
		} `xml:"url"`
	} `xml:"episodeguide"`
	Season  string `xml:"season"`
	Episode string `xml:"episode"`
	Status  string `xml:"status"`
} 

type Movie struct {
	XMLName       xml.Name `xml:"movie"`
	Plot          string   `xml:"plot"`
	Lockdata      string   `xml:"lockdata"`
	Dateadded     string   `xml:"dateadded"`
	Title         string   `xml:"title"`
	Originaltitle string   `xml:"originaltitle"`
	Director      string   `xml:"director"`
	Writer        []string `xml:"writer"`
	Credits       []string `xml:"credits"`
	Trailer       []string `xml:"trailer"`
	Rating        string   `xml:"rating"`
	Year          string   `xml:"year"`
	Mpaa          string   `xml:"mpaa"`
	Imdbid        string   `xml:"imdbid"`
	Tmdbid        string   `xml:"tmdbid"`
	Premiered     string   `xml:"premiered"`
	Releasedate   string   `xml:"releasedate"`
	Criticrating  string   `xml:"criticrating"`
	Runtime       string   `xml:"runtime"`
	Tagline       string   `xml:"tagline"`
	Country       []string `xml:"country"`
	Genre         []string `xml:"genre"`
	Studio        []string `xml:"studio"`
	Tag           []string `xml:"tag"`
	Art           struct {
		Poster string `xml:"poster"`
		Fanart string `xml:"fanart"`
	} `xml:"art"`
	Actor []struct {
		Name      string `xml:"name"`
		Role      string `xml:"role"`
		Type      string `xml:"type"`
		Sortorder string `xml:"sortorder"`
		Thumb     string `xml:"thumb"`
	} `xml:"actor"`
	ID       string `xml:"id"`
	Fileinfo struct {
		Streamdetails struct {
			Video struct {
				Codec             string `xml:"codec"`
				Micodec           string `xml:"micodec"`
				Bitrate           string `xml:"bitrate"`
				Width             string `xml:"width"`
				Height            string `xml:"height"`
				Aspect            string `xml:"aspect"`
				Aspectratio       string `xml:"aspectratio"`
				Framerate         string `xml:"framerate"`
				Language          string `xml:"language"`
				Scantype          string `xml:"scantype"`
				Default           string `xml:"default"`
				Forced            string `xml:"forced"`
				Duration          string `xml:"duration"`
				Durationinseconds string `xml:"durationinseconds"`
			} `xml:"video"`
			Audio struct {
				Codec        string `xml:"codec"`
				Micodec      string `xml:"micodec"`
				Bitrate      string `xml:"bitrate"`
				Language     string `xml:"language"`
				Scantype     string `xml:"scantype"`
				Channels     string `xml:"channels"`
				Samplingrate string `xml:"samplingrate"`
				Default      string `xml:"default"`
				Forced       string `xml:"forced"`
			} `xml:"audio"`
		} `xml:"streamdetails"`
	} `xml:"fileinfo"`
} 

var watcher *fsnotify.Watcher

func printTime(s string, args ...interface{}) {
	fmt.Printf(time.Now().Format("15:04:05.0000")+" "+s+"\n", args...)
}

func updateNfoFile(nfoFilePath string){
// This function will open the specified NFO file and modify the MPAA rating if it starts with GB-
// In situations where this is the case the GB- will be replaced with UK:

	// get just the filename
	filePathParts := strings.Split(nfoFilePath,"/")
	fileName := filePathParts[len(filePathParts)-1]

	// open the XML file
	xmlFile, err := os.Open(nfoFilePath)

	// catch and display any errors
	if err != nil {
		log.Println(err)
	} else {
		// if the file was opened successfully then defer the closing of our xmlFile
		// so that we can parse it later on
		defer xmlFile.Close()

		// read our opened xmlFile as a byte array.
		byteValue, err := ioutil.ReadAll(xmlFile)

		if err != nil{
			fmt.Println("An error occurred whilst reading the file.")
			fmt.Printf("%v\n",err)
		}

		if strings.ToLower(fileName) == "tvshow.nfo"{
			// unmarshal the file into the struct
			var tvshow Tvshow
			xml.Unmarshal(byteValue, &tvshow)
			
			// if the mpaa rating starts with GB- then...
			if strings.HasPrefix(tvshow.Mpaa,"GB-"){
				
				// note the original MPAA rateing
				origMPAA := tvshow.Mpaa

				// swap it for UK:
				tvshow.Mpaa = strings.ReplaceAll(origMPAA,"GB-","UK:")
				
				// marshall from XML back into a byte array
				file, _ := xml.MarshalIndent(tvshow, "", "  ")
				
				// write out the XML and the standard header back to the file
				_ = ioutil.WriteFile(nfoFilePath, []byte(Header + string(file)), 0644)

				log.Printf("  The MPAA rating in the file %s was changed from %s to %s.\n",fileName,origMPAA,tvshow.Mpaa)
			}
		} else {
			// unmarshal the file into the struct
			var movie Movie
			xml.Unmarshal(byteValue, &movie)
			
			// if the mpaa rating starts with GB- then...
			if strings.HasPrefix(movie.Mpaa,"GB-"){
				
				// note the original MPAA rateing
				origMPAA := movie.Mpaa

				// swap it for UK:
				movie.Mpaa = strings.ReplaceAll(origMPAA,"GB-","UK:")
				
				// marshall from XML back into a byte array
				file, _ := xml.MarshalIndent(movie, "", "  ")
				
				// write out the XML and the standard header back to the file
				_ = ioutil.WriteFile(nfoFilePath, []byte(Header + string(file)), 0644)

				log.Printf("  The MPAA rating in the file %s was changed from %s to %s.\n",fileName,origMPAA,movie.Mpaa)
			}
		}
	}
}

// main
func main() {

	var folders string

	// get the folders to be monitored that were specified with the -d option
	flag.StringVar(&folders, "d", "", "")
	flag.Parse()

	// if none were provided then exit with a suitable warning
	if strings.Trim(folders," ") == ""{
		log.Println("Please specify one or more folders to monitor by using the -d option.")
		log.Println("To specify multiple paths please ensure each is seperated by a comma.")
		os.Exit(0)
	} 

	// split the directories into a slice
	xDirs := strings.Split(folders,",")

	// record the current time so we can identify how long th einitial scan takes
	startTime := time.Now()
	log.Println("The nfo-renamer is starting.....")

	var doF = func(xpath string, xinfo fs.DirEntry, xerr error) error {

		// first thing to do, check error. and decide what to do about it
		if xerr != nil {
			log.Printf("error [%v] at a path [%q]\n", xerr, xpath)
			return xerr
		}

		// only process nfo files
		if !xinfo.IsDir() {
			fileExtension := strings.ToLower(filepath.Ext(xpath))
			if fileExtension == ".nfo"{
				nfoFilePath := fmt.Sprintf("%v/%v",filepath.Dir(xpath),xinfo.Name())
				updateNfoFile(nfoFilePath)
			}
		}
		return nil
	}

	// scan all files and folders under the starting path
	for _,v := range xDirs{
		log.Printf("  Processing the folder %s\n",v)
		err := filepath.WalkDir(v, doF)

		if err != nil {
			log.Printf("error walking the path %q: %v\n", v, err)
		}
		
		// display the time taken for the initial scan to complete
		elapsed := time.Since(startTime)
		log.Printf("  Scanning for nfo files took %s\n\n", elapsed)
	}

	// creates a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()	

	log.Println(" Creating the file system watchers")

	//err := watchDir(watcher *fsnotify.Watcher, xDir, reload func()) error {
	for _,v := range xDirs{
		watcher.Add(v)
	}

	// starting at the root of the project, walk each file/directory searching for
	// directories
	//if err := filepath.Walk("/home/pwes/code/nfo-updated/test", watchDir); err != nil {
	for _,v := range xDirs{
		if err := filepath.Walk(v, watchDir); err != nil {
			log.Println("ERROR", err)
		}
	}

	log.Println(" File system watchers now created")

	done := make(chan bool)

	//
	go func() {
		var (
			// Wait 100ms for new events; each new event resets the timer.
			waitFor = 100 * time.Millisecond
	
			// Keep track of the timers, as path → timer.
			mu     sync.Mutex
			timers = make(map[string]*time.Timer)
	
			// Callback we run.
			printEvent = func(e fsnotify.Event) {
				//printTime(e.String())

				// get file path in lower case
				filePath := strings.ToLower(e.Name)
					
				// if the update was to an nfo file
				if strings.HasSuffix(filePath,".nfo") {
					updateNfoFile(e.Name)
				}

				// Don't need to remove the timer if you don't have a lot of files.
				mu.Lock()
				delete(timers, e.Name)
				mu.Unlock()
			}
		)

		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				//log.Printf("EVENT! %#v\n", event)

				// get file path in lower case
				filePath := strings.ToLower(event.Name)
					
				// if the update was to an nfo file
				if strings.HasSuffix(filePath,".nfo") {
					if event.Has(fsnotify.Create){
						updateNfoFile(event.Name)
					}
				} else {
					if event.Has(fsnotify.Create){
						fileInfo, err := os.Stat(event.Name)
						if err == nil{
							//log.Printf("EVENT! %#v\n", event)
							if fileInfo.IsDir(){
								watcher.Add(event.Name)
							}
						}
					}
				}

							// Get timer.
				mu.Lock()
				t, ok := timers[event.Name]
				mu.Unlock()

				// No timer yet, so create one.
				if !ok {
					t = time.AfterFunc(math.MaxInt64, func() { printEvent(event) })
					t.Stop()

					mu.Lock()
					timers[event.Name] = t
					mu.Unlock()
				}

				// Reset the timer for this path, so it will start from 100ms again.
				t.Reset(waitFor)

				// watch for errors
			case err := <-watcher.Errors:
				log.Println("ERROR", err)
			}
		}
	}()

	<-done
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}