package qdata

// Key mirrors Qt::Key.
type Key int

const (
	// ---- Unicode Basic Latin block (0x00-0x7f) ---------------------------

	Key_Space        Key = 0x20
	Key_Any          Key = Key_Space
	Key_Exclam       Key = 0x21
	Key_QuoteDbl     Key = 0x22
	Key_NumberSign   Key = 0x23
	Key_Dollar       Key = 0x24
	Key_Percent      Key = 0x25
	Key_Ampersand    Key = 0x26
	Key_Apostrophe   Key = 0x27
	Key_ParenLeft    Key = 0x28
	Key_ParenRight   Key = 0x29
	Key_Asterisk     Key = 0x2a
	Key_Plus         Key = 0x2b
	Key_Comma        Key = 0x2c
	Key_Minus        Key = 0x2d
	Key_Period       Key = 0x2e
	Key_Slash        Key = 0x2f
	Key_0            Key = 0x30
	Key_1            Key = 0x31
	Key_2            Key = 0x32
	Key_3            Key = 0x33
	Key_4            Key = 0x34
	Key_5            Key = 0x35
	Key_6            Key = 0x36
	Key_7            Key = 0x37
	Key_8            Key = 0x38
	Key_9            Key = 0x39
	Key_Colon        Key = 0x3a
	Key_Semicolon    Key = 0x3b
	Key_Less         Key = 0x3c
	Key_Equal        Key = 0x3d
	Key_Greater      Key = 0x3e
	Key_Question     Key = 0x3f
	Key_At           Key = 0x40
	Key_A            Key = 0x41
	Key_B            Key = 0x42
	Key_C            Key = 0x43
	Key_D            Key = 0x44
	Key_E            Key = 0x45
	Key_F            Key = 0x46
	Key_G            Key = 0x47
	Key_H            Key = 0x48
	Key_I            Key = 0x49
	Key_J            Key = 0x4a
	Key_K            Key = 0x4b
	Key_L            Key = 0x4c
	Key_M            Key = 0x4d
	Key_N            Key = 0x4e
	Key_O            Key = 0x4f
	Key_P            Key = 0x50
	Key_Q            Key = 0x51
	Key_R            Key = 0x52
	Key_S            Key = 0x53
	Key_T            Key = 0x54
	Key_U            Key = 0x55
	Key_V            Key = 0x56
	Key_W            Key = 0x57
	Key_X            Key = 0x58
	Key_Y            Key = 0x59
	Key_Z            Key = 0x5a
	Key_BracketLeft  Key = 0x5b
	Key_Backslash    Key = 0x5c
	Key_BracketRight Key = 0x5d
	Key_AsciiCircum  Key = 0x5e
	Key_Underscore   Key = 0x5f
	Key_QuoteLeft    Key = 0x60
	Key_BraceLeft    Key = 0x7b
	Key_Bar          Key = 0x7c
	Key_BraceRight   Key = 0x7d
	Key_AsciiTilde   Key = 0x7e

	// ---- Unicode Latin-1 Supplement block (0x80-0xff) --------------------

	Key_nobreakspace        Key = 0x0a0
	Key_exclamdown          Key = 0x0a1
	Key_cent                Key = 0x0a2
	Key_sterling            Key = 0x0a3
	Key_currency            Key = 0x0a4
	Key_yen                 Key = 0x0a5
	Key_brokenbar           Key = 0x0a6
	Key_section             Key = 0x0a7
	Key_diaeresis           Key = 0x0a8
	Key_copyright           Key = 0x0a9
	Key_ordfeminine         Key = 0x0aa
	Key_guillemotleft       Key = 0x0ab
	Key_notsign             Key = 0x0ac
	Key_hyphen              Key = 0x0ad
	Key_registered          Key = 0x0ae
	Key_macron              Key = 0x0af
	Key_degree              Key = 0x0b0
	Key_plusminus           Key = 0x0b1
	Key_twosuperior         Key = 0x0b2
	Key_threesuperior       Key = 0x0b3
	Key_acute               Key = 0x0b4
	Key_micro               Key = 0x0b5
	Key_paragraph           Key = 0x0b6
	Key_periodcentered      Key = 0x0b7
	Key_cedilla             Key = 0x0b8
	Key_onesuperior         Key = 0x0b9
	Key_masculine           Key = 0x0ba
	Key_guillemotright      Key = 0x0bb
	Key_onequarter          Key = 0x0bc
	Key_onehalf             Key = 0x0bd
	Key_threequarters       Key = 0x0be
	Key_questiondown        Key = 0x0bf
	Key_Agrave              Key = 0x0c0
	Key_Aacute              Key = 0x0c1
	Key_Acircumflex         Key = 0x0c2
	Key_Atilde              Key = 0x0c3
	Key_Adiaeresis          Key = 0x0c4
	Key_Aring               Key = 0x0c5
	Key_AE                  Key = 0x0c6
	Key_Ccedilla            Key = 0x0c7
	Key_Egrave              Key = 0x0c8
	Key_Eacute              Key = 0x0c9
	Key_Ecircumflex         Key = 0x0ca
	Key_Ediaeresis          Key = 0x0cb
	Key_Igrave              Key = 0x0cc
	Key_Iacute              Key = 0x0cd
	Key_Icircumflex         Key = 0x0ce
	Key_Idiaeresis          Key = 0x0cf
	Key_ETH                 Key = 0x0d0
	Key_Ntilde              Key = 0x0d1
	Key_Ograve              Key = 0x0d2
	Key_Oacute              Key = 0x0d3
	Key_Ocircumflex         Key = 0x0d4
	Key_Otilde              Key = 0x0d5
	Key_Odiaeresis          Key = 0x0d6
	Key_multiply            Key = 0x0d7
	Key_Ooblique            Key = 0x0d8
	Key_Ugrave              Key = 0x0d9
	Key_Uacute              Key = 0x0da
	Key_Ucircumflex         Key = 0x0db
	Key_Udiaeresis          Key = 0x0dc
	Key_Yacute              Key = 0x0dd
	Key_THORN               Key = 0x0de
	Key_ssharp              Key = 0x0df
	Key_division            Key = 0x0f7
	Key_ydiaeresis          Key = 0x0ff

	// ---- Custom keys — base 0x01000000 -----------------------------------

	Key_Escape      Key = 0x01000000
	Key_Tab         Key = 0x01000001
	Key_Backtab     Key = 0x01000002
	Key_Backspace   Key = 0x01000003
	Key_Return      Key = 0x01000004
	Key_Enter       Key = 0x01000005
	Key_Insert      Key = 0x01000006
	Key_Delete      Key = 0x01000007
	Key_Pause       Key = 0x01000008
	Key_Print       Key = 0x01000009
	Key_SysReq      Key = 0x0100000a
	Key_Clear       Key = 0x0100000b
	Key_Home        Key = 0x01000010
	Key_End         Key = 0x01000011
	Key_Left        Key = 0x01000012
	Key_Up          Key = 0x01000013
	Key_Right       Key = 0x01000014
	Key_Down        Key = 0x01000015
	Key_PageUp      Key = 0x01000016
	Key_PageDown    Key = 0x01000017

	// ---- Modifier keys ---------------------------------------------------

	Key_Shift       Key = 0x01000020
	Key_Control     Key = 0x01000021
	Key_Meta        Key = 0x01000022
	Key_Alt         Key = 0x01000023
	Key_CapsLock    Key = 0x01000024
	Key_NumLock     Key = 0x01000025
	Key_ScrollLock  Key = 0x01000026

	// ---- Function keys F1..F35 -------------------------------------------

	Key_F1  Key = 0x01000030
	Key_F2  Key = 0x01000031
	Key_F3  Key = 0x01000032
	Key_F4  Key = 0x01000033
	Key_F5  Key = 0x01000034
	Key_F6  Key = 0x01000035
	Key_F7  Key = 0x01000036
	Key_F8  Key = 0x01000037
	Key_F9  Key = 0x01000038
	Key_F10 Key = 0x01000039
	Key_F11 Key = 0x0100003a
	Key_F12 Key = 0x0100003b
	Key_F13 Key = 0x0100003c
	Key_F14 Key = 0x0100003d
	Key_F15 Key = 0x0100003e
	Key_F16 Key = 0x0100003f
	Key_F17 Key = 0x01000040
	Key_F18 Key = 0x01000041
	Key_F19 Key = 0x01000042
	Key_F20 Key = 0x01000043
	Key_F21 Key = 0x01000044
	Key_F22 Key = 0x01000045
	Key_F23 Key = 0x01000046
	Key_F24 Key = 0x01000047
	Key_F25 Key = 0x01000048
	Key_F26 Key = 0x01000049
	Key_F27 Key = 0x0100004a
	Key_F28 Key = 0x0100004b
	Key_F29 Key = 0x0100004c
	Key_F30 Key = 0x0100004d
	Key_F31 Key = 0x0100004e
	Key_F32 Key = 0x0100004f
	Key_F33 Key = 0x01000050
	Key_F34 Key = 0x01000051
	Key_F35 Key = 0x01000052

	// ---- Extra keys ------------------------------------------------------

	Key_Super_L      Key = 0x01000053
	Key_Super_R      Key = 0x01000054
	Key_Menu         Key = 0x01000055
	Key_Hyper_L      Key = 0x01000056
	Key_Hyper_R      Key = 0x01000057
	Key_Help         Key = 0x01000058
	Key_Direction_L  Key = 0x01000059
	Key_Direction_R  Key = 0x01000060

	// ---- International & multi-key character composition -----------------

	Key_AltGr             Key = 0x01001103
	Key_Multi_key         Key = 0x01001120
	Key_Codeinput         Key = 0x01001137
	Key_SingleCandidate   Key = 0x0100113c
	Key_MultipleCandidate Key = 0x0100113d
	Key_PreviousCandidate Key = 0x0100113e

	// ---- Misc functions --------------------------------------------------

	Key_Mode_switch Key = 0x0100117e

	// ---- Japanese keyboard support ---------------------------------------

	Key_Kanji             Key = 0x01001121
	Key_Muhenkan          Key = 0x01001122
	Key_Henkan            Key = 0x01001123
	Key_Romaji            Key = 0x01001124
	Key_Hiragana          Key = 0x01001125
	Key_Katakana          Key = 0x01001126
	Key_Hiragana_Katakana Key = 0x01001127
	Key_Zenkaku           Key = 0x01001128
	Key_Hankaku           Key = 0x01001129
	Key_Zenkaku_Hankaku   Key = 0x0100112a
	Key_Touroku           Key = 0x0100112b
	Key_Massyo            Key = 0x0100112c
	Key_Kana_Lock         Key = 0x0100112d
	Key_Kana_Shift        Key = 0x0100112e
	Key_Eisu_Shift        Key = 0x0100112f
	Key_Eisu_toggle       Key = 0x01001130

	// ---- Korean keyboard support -----------------------------------------

	Key_Hangul           Key = 0x01001131
	Key_Hangul_Start     Key = 0x01001132
	Key_Hangul_End       Key = 0x01001133
	Key_Hangul_Hanja     Key = 0x01001134
	Key_Hangul_Jamo      Key = 0x01001135
	Key_Hangul_Romaja    Key = 0x01001136
	Key_Hangul_Jeonja    Key = 0x01001138
	Key_Hangul_Banja     Key = 0x01001139
	Key_Hangul_PreHanja  Key = 0x0100113a
	Key_Hangul_PostHanja Key = 0x0100113b
	Key_Hangul_Special   Key = 0x0100113f

	// ---- Dead keys -------------------------------------------------------

	Key_Dead_Grave              Key = 0x01001250
	Key_Dead_Acute              Key = 0x01001251
	Key_Dead_Circumflex         Key = 0x01001252
	Key_Dead_Tilde              Key = 0x01001253
	Key_Dead_Macron             Key = 0x01001254
	Key_Dead_Breve              Key = 0x01001255
	Key_Dead_Abovedot           Key = 0x01001256
	Key_Dead_Diaeresis          Key = 0x01001257
	Key_Dead_Abovering          Key = 0x01001258
	Key_Dead_Doubleacute        Key = 0x01001259
	Key_Dead_Caron              Key = 0x0100125a
	Key_Dead_Cedilla            Key = 0x0100125b
	Key_Dead_Ogonek             Key = 0x0100125c
	Key_Dead_Iota               Key = 0x0100125d
	Key_Dead_Voiced_Sound       Key = 0x0100125e
	Key_Dead_Semivoiced_Sound   Key = 0x0100125f
	Key_Dead_Belowdot           Key = 0x01001260
	Key_Dead_Hook               Key = 0x01001261
	Key_Dead_Horn               Key = 0x01001262
	Key_Dead_Stroke             Key = 0x01001263
	Key_Dead_Abovecomma         Key = 0x01001264
	Key_Dead_Abovereversedcomma Key = 0x01001265
	Key_Dead_Doublegrave        Key = 0x01001266
	Key_Dead_Belowring          Key = 0x01001267
	Key_Dead_Belowmacron        Key = 0x01001268
	Key_Dead_Belowcircumflex    Key = 0x01001269
	Key_Dead_Belowtilde         Key = 0x0100126a
	Key_Dead_Belowbreve         Key = 0x0100126b
	Key_Dead_Belowdiaeresis     Key = 0x0100126c
	Key_Dead_Invertedbreve      Key = 0x0100126d
	Key_Dead_Belowcomma         Key = 0x0100126e
	Key_Dead_Currency           Key = 0x0100126f
	Key_Dead_a                  Key = 0x01001280
	Key_Dead_A                  Key = 0x01001281
	Key_Dead_e                  Key = 0x01001282
	Key_Dead_E                  Key = 0x01001283
	Key_Dead_i                  Key = 0x01001284
	Key_Dead_I                  Key = 0x01001285
	Key_Dead_o                  Key = 0x01001286
	Key_Dead_O                  Key = 0x01001287
	Key_Dead_u                  Key = 0x01001288
	Key_Dead_U                  Key = 0x01001289
	Key_Dead_Small_Schwa        Key = 0x0100128a
	Key_Dead_Capital_Schwa      Key = 0x0100128b
	Key_Dead_Greek              Key = 0x0100128c
	Key_Dead_Lowline            Key = 0x01001290
	Key_Dead_Aboveverticalline  Key = 0x01001291
	Key_Dead_Belowverticalline  Key = 0x01001292
	Key_Dead_Longsolidusoverlay Key = 0x01001293

	// ---- Multimedia / internet keys --------------------------------------

	Key_Back     Key = 0x01000061
	Key_Forward  Key = 0x01000062
	Key_Stop     Key = 0x01000063
	Key_Refresh  Key = 0x01000064

	Key_VolumeDown Key = 0x01000070
	Key_VolumeMute Key = 0x01000071
	Key_VolumeUp   Key = 0x01000072
	Key_BassBoost  Key = 0x01000073
	Key_BassUp     Key = 0x01000074
	Key_BassDown   Key = 0x01000075
	Key_TrebleUp   Key = 0x01000076
	Key_TrebleDown Key = 0x01000077

	Key_MediaPlay            Key = 0x01000080
	Key_MediaStop            Key = 0x01000081
	Key_MediaPrevious        Key = 0x01000082
	Key_MediaNext            Key = 0x01000083
	Key_MediaRecord          Key = 0x01000084
	Key_MediaPause           Key = 0x01000085
	Key_MediaTogglePlayPause Key = 0x01000086

	Key_HomePage  Key = 0x01000090
	Key_Favorites Key = 0x01000091
	Key_Search    Key = 0x01000092
	Key_Standby   Key = 0x01000093
	Key_OpenUrl   Key = 0x01000094

	Key_LaunchMail  Key = 0x010000a0
	Key_LaunchMedia Key = 0x010000a1

	// ---- Launch0..LaunchF ------------------------------------------------

	Key_Launch0 Key = 0x010000a2
	Key_Launch1 Key = 0x010000a3
	Key_Launch2 Key = 0x010000a4
	Key_Launch3 Key = 0x010000a5
	Key_Launch4 Key = 0x010000a6
	Key_Launch5 Key = 0x010000a7
	Key_Launch6 Key = 0x010000a8
	Key_Launch7 Key = 0x010000a9
	Key_Launch8 Key = 0x010000aa
	Key_Launch9 Key = 0x010000ab
	Key_LaunchA Key = 0x010000ac
	Key_LaunchB Key = 0x010000ad
	Key_LaunchC Key = 0x010000ae
	Key_LaunchD Key = 0x010000af
	Key_LaunchE Key = 0x010000b0
	Key_LaunchF Key = 0x010000b1

	Key_MonBrightnessUp        Key = 0x010000b2
	Key_MonBrightnessDown      Key = 0x010000b3
	Key_KeyboardLightOnOff     Key = 0x010000b4
	Key_KeyboardBrightnessUp   Key = 0x010000b5
	Key_KeyboardBrightnessDown Key = 0x010000b6
	Key_PowerOff               Key = 0x010000b7
	Key_WakeUp                 Key = 0x010000b8
	Key_Eject                  Key = 0x010000b9
	Key_ScreenSaver            Key = 0x010000ba
	Key_WWW                    Key = 0x010000bb
	Key_Memo                   Key = 0x010000bc
	Key_LightBulb              Key = 0x010000bd
	Key_Shop                   Key = 0x010000be
	Key_History                Key = 0x010000bf
	Key_AddFavorite            Key = 0x010000c0
	Key_HotLinks               Key = 0x010000c1
	Key_BrightnessAdjust       Key = 0x010000c2
	Key_Finance                Key = 0x010000c3
	Key_Community              Key = 0x010000c4
	Key_AudioRewind            Key = 0x010000c5
	Key_BackForward            Key = 0x010000c6
	Key_ApplicationLeft        Key = 0x010000c7
	Key_ApplicationRight       Key = 0x010000c8
	Key_Book                   Key = 0x010000c9
	Key_CD                     Key = 0x010000ca
	Key_Calculator             Key = 0x010000cb
	Key_ToDoList               Key = 0x010000cc
	Key_ClearGrab              Key = 0x010000cd
	Key_Close                  Key = 0x010000ce
	Key_Copy                   Key = 0x010000cf
	Key_Cut                    Key = 0x010000d0
	Key_Display                Key = 0x010000d1
	Key_DOS                    Key = 0x010000d2
	Key_Documents              Key = 0x010000d3
	Key_Excel                  Key = 0x010000d4
	Key_Explorer               Key = 0x010000d5
	Key_Game                   Key = 0x010000d6
	Key_Go                     Key = 0x010000d7
	Key_iTouch                 Key = 0x010000d8
	Key_LogOff                 Key = 0x010000d9
	Key_Market                 Key = 0x010000da
	Key_Meeting                Key = 0x010000db
	Key_MenuKB                 Key = 0x010000dc
	Key_MenuPB                 Key = 0x010000dd
	Key_MySites                Key = 0x010000de
	Key_News                   Key = 0x010000df
	Key_OfficeHome             Key = 0x010000e0
	Key_Option                 Key = 0x010000e1
	Key_Paste                  Key = 0x010000e2
	Key_Phone                  Key = 0x010000e3
	Key_Calendar               Key = 0x010000e4
	Key_Reply                  Key = 0x010000e5
	Key_Reload                 Key = 0x010000e6
	Key_RotateWindows          Key = 0x010000e7
	Key_RotationPB             Key = 0x010000e8
	Key_RotationKB             Key = 0x010000e9
	Key_Save                   Key = 0x010000ea
	Key_Send                   Key = 0x010000eb
	Key_Spell                  Key = 0x010000ec
	Key_SplitScreen            Key = 0x010000ed
	Key_Support                Key = 0x010000ee
	Key_TaskPane               Key = 0x010000ef
	Key_Terminal               Key = 0x010000f0
	Key_Tools                  Key = 0x010000f1
	Key_Travel                 Key = 0x010000f2
	Key_Video                  Key = 0x010000f3
	Key_Word                   Key = 0x010000f4
	Key_Xfer                   Key = 0x010000f5
	Key_ZoomIn                 Key = 0x010000f6
	Key_ZoomOut                Key = 0x010000f7
	Key_Away                   Key = 0x010000f8
	Key_Messenger              Key = 0x010000f9
	Key_WebCam                 Key = 0x010000fa
	Key_MailForward            Key = 0x010000fb
	Key_Pictures               Key = 0x010000fc
	Key_Music                  Key = 0x010000fd
	Key_Battery                Key = 0x010000fe
	Key_Bluetooth              Key = 0x010000ff
	Key_WLAN                   Key = 0x01000100
	Key_UWB                    Key = 0x01000101
	Key_AudioForward           Key = 0x01000102
	Key_AudioRepeat            Key = 0x01000103
	Key_AudioRandomPlay        Key = 0x01000104
	Key_Subtitle               Key = 0x01000105
	Key_AudioCycleTrack        Key = 0x01000106
	Key_Time                   Key = 0x01000107
	Key_Hibernate              Key = 0x01000108
	Key_View                   Key = 0x01000109
	Key_TopMenu                Key = 0x0100010a
	Key_PowerDown              Key = 0x0100010b
	Key_Suspend                Key = 0x0100010c
	Key_ContrastAdjust         Key = 0x0100010d

	// ---- LaunchG..LaunchH, touchpad, mic, media-coloured ---------------

	Key_LaunchG       Key = 0x0100010e
	Key_LaunchH       Key = 0x0100010f
	Key_TouchpadToggle Key = 0x01000110
	Key_TouchpadOn    Key = 0x01000111
	Key_TouchpadOff   Key = 0x01000112
	Key_MicMute       Key = 0x01000113
	Key_Red           Key = 0x01000114
	Key_Green         Key = 0x01000115
	Key_Yellow        Key = 0x01000116
	Key_Blue          Key = 0x01000117
	Key_ChannelUp     Key = 0x01000118
	Key_ChannelDown   Key = 0x01000119
	Key_Guide         Key = 0x0100011a
	Key_Info          Key = 0x0100011b
	Key_Settings      Key = 0x0100011c
	Key_MicVolumeUp   Key = 0x0100011d
	Key_MicVolumeDown Key = 0x0100011e
	Key_Keyboard      Key = 0x0100011f
	Key_New           Key = 0x01000120
	Key_Open          Key = 0x01000121
	Key_Find          Key = 0x01000122
	Key_Undo          Key = 0x01000123
	Key_Redo          Key = 0x01000124

	Key_MediaLast Key = 0x0100ffff

	// ---- Keypad navigation ----------------------------------------------

	Key_Select Key = 0x01010000
	Key_Yes    Key = 0x01010001
	Key_No     Key = 0x01010002

	// ---- Newer misc keys -------------------------------------------------

	Key_Cancel  Key = 0x01020001
	Key_Printer Key = 0x01020002
	Key_Execute Key = 0x01020003
	Key_Sleep   Key = 0x01020004
	Key_Play    Key = 0x01020005
	Key_Zoom    Key = 0x01020006
	Key_Exit    Key = 0x0102000a

	// ---- Device keys -----------------------------------------------------

	Key_Context1         Key = 0x01100000
	Key_Context2         Key = 0x01100001
	Key_Context3         Key = 0x01100002
	Key_Context4         Key = 0x01100003
	Key_Call             Key = 0x01100004
	Key_Hangup           Key = 0x01100005
	Key_Flip             Key = 0x01100006
	Key_ToggleCallHangup Key = 0x01100007
	Key_VoiceDial        Key = 0x01100008
	Key_LastNumberRedial Key = 0x01100009

	Key_Camera      Key = 0x01100020
	Key_CameraFocus Key = 0x01100021

	// ---- Unknown ---------------------------------------------------------

	Key_unknown Key = 0x01ffffff
)
