# Logging in to Deezer using the mobile API

On the desktop versions of Deezer, logging in requires a Captcha. However, on
the mobile versions no Captcha is required. This is because on mobile the app
uses a different endpoint to log in, but encrypts the login parameters instead
(using a hardcoded key). This Gist details how to obtain the login parameter
encryption key (what I call the "gateway key") for the Android version.

Note that no keys will be posted here due to fear of DMCA takedowns. All keys
are easily obtainable by downloading the Android APK or iOS IPA and inspecting
the resources, and some keys can be obtained directly from the Deezer website
JS source code. How to do that is left as an exercise to the reader.

## Getting the API key

The API key is a plaintext string that's 64 characters long, and is sent with
every request. Some app requests are sent over plain HTTP so it should be
relatively easy to obtain them that way, otherwise you will have to scan the
binaries.

## Getting the gateway key (Android)

The gateway key is stored inside an icon asset within the APK.

1. Extract the file `assets/icon2.png` from the APK
2. Run the following Python script with the path to the icon as the first
   argument:

~~~py
import sys

MAGIC = b'PLTE'

fpath = sys.argv[1]
if not fpath:
   print('Usage: android_gw_key.py [path-to-icon2.png]\n')
   exit()

icon_buf = bytearray(open(fpath, 'rb').read())
magic_idx = icon_buf.find(MAGIC)
if magic_idx < 0:
   print('Magic phrase not found!')
   exit()

crypted = icon_buf[magic_idx:]
out = bytearray(300)

for i in range(len(out)):
   pm = i + len(MAGIC)
   if i >= 90:
       out[i - 90] = crypted[pm] ^ 0x40
   crypted[pm] = int((crypted[pm] + crypted[pm + 1] + crypted[pm + 2]) / 3)

print(out[1:17].decode('utf8'))
~~~

3. The key will be printed out.

## Getting the gateway key (iOS)

TBD...
