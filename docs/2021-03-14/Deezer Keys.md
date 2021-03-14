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
2. Run the Python script with the path to the icon as the first argument
3. The key will be printed out.

## Getting the gateway key (iOS)

TBD...
