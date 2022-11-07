# JellyFin NFO File Updater
JellyFin is a media management and streaming solution which includes a huge number of great features, including the ability to automatically download media info and use it to generate .nfo files.

However, in my case I hit an issue whereby my Kodi installations wouldn't display the correct images for the MPAA studio ratings.  It turns out that this is because Kodi expects UK ratings to be formatted as UK:[Rating] (e.g. UK:PG, UK:18 etc) but JellyFin creates them in the format GB-[Rating] (e.g. GB-PG, GB-18).

This small utility program will start by updating any nfo files which it finds that are incorrectly formatted.  Once done it creates notify watchers so that it will trigger the same check to be performed on any new nfo files which get created.

## Usage
In my configuration most of my media software is running as Docker containers directly on the Synology NAS.  The easiest configuration for me was as follows.

1. Copy the binary named jellyfin-nfo-updater to the NAS using either a mountpoint under WSL or a utility such as WinSCP.
2. Make sure you place the binary under /usr/sbin.  You might need to use sudo to move it here.  Also make sure the chmod settings have it defined as executable.
3. Create a new Triggered Task in the Synology Task Scheduler.  It should execute at boot-up and the command line should look something like this:
```
/usr/sbin/jellyfin-nfo-updater -d '/volume1/media/Movies,/volume1/media/TV Shows'
```
4.  Modify the above command line based on the folders you need monitoring.
5. Finally reboot the NAS and when it starts the process to do the renaming should be running in the background (you can verify this by connecting via SSH to the NAS and doing a 'ps -ef | grep jelly' command.)
