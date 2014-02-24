=======
gpxcat
======

Gpxcat is an utility for concatinating, simplifing and splitting GPX tracks.

## Remove all segments and split the tracks

In this example all tracks are concatinated, individual segments removed and the track is split.
The result of this operation is output to the "out" directory.
By default the tracks are split by day, according to the -format parameter.

```bash
./gpxcat -split files/*.gpx
```

## Resample all tracks from Colombia 

In this example all the tracks from one country are concatinated and resampled.
By default 500 points are kept, you can alter this using the -keep N flag.

```bash
./gpxcat -resample -keep 400 out/colombia/*.gpx > colombia-small.gpx
```