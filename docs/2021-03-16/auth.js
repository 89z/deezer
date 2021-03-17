// Calling this function:

// mobileClient.signInWithEmail(

// REQUEST
{
  method: 'get',
  url: 'https://api.deezer.com/1.0/gateway.php?method=mobile_auth&api_key=4VCYIJUCDLOUELGD1V8WBVYBNVDYOXEWSLLZDONGBBDFVXTZJRXPR29JRLQFO6ZE&output=3&uniq_id=30ff2d2e124932f8d51d601eef91eb5a',
  headers: {
    'User-Agent': 'Deezer/6.1.22.49 (Android; 9; Tablet; us) innotek GmbH VirtualBox'
  }
}

// RESPONSE
{
  error: [],
  results: {
    LOGIN_TEXT: 'Listen to music. Try Flow.',
    TOKEN: '75529d6ee8c8192c8b7ed46aad22de67fba6be11f9161afdf8c2504179b42a561ce7e00244ecf3567fb52aa4a9c6994f21a2c7b4e5f7d69545963012a11506f9b86ccaa4b0b7555406fb85ffb5cf40f53b01e7e9fbdd28baeaf43508c8fb7a5950cd44056c365701f69af35276457962',
    POLICY: {
      S_MOD: '1',
      S_SMARTRADIO: '1',
      S_RADIO: '1',
      S_STREAMING: '0',
      S_PREMIUM: '0',
      S_ALC: '0'
    },
    TIMESTAMP: 1615945802,
    PLATFORM: {
      name: 'android',
      family: 'tablet',
      electron: false,
      os_version: '9',
      os_name: 'android',
      version: 2,
      ua: 'Deezer/6.1.22.49 (Android; 9; Tablet; us) innotek GmbH VirtualBox',
      model: 'innotek GmbH VirtualBox',
      app_version: '6.1.22.49',
      app_name: 'Deezer',
      external_webview: null,
      lang: 'us',
      browser: false
    },
    COUNTRY: 'US',
    PHONE_COUNTRY: {
      LABEL: 'United States of America',
      PHONE_CODE: '+1',
      COUNTRY_ISO: 'US'
    },
    CONFIG: {
      TOTAL_TRACKS: '73',
      HOST_STREAM_MOB: 'http://e-cdn-proxy-{0}.deezer.com/mobile/1/',
      HOST_STREAM_PREVIEW: 'http://e-cdn-preview-{0}.deezer.com/stream/{1}-{2}.mp3',
      HOST_CUSTOS_MOB: 'http://e-cdn-content.dzcdn.net/mobile/custos/',
      HOST_API_MOB: 'https://api.deezer.com/1.0/gateway.php',
      HOST_FILES: 'http://e-cdn-files.dzcdn.net',
      HOST_CONTENT: 'http://e-cdn-content.dzcdn.net',
      HOST_IMAGES: 'http://cdn-images.dzcdn.net/images',
      HOST_LIVE: 'http://live.deezer.com',
      HOST_MESSAGING: 'messaging.deezer.com',
      HOST_SITE: 'https://www.deezer.com',
      HOST_UPLOAD: 'https://upload.deezer.com',
      URL_REGISTER: 'https://www.deezer.com/register.php',
      URL_FACEBOOK_LOGIN: 'https://www.deezer.com/facebook.php?token=%@&fb_user_id=%@&from_mobile=1',
      URL_PASSWORD_RESET: 'https://www.deezer.com/password/reset',
      URL_SUPPORT: 'http://support.deezer.com',
      URL_PAYMENT: 'https://m.deezer.com/payment/?webview=true',
      URL_MEDIA: 'https://media.deezer.com',
      GAIN: { TARGET: '-15', ADS: '-13' },
      URL_VENDOR_LIST: 'https://www.deezer.com/us/legal/vendor-list',
      GZIP: false,
      POP: 'fr',
      FB_ONLY: 0,
      FB_PUBLISH: 0,
      RESTRICTED: false,
      PLAYLIST_WELCOME_ID: '1014073671',
      REGISTRATION: {},
      POLICY: 'MOBILE',
      APPSTUDIO: false,
      PAYMENT_METHODS: {
        VISA: true,
        MASTERCARD: true,
        AMEX: true,
        PAYPAL: true,
        IOS: true
      },
      AD_RULES: {
        audio: { default: { start: 1, interval: 3, unit: 'track' } },
        display: { interstitial: { start: 900, interval: 900, unit: 'sec' } },
        big_native_ads_home: {
          iphone: { enabled: false },
          ipad: { enabled: false },
          android: { enabled: false },
          android_tablet: { enabled: false }
        }
      },
      START_IMAGE: { TYPE: 'misc', MD5: '3f70c0d78d6df6c16f1f31c171096803' },
      REG_FLOW: [ 'email', 'msisdn' ],
      LOGIN_FLOW: [ 'email', 'msisdn' ],
      VOICECALLBACK: false
    },
    USER_TOKEN: 'b609152d2f84159094f6b5d05078a1674462c8443267b211086699d70c35ecfd',    OFFER: {},
    LOGIN_ART: [
      '64ff3c74610e2e6ed60b058d99605ee6',
      '0cf38443f4108ca1c9043ade46c85f2a',
      'f0ad5828166c8f61bae16f91d21797d5',
      '928f6f7822b448def8ddbcd8b39955c7'
    ],
    LOGIN_ART_NEW: [
      '40bc81c691ef92e49f9f22518b4d269b',
      '84e5addaaecfb5611539c0415a4689d6',
      'bff018772fb06885e3ee938bd6cbbc4c',
      '6f5acb294b6ffb18a7024b525bb7424c'
    ],
    LOGIN_ART_0115: [ '2747dfc240d63a71e4c5256de92fa7e0' ],
  }
}

// REQUEST
{
  method: 'get',
  url: 'https://api.deezer.com/1.0/gateway.php?method=api_checkToken&api_key=4VCYIJUCDLOUELGD1V8WBVYBNVDYOXEWSLLZDONGBBDFVXTZJRXPR29JRLQFO6ZE&output=3&auth_token=a1b857c5d8859d331c77e78e6ab00641e8c9a23e1c84e2774f3857ad12587a5392317e8514732be3ce34459b29870f440f9932964c271bed20156e7d61901ce2',
  headers: {
    'User-Agent': 'Deezer/6.1.22.49 (Android; 9; Tablet; us) innotek GmbH VirtualBox'
  }
}

// RESPONSE
{ error: [], results: 'frd0d4e91d686ca6b4dac9098297436574e9853f' }

// REQUEST
{
  method: 'post',
  url: 'https://api.deezer.com/1.0/gateway.php?method=mobile_userAuth&api_key=4VCYIJUCDLOUELGD1V8WBVYBNVDYOXEWSLLZDONGBBDFVXTZJRXPR29JRLQFO6ZE&output=3&input=3&sid=frd0d4e91d686ca6b4dac9098297436574e9853f',
  headers: {
    'User-Agent': 'Deezer/6.1.22.49 (Android; 9; Tablet; us) innotek GmbH VirtualBox'
  },
  data: {
    mail: 'srpen6@gmail.com',
    password: '1668557f96c7dc0f3027324beb59f347',
    device_serial: '',
    platform: 'innotek GmbH_x86_64_9',
    custo_version_id: '',
    custo_partner: '',
    model: 'VirtualBox',
    device_name: 'VirtualBox',
    device_os: 'Android',
    device_type: 'tablet',
    google_play_services_availability: '1',
    consent_string: ''
  }
}

// RESPONSE
{
  error: [],
  results: {
    USER_ID: 4239375622,
    PREMIUM: {
      OFFER_ID: '1',
      DATE_START: '2021-03-16 21:28:13',
      DATE_END: '2021-06-16 21:28:13',
      PLATEFORM: 'ogone',
      TRY_AND_BUY: '1',
      COMMITMENT: '0',
      RENEWAL: '0',
      STATUS: true,
      RANDOM: '6e605b0fbbc2cf8a0e1c1de5ade9fc61e67583476f0fb4d261de7b9731e492030c247407bf4a2ae1329c7e85e229cf7b79d00412bc29a3bcd93775377c1423786ae318701ca445a69fa4d202ba4e3a00c0cfdab80aa9d00a0c0d795b3c071f953d3bbf0e6214b070be97850008e34f69d67f94dc9a3eafac65575a6dec3cd04adc6cc7f7b0aac3aefa6b15bd2bf98c429df8f19913a103d3a05cc815702560db9bf0973c4d1f5449f421ded97bc801b24c0ddc2c7bb5de2e0e460c13a1e4384186365578144b194f4999729632b8b438eca92559e90588fca9eec46cee35fd7a004ea825b97c19d141bc54795a1b46f1a68925d7d94f07ce05eb6a7e87216c8fd6871ba272254bb897c2abf40076d1ae57df552593f58fb9fe04f67c9ef718e2b89d0b838943dd609a5699cc3e028b67dca0420c3c64bda41c605fa55a603f44b329f504234c11afaddd8d9fec03d952db2a67fa69ca09640006d57105bbe51a5a15b0c66ac516ca48cb34bed87f98e7bb45433da76b73169288a88173760a2a224e054a350d951ffd76a0c2597546cb1c08e2a49016b63a03fce0d14e3c979c626feb390967fc63c8488a0a46141dde8b22b8f2d49834a09a4a90519f9f6d75b4c2c34a0cf354a1c53d185a6ef620dcd0d9bea2c844461f1f23e84e92d652323448bbbaa47a41be737f73073478f443431eeee419cd5f4573d5cfdb34e15f21e6e88868e02e01e269cc94408a393e12671666e0be6511e721b153013a5ae763677dca5a107b57406ea47c5beec7287aff36afef5adcea769a5d29d21334c0d95827571a29a59ae2487d5a2261f83272c7f7ff12d1b88e501ca46fc2c1e58fac1c0aeb26ab030d2cf1144a0fa47597a32773158d7a9394209af7265c58dd593e25d94f61e3bd376881f15b40b16e0905d83f7de5b51d2391f4ffde51a422f43e411d4595f4f1996b0bc574ef8d2c74556b030639ac143e598da34f7a870123829c97cd8805a5861e2dcd3025dcd6b02981ac469b489d17f36cf1bd882a420f2d8c95c71b551e309b95495e341501654e2c6a81d2d6db071b3c0ca3c783af57cfe503d5470074129a7169ca904a2784bfad78c72b03f7a08753c3512d9f1218254f59a491877df9417ff8d0b92eb22c5389cbede155ba9dbe0a1c4a4b7a68011abc8ac7a49ea17298d198e3563fe36a6336c534f749c0cec13fe2b6dc2a6444987bf7006874f3cbac04a249e6fa01e89e4185d446b1414787e5b7ebb08378e207014319d1e4502a11edcc98d1a9395ac07705a9cd47a73d7e905d416cad1fc03d5d748993195d9ea3cfcf3f011d0502c12a3e1b26775cc15341cd79855da2397cdea0d0f3ba8d0f731fc3e64274cacc8ad2e6358e268f45ef2706ad3487bc3bc057377447ca0ecd6d7da45493435d49efd4bd05a11bab98d7b17b558c1911860f118977defd2705286f04367fe1f2141f57b216dad8e292c41f5a4fabae88143fad15d410efc053657c904bd58df2d3760abf2bd3e606b910385e9ba5087ff5e979a5e4a8c104bbd29881f7d1aad94075833cddaf8c107eef578d510b3bc2853faa99b2aa5f4c7bfe39ee00b2ad38f729',
      OFFER_NAME: 'Deezer Premium',
      IS_B2B: false,
      TRY_AND_BUY_CC_DONE: true
    },
    COUNTRY: 'US',
    RECOMMENDATION_COUNTRY: 'US',
    ZONE: 'NORTH_AMERICA',
    EMAIL: 'srpen6@gmail.com',
    SEX: 'M',
    FB_USER_ID: '',
    FB_TOKEN: '',
    GG_USER_ID: '',
    GG_TOKEN: '',
    TWITTER_TOKEN: '',
    TWITTER_SECRET: '',
    USER_TOKEN: 'f532e274bf966c6e296c76ad0c02814b52b4238c266fc2f26e27150e4e4c1ad1',    NOTIFICATION_UNREAD: 3,
    EXPLICIT_CONTENT_LEVEL: 'explicit_display',
    EXPLICIT_CONTENT_LEVELS_AVAILABLE: [
      'explicit_display',
      'explicit_no_recommendation',
      'explicit_hide'
    ],
    CAN_EDIT_EXPLICIT_CONTENT_LEVEL: true,
    FIRSTNAME: '',
    LASTNAME: '',
    BLOG_NAME: 'srpen6',
    USER_PICTURE: '',
    AGE: 38,
    IS_KID: false,
    USER_PHONE: {},
    ARL: 'a9e22d64eba788f932c7624e06df78f0fe5d3adbf1417563c5e290a2e2febcf0e955b5e464b52814b2cbaefdf7c260b5f1401b14ef5384f1efaa02a341ac35465d7ebec81b87cf2f7cc855d3a580fc19320a243560e9f5642dea485dc7c8b45c',
    MULTI_ACCOUNT: {
      enabled: false,
      active: false,
      max_children: null,
      parent: null,
      is_sub_account: false
    },
    CAN_BE_CONVERTED_TO_INDEPENDENT: false,
    ABTEST: {
      discovery_algorithms: {
        id: 'discovery_card',
        option: 'svd_discovery',
        behaviour: 'new discovery algorithm, based on nearest neighbours search'
      },
      flow_algorithms: {
        id: 'Flow-2021-february',
        option: 'C',
        behaviour: 'default algorithm, using v2 tags for reco, all tags'
      },
      showstab_newbie: {
        id: 'showstab_newbie',
        option: 'default',
        behaviour: 'default showstab'
      },
      triforce_queuelist_ui: {
        id: 'triforce_queuelist_ui',
        option: 'default',
        behaviour: 'Cloud based queue list deactivated'
      },
      android_exoplayer2: {
        id: 'android_exoplayer2',
        option: 'A',
        behaviour: 'Legacy audio engine'
      },
      playlists_suggestion_algorithm_october_2020: {
        id: 'playlists_suggestion_algorithm_october_2020',
        option: 'default',
        behaviour: 'Suggest playlists thompson sampling reinforcement learning, fallback on coldstart models'
      },
      share_android_image_preview: {
        id: 'share_android_image_preview',
        option: 'default',
        behaviour: 'default sharing'
      },
      premium_tab_direct_link: {
        id: 'premium_tab_direct_link',
        option: 'direct_link_to_payment',
        behaviour: 'Direct link to payment'
      }
    },
    DEVICE_ALREADY_LINKED: '0',
    MAX_NB_DEVICES: '3',
    AUDIO_SETTINGS: {
      default_preset: 'fast',
      default_download_on_mobile_network: false,
      presets: [
        {
          mobile_download: 'standard',
          mobile_streaming: 'standard',
          wifi_download: 'standard',
          wifi_streaming: 'standard',
          id: 'fast',
          title: 'Compact',
          description: 'Adjust audio quality to reduce data usage.'
        },
        {
          mobile_download: 'standard',
          mobile_streaming: 'standard',
          wifi_download: 'high',
          wifi_streaming: 'high',
          id: 'balanced',
          title: 'Balanced',
          description: 'Improve audio quality when connected to WiFi.'
        },
        {
          mobile_download: 'high',
          mobile_streaming: 'high',
          wifi_download: 'high',
          wifi_streaming: 'high',
          id: 'quality',
          title: 'Better',
          description: 'Best audio quality available.'
        }
      ],
      connected_device_streaming_preset: 'high'
    },
    AD_CONFIGS: {
      audio: { default: { start: 1, interval: 3, unit: 'track' } },
      display: { interstitial: { start: 900, interval: 900, unit: 'sec' } }
    },
    DESCRIPTION: 'srpen6@gmail.com',
    PROFIL: { ROAMING: '0', NETWORK: '-', ACTIVE: {}, CUSTOM: {} },
    CUSTO: {
      DATA: {
        id: '1',
        name: 'premiumplus_FR',
        lang: {
          locale: [
            'en',    'fr', 'nl',
            'ru',    'pt', 'es',
            'it',    'de', 'pl',
            'pt-br', 'ro', 'th',
            'ms',    'id', 'hu',
            'tr'
          ],
          default_locale: 'en'
        },
        buy_itunes: { activated: 'true' },
        startup_push: {
          action: [
            {
              condition: 'TRY_AND_BUY_MOB_OLD',
              priority: '1',
              type: 'POPUP',
              popup: { text_id: 'tryandbuy.mobile' }
            },
            {
              condition: 'android_generic_3.000',
              priority: '2',
              type: 'POPUP',
              popup: {
                text_id: 'update',
                button: {
                  text_id: 'update.ok',
                  action: [ { condition: 'NONE', type: 'CLOSE_POPUP' }, {} ]
                }
              }
            },
            {
              condition: 'android_generic_3.001',
              priority: '3',
              type: 'POPUP',
              popup: {
                text_id: 'update',
                button: {
                  text_id: 'update.ok',
                  action: [ { condition: 'NONE', type: 'CLOSE_POPUP' }, {} ]
                }
              }
            },
            {
              condition: 'android_generic_3.002',
              priority: '4',
              type: 'POPUP',
              popup: {
                text_id: 'update',
                button: {
                  text_id: 'update.ok',
                  action: [ { condition: 'NONE', type: 'CLOSE_POPUP' }, {} ]
                }
              }
            },
            {}
          ]
        }
      },
      VERSION_KEY: 'WW_1__18'
    },
    CUSTO_CONDITIONS: {
      STORE_DIRECT: true,
      _: true,
      TRYNBUY_DONE: true,
      PARTNER_DISPLAY_EXTERNAL_CONTENT_WARNING: false
    },
    DEFAULT_PLAYLIST_ID: '8796668262',
    APPSTUDIO: { apps: true, store: true },
    WAITING_FOR_PAYMENT: false,
    CUSTOM_TEXTS: [],
    OFFER_ELIGIBILITIES: [
      {
        type: 'subscription',
        offer: {
          id: 100001,
          name: 'Deezer Family',
          description: [
            '6 Deezer Premium profiles, 6 happy people, just $14.99/month',
            "The best part is, there's no contract!"
          ],
          duration: 1,
          price: { amount: '14.99', currency: 'USD', display: '$14.99' },
          multi_account_max_allowed: 6
        },
        cta: {
          label: 'Subscribe to Deezer Family',
          label_extend: 'Subscribe for $14.99/month ',
          log_name: 'subscription',
          subscribe_uri: 'https://www.deezer.com/payment/go.php?cip=NDVlOWIxNmViNmIzNWU1YjRlMDljMmI3ZjEzOGFiZjc5MmZjZmFiZjJiODNkOWI3MGEyNGI4M2FmMDhhZmY0ZDJmOWY2MTFlMjhkY2IyZmNkNGZlZmRhZjVjMGU5YzQyYzMxNTdlMGJiMDE4OGM5YWIxODM5YjgyODkwZTMwYzllZDY5MzY0YzRmM2E5NGU4MTc0MmYwNTZmNzkyMWVjMw%3D%3D&webview=1'
        },
        category: 'deezer'
      }
    ],
    CONVERSION_ENTRYPOINTS: {
      SETTINGS_MYACCOUNT: {
        deeplink: 'deezer://www.deezer.com/offerwall/WANT_TO_SUBSCRIBE_FROM_SETTINGS_MYACCOUNT',
        cta_label: 'Change plan',
        offer_type: 'FAMILY'
      },
      SETTINGS_AUDIO_HIFI: {
        deeplink: 'deezer://www.deezer.com/offerwall/SETTINGS_AUDIO_HIFI?focus=hifi',
        cta_label: 'Get High Fidelity sound with Deezer HiFi',
        description: '$14.99/month TRY IT NOW'
      }
    },
    OFFER_BOX_URL: 'https://www.deezer.com/us/paywalls?app_version=6.1.22.49&ts=1615945802',
    APPCUSTO_CHECKSUM: '20688b10ac10a8cdc38ca066f9373222',
    CHECKSUMS: { OFFERWALL: 'd0abd39bca578263debb2f50dfb27ff6' },
    SHOW_OPTIN_BUTTON: {
      URL: 'https://www.deezer.com/us/optin-box-inapp',
      LABEL: 'Notification preferences'
    },
    HAS_CONSENT: false
  }
};
