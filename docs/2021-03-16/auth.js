signInWithEmail BEGIN
REQUEST {
  method: 'post',
  url: 'https://api.deezer.com/1.0/gateway.php?method=mobile_userAuth&api_key=4VCYIJUCDLOUELGD1V8WBVYBNVDYOXEWSLLZDONGBBDFVXTZJRXPR29JRLQFO6ZE&output=3&input=3&sid=fr737eb0c0365971df92b1d18554382a04ab95b1',
  headers: {
    'User-Agent': 'Deezer/6.1.22.49 (Android; 9; Tablet; us) innotek GmbH VirtualBox'
  },
  data: {
    mail: 'srpen6@gmail.com',
    password: '78e715eb052ff77f9b92f3eb59a1360e',
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
RESPONSE
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
      RANDOM: '6e605b0fbbc2cf8a0e1c1de5ade9fc61e67583476f0fb4d261de7b9731e492030c247407bf4a2ae1329c7e85e229cf7b27913acae322e65e1a61126c95029b67a1fbba28f59e68e8145a83b0cbd5f1e99b0f37b905bc7f93ff0d60bec2b83be33913722ebba35cff322fab32a35915fb31c23b6aa4f5fb6c6d783bae8bea1badf603fcc049a1d9c3596d1462f91c9c7558ed68b49d41da9196173eb65d3c404a69ce7566cbaa1eb8af9231dd77cd10874c0ddc2c7bb5de2e0e460c13a1e438413288a401c4d869071209ff11b4310b69eca92559e90588fca9eec46cee35fd7a004ea825b97c19d141bc54795a1b46f1a68925d7d94f07ce05eb6a7e87216c8fd6871ba272254bb897c2abf40076d1ae57df552593f58fb9fe04f67c9ef718e2b89d0b838943dd609a5699cc3e028b678b5b819ac2a0459814a53816ba278be65a17887ec78e653b6ccbaa7a440bfafbdb2a67fa69ca09640006d57105bbe51a5a15b0c66ac516ca48cb34bed87f98e7bb45433da76b73169288a88173760a2a224e054a350d951ffd76a0c2597546cb1c08e2a49016b63a03fce0d14e3c979c626feb390967fc63c8488a0a46141dde8b22b8f2d49834a09a4a90519f9f6d75b4c2c34a0cf354a1c53d185a6ef620dcd0d9bea2c844461f1f23e84e92d652323448bbbaa47a41be737f73073478f443431eeee419cd5f4573d5cfdb34e15f21e6e88868e02e01e269cc94408a393e12671666e0be6511e721b153013a5ae763677dca5a107b57406ea47c5beec7287aff36afef5adcea769a5d29d21334c0d95827571a29a59ae2487d5a2261f83272c7f7ff12d1b88e501ca46fc2c1e58fac1c0aeb26ab030d2cf1144a0fa47597a32773158d7a9394209af7265c58dd593e25d94f61e3bd376881f15b40b16e0905d83f7de5b51d2391f4ffde51a422f43e411d4595f4f1996b0bc574ef8d2c74556b030639ac143e598da34f7a870123829c97cd8805a5861e2dcd3025dcd6b02981ac469b489d17f36cf1bd882a420f2d8c95c71b551e309b95495e341501654e2c6a81d2d6db071b3c0ca3c783af57cfe503d5470074129a7169ca904a2784bfad78c72b03f7a08753c3512d9f1218254f59a491877df9417ff8d0b92eb22c5389cbede155ba9dbe0a1c4a4b7a68011abc8ac7a49ea17298d198e3563fe36a6336c534f749c0cec13fe2b6dc2a6444987bf7006874f3cbac04a249e6fa01e89e4185d446b1414787e5b7ebb08378e207014319d1e4502a11edcc98d1a9395ac07705a9cd47a73d7e905d416cad1fc03da64706644c75cc49bb970c5a5609d574f5ad38e0d2356f7f13911d1d1ae9be3adea0d0f3ba8d0f731fc3e64274cacc8ad2e6358e268f45ef2706ad3487bc3bc057377447ca0ecd6d7da45493435d49efd4bd05a11bab98d7b17b558c1911860f118977defd2705286f04367fe1f2141f57b216dad8e292c41f5a4fabae88143f00ef0aa8f9652c4c72e75f746926eb280abf2bd3e606b910385e9ba5087ff5e932f6ea7f63f7db6526e2550d1fb73bf1ae67910977068b7d59ae953b92913426363ea083290ccfef77aefa5856da4d83',
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
    USER_TOKEN: '012de8bc303c919af748e211fcc8b803f5214dfa50b5edc449545749c843b3b6',
    NOTIFICATION_UNREAD: 3,
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
    ARL: '3172c43f6df7b546bf012225e56998e65611fa42c0e1c06b9eac1370ca8d9b9385958c65ba604acf1c385988a2c0f0237760ea5ca4d9aef246c49725ecd71e6246f34bf736068d1e9487f49abef1819de907f35eb65e529a5d5ebe9b8da977cc',
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
    OFFER_BOX_URL: 'https://www.deezer.com/us/paywalls?app_version=6.1.22.49&ts=1615952269',
    APPCUSTO_CHECKSUM: '20688b10ac10a8cdc38ca066f9373222',
    CHECKSUMS: { OFFERWALL: 'd0abd39bca578263debb2f50dfb27ff6' },
    SHOW_OPTIN_BUTTON: {
      URL: 'https://www.deezer.com/us/optin-box-inapp',
      LABEL: 'Notification preferences'
    },
    HAS_CONSENT: false
  }
}
signInWithEmail END
