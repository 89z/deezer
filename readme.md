# Deezer

Download audio from Deezer

https://www.deezer.com/us/track/75498415

## API

getData:

given a SNG_ID, return Results.Data

getSource:

given SNG_ID, Results.Data and Quality, return Audio URL

decryptAudio:

given SNG_ID and Audio URL, return io.Reader

https://github.com/godeezer/lib/blob/8e011608/deezer/crypto.go#L87-L97

## Thanks

- https://github.com/moon004/Go-deezer-downloader
- https://pkg.go.dev/github.com/godeezer/lib/deezer
- https://pkg.go.dev/github.com/joshbarrass/deezerdl/pkg/deezer

## AES

Used for URL Generation

Go to the "source" panel of browser debugger, to show the deezer [webworker
source code](https://imgur.com/pwS370Q.png). [pretty
print](https://i.imgur.com/P3eaAf3.png) the source. put a breakpoint on the
[latest occurrence of the "JSON" word](https://i.imgur.com/E0UJDwX.png). Play a
new track, to trigger the breakpoint You see the keys in the [the current
scope](https://i.imgur.com/1e6P98L.png).

~~~js
[...key].map(e=>String.fromCharCode(e)).join('');
~~~

## CBC

Used for track deciphering

Same method as the [DZR_AES](#DZR_AES), but the CBC keys are [split in 2
arrays](https://i.imgur.com/1e6P98L.png) that need to be reassembled using the
following script

~~~js
let a=[97, ... 103];
let b=[49, ... 52];

[].
concat(...a.map((c,i)=>[a[i],b[i]]).reverse()).
map(e=>String.fromCharCode(e)).
join(''):
~~~
