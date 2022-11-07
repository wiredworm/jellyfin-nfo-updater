# JellyFin NFO File Updater
JellyFin is a media management and streaming solution which includes a huge number of great features, including the ability to automatically download media info and use it to generate .nfo files.

However, in my case I hit an issue whereby my Kodi installations wouldn't display the correct images for the MPAA studio ratings.  It turns out that this is because Kodi expects UK ratings to be formatted as UK:[Rating] (e.g. UK:PG, UK:18 etc) but JellyFin creates them in the format GB-[Rating] (e.g. GB-PG, GB-18).

This small utility program will start by updating any nfo files which it finds that are incorrectly formatted.  Once done it creates notify watchers so that it will trigger the same check to be performed on any new nfo files which get created.
