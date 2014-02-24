=======
gpxcat
======

Gpxcat is an utility for concatinating, simplifing and splitting GPX tracks.

# Examples

These are examples of how to use this tool.

## Remove all segments and split the tracks

In this example all tracks are concatinated, individual segments removed and the track is split.
The result of this operation is output to the "out" directory.

```bash
./gpxcat -split files/*.gpx
```

## Resample all tracks from Colombia 

In this example all the tracks from one country are concatinated and resampled.

```bash
./gpxcat -resample -keep 400 out/colombia/*.gpx > colombia-small.gpx
```