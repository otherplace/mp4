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
$./mp4tool info ../sample/meta.test1.mp4 
{"ftyp":{"MajorBrand":"isom","MinorVersion":512,"CompatibleBrands":["isom","iso2","avc1","mp41"]},"moov":{"mvhd":{"Version":0,"Flags":[0,0,0],"CreationTime":0,"ModificationTime":0,"Timescale":1000,"Duration":2024,"NextTrackId":0,"Rate":65536,"Volume":256},"Iods":null,"trak":[{"tkhd":{"Version":0,"Flags":[0,0,3],"CreationTime":0,"ModificationTime":0,"TrackId":1,"Duration":2000,"Layer":0,"AlternateGroup":0,"Volume":0,"Matrix":"AAEAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAABAAAAA","Width":251658240,"Height":112328704},"mdia":{"mdhd":{"Version":0,"Flags":[0,0,0],"CreationTime":0,"ModificationTime":0,"Timescale":12288,"Duration":24576,"Language":5575},"hdlr":{"Version":0,"Flags":[0,0,0],"PreDefined":0,"HandlerType":"vide","Name":"VideoHandler\u0000"},"minf":{"vmhd":{"Version":0,"Flags":[0,0,1],"GraphicsMode":0,"OpColor":[0,0,0]},"smhd":null,"stbl":{"Stsd":{"Version":0,"Flags":[0,0,0]},"Stts":{"Version":0,"Flags":[0,0,0],"SampleCount":[48],"SampleTimeDelta":[512]},"Stss":{"Version":0,"Flags":[0,0,0],"SampleNumber":[1]},"Stsc":{"Version":0,"Flags":[0,0,0],"FirstChunk":[1,2],"SamplesPerChunk":[2,1],"SampleDescriptionID":[1,1]},"Stsz":{"Version":0,"Flags":[0,0,0],"SampleUniformSize":0,"SampleNumber":48,"SampleSize":[131933,81064,26236,21269,19355,66845,23874,19672,20042,53964,22727,18513,18729,63390,23401,19837,20071,67166,28061,21088,20469,58968,23139,19383,16762,60693,22502,17759,17155,48588,21095,14786,15184,55622,21470,17784,15677,42429,19281,14030,14648,44803,22681,15870,16329,35535,17118,16898]},"Stco":{"Version":0,"Flags":[0,0,0],"ChunkOffset":[48,213335,240186,262080,282106,349672,373908,394315,415079,469755,493235,512113,531588,595720,619873,640451,660902,728830,757610,779441,800655,860003,883899,904029,921482,982913,1005772,1024267,1042234,1091531,1113320,1128446,1144378,1200794,1223027,1241516,1257562,1300738,1320799,1335555,1350568,1396100,1419568,1436161,1453218,1489120,1506989]},"Ctts":{"Version":0,"Flags":[0,0,0],"SampleCount":[1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,2],"SampleOffset":[1024,2560,1024,0,512,2560,1024,0,512,2560,1024,0,512,2560,1024,0,512,2560,1024,0,512,2560,1024,0,512,2560,1024,0,512,2560,1024,0,512,2560,1024,0,512,2560,1024,0,512,2560,1024,0,512,2048,512]}},"dinf":{"Dref":{"Version":0,"Flags":[0,0,0]}},"Hdlr":null}},"edts":{"elst":{"Version":0,"Flags":[0,0,0],"SegmentDuration":[2000],"MediaTime":[1024],"MediaRateInteger":[1],"MediaRateFraction":[0]}}},{"tkhd":{"Version":0,"Flags":[0,0,3],"CreationTime":0,"ModificationTime":0,"TrackId":2,"Duration":2024,"Layer":0,"AlternateGroup":1,"Volume":256,"Matrix":"AAEAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAABAAAAA","Width":0,"Height":0},"mdia":{"mdhd":{"Version":0,"Flags":[0,0,0],"CreationTime":0,"ModificationTime":0,"Timescale":44100,"Duration":89224,"Language":5575},"hdlr":{"Version":0,"Flags":[0,0,0],"PreDefined":0,"HandlerType":"soun","Name":"SoundHandler\u0000"},"minf":{"vmhd":null,"smhd":{"Version":0,"Flags":[0,0,0],"Balance":0},"stbl":{"Stsd":{"Version":0,"Flags":[0,0,0]},"Stts":{"Version":0,"Flags":[0,0,0],"SampleCount":[87,1],"SampleTimeDelta":[1024,136]},"Stss":null,"Stsc":{"Version":0,"Flags":[0,0,0],"FirstChunk":[1,2,6,7,11,12,16,17,21,22,26,27,31,32,36,37,40,41,45,46,47],"SamplesPerChunk":[1,2,1,2,1,2,1,2,1,2,1,2,1,2,1,2,1,2,1,2,6],"SampleDescriptionID":[1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1]},"Stsz":{"Version":0,"Flags":[0,0,0],"SampleUniformSize":0,"SampleNumber":88,"SampleSize":[290,322,293,306,319,335,336,370,351,362,372,363,362,360,354,358,371,382,365,348,398,347,395,384,368,379,362,380,382,380,359,360,376,367,375,370,380,383,374,381,366,348,343,382,356,357,372,364,411,401,371,338,363,331,340,351,397,393,401,385,378,368,337,369,371,376,396,384,366,360,365,346,383,389,398,359,364,362,366,367,373,378,368,402,384,398,484,7]},"Stco":{"Version":0,"Flags":[0,0,0],"ChunkOffset":[213045,239571,261455,281435,348951,373546,393580,414357,469043,492482,511748,530842,594978,619121,639710,660522,728068,756891,778698,799910,859623,883142,903282,920791,982175,1005415,1023531,1041422,1090822,1112626,1128106,1143630,1200000,1222264,1240811,1257193,1299991,1320019,1334829,1350203,1395371,1418781,1435438,1452490,1488753,1506238,1523887]},"Ctts":null},"dinf":{"Dref":{"Version":0,"Flags":[0,0,0]}},"Hdlr":null}},"edts":{"elst":{"Version":0,"Flags":[0,0,0],"SegmentDuration":[2000],"MediaTime":[1024],"MediaRateInteger":[1],"MediaRateFraction":[0]}}}]},"mdat":{"ContentSize":1525882},"free":[{}]}

$./mp4tool info -t ../sample/meta.test1.mp4 
Box type: ftyp
+- Major brand: isom
+- Minor version: 0x200
+- Compatible brands: sizes = 4
 +- [0] : isom
 +- [1] : iso2
 +- [2] : avc1
 +- [3] : mp41
Movie Header:
 Timescale: 1000 units/sec
 Duration: 2024 units (2s)
 Rate: 1.0
 Volume: 1.0
Track 0
Track Header:
 Duration: 2000 units
 WxH: 3840.0x1714.0
Segment Duration:
 #0: 2000 units
Media Header:
 Timescale: 12288 units/sec
 Duration: 24576 units (2s)
Sample to Chunk:
 #0 : 2 samples per chunk starting @chunk #1 
 #1 : 1 samples per chunk starting @chunk #2 
Time to sample:
 #0 : 48 samples with duration 512 units
Samples : 48 total samples
Key frames:
 #0 : sample #1
Chunk byte offsets:
 #0 : starts at 48
 #1 : starts at 213335
 #2 : starts at 240186
 ...
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
