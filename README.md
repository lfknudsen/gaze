# Gaze

List all video files in the current directory, displaying their resolution,
bit-rate, and frame-rate.

## Dependencies

Requires `ffprobe` (i.e. `ffmpeg`) installed on your PATH.

## Compile

```
go build
```

## Usage

```
gaze [OPTIONS]

Options:
  -c    run concurrently (default true)
  -h    print this help message and exit
  -s string
        how to sort the output. Options are:
         'na': By name, ascending.
         'nd': By name, descending.
         'fa': By fps, ascending.
         'fd': By fps, descending.
         'ba': By bitrate, ascending.
         'bd': By bitrate, descending.
         'ra': By resolution, ascending.
         'rd': By resolution, descending.
         'wa': By width, ascending.
         'wd': By width, descending.
         'ha': By height, ascending.
         'hd': By height, descending.
         (default "na")
```

## Author

lfknudsen