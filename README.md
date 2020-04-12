# mp4

A encoder/decoder class, io.Reader and io.Writer compatible, usable for HTTP pseudo streaming

For the complete MP4 specifications, see https://standards.iso.org/ittf/PubliclyAvailableStandards/c068960_ISO_IEC_14496-12_2015.zip 


## Warning

Some boxes can have multiple formats (ctts, elst, tkhd, ...). Only the version 0 of those boxes is currently decoded (see https://github.com/jfbus/mp4/issues/7).
Version 1 will be supported, and this will break a few things (e.g. some uint32 attributes will switch to uint64).

## CLI

A CLI can be found in cli/mp4tool.go

It can :

* Display info about a media
```
mp4tool info file.mp4
```
* Copy a video (decode it and reencode it to another file, useful for debugging)
```
mp4tool copy in.mp4 out.mp4
```
* Generate a clip
```
mp4tool clip --start 10 --duration 30 in.mp4 out.mp4
```

(if you really want to generate a clip, you should use ffmpeg, you will ge better results)

## LICENSE

See LICENSE
