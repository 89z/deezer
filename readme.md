# Deezer

> Listen, um, I don’t really know what’s going on here, but I’m leaving. I
> don’t know where exactly, but I’m gonna start over.
>
> Come with me. Please.
>
> [Paint it Black (2016)](//wikipedia.org/wiki/Paint_It_Black_%282016_film%29)

Download audio from Deezer

## API

GetData:

Given a SNG_ID, return Results.Data

GetSource:

Given SNG_ID, Results.Data and quality, return audio URL

Decrypt:

Given SNG_ID and byte slice, decrypt byte slice in place

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
