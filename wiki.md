# Wiki

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
