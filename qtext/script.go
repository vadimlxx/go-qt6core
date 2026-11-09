package qtext

// Script mirrors QChar::Script from Qt 6.11 (Unicode 17.0 data).
type Script int

const (
    Script_Unknown Script = iota
    Script_Inherited
    Script_Common

    Script_Latin
    Script_Greek
    Script_Cyrillic
    Script_Armenian
    Script_Hebrew
    Script_Arabic
    Script_Syriac
    Script_Thaana
    Script_Devanagari
    Script_Bengali
    Script_Gurmukhi
    Script_Gujarati
    Script_Oriya
    Script_Tamil
    Script_Telugu
    Script_Kannada
    Script_Malayalam
    Script_Sinhala
    Script_Thai
    Script_Lao
    Script_Tibetan
    Script_Myanmar
    Script_Georgian
    Script_Hangul
    Script_Ethiopic
    Script_Cherokee
    Script_CanadianAboriginal
    Script_Ogham
    Script_Runic
    Script_Khmer
    Script_Mongolian
    Script_Hiragana
    Script_Katakana
    Script_Bopomofo
    Script_Han
    Script_Yi
    Script_OldItalic
    Script_Gothic
    Script_Deseret
    Script_Tagalog
    Script_Hanunoo
    Script_Buhid
    Script_Tagbanwa
    Script_Coptic

    // Unicode 4.0
    Script_Limbu
    Script_TaiLe
    Script_LinearB
    Script_Ugaritic
    Script_Shavian
    Script_Osmanya
    Script_Cypriot
    Script_Braille

    // Unicode 4.1
    Script_Buginese
    Script_NewTaiLue
    Script_Glagolitic
    Script_Tifinagh
    Script_SylotiNagri
    Script_OldPersian
    Script_Kharoshthi

    // Unicode 5.0
    Script_Balinese
    Script_Cuneiform
    Script_Phoenician
    Script_PhagsPa
    Script_Nko

    // Unicode 5.1
    Script_Sundanese
    Script_Lepcha
    Script_OlChiki
    Script_Vai
    Script_Saurashtra
    Script_KayahLi
    Script_Rejang
    Script_Lycian
    Script_Carian
    Script_Lydian
    Script_Cham

    // Unicode 5.2
    Script_TaiTham
    Script_TaiViet
    Script_Avestan
    Script_EgyptianHieroglyphs
    Script_Samaritan
    Script_Lisu
    Script_Bamum
    Script_Javanese
    Script_MeeteiMayek
    Script_ImperialAramaic
    Script_OldSouthArabian
    Script_InscriptionalParthian
    Script_InscriptionalPahlavi
    Script_OldTurkic
    Script_Kaithi

    // Unicode 6.0
    Script_Batak
    Script_Brahmi
    Script_Mandaic

    // Unicode 6.1
    Script_Chakma
    Script_MeroiticCursive
    Script_MeroiticHieroglyphs
    Script_Miao
    Script_Sharada
    Script_SoraSompeng
    Script_Takri

    // Unicode 7.0
    Script_CaucasianAlbanian
    Script_BassaVah
    Script_Duployan
    Script_Elbasan
    Script_Grantha
    Script_PahawhHmong
    Script_Khojki
    Script_LinearA
    Script_Mahajani
    Script_Manichaean
    Script_MendeKikakui
    Script_Modi
    Script_Mro
    Script_OldNorthArabian
    Script_Nabataean
    Script_Palmyrene
    Script_PauCinHau
    Script_OldPermic
    Script_PsalterPahlavi
    Script_Siddham
    Script_Khudawadi
    Script_Tirhuta
    Script_WarangCiti

    // Unicode 8.0
    Script_Ahom
    Script_AnatolianHieroglyphs
    Script_Hatran
    Script_Multani
    Script_OldHungarian
    Script_SignWriting

    // Unicode 9.0
    Script_Adlam
    Script_Bhaiksuki
    Script_Marchen
    Script_Newa
    Script_Osage
    Script_Tangut

    // Unicode 10.0
    Script_MasaramGondi
    Script_Nushu
    Script_Soyombo
    Script_ZanabazarSquare

    // Unicode 12.1
    Script_Dogra
    Script_GunjalaGondi
    Script_HanifiRohingya
    Script_Makasar
    Script_Medefaidrin
    Script_OldSogdian
    Script_Sogdian
    Script_Elymaic
    Script_Nandinagari
    Script_NyiakengPuachueHmong
    Script_Wancho

    // Unicode 13.0
    Script_Chorasmian
    Script_DivesAkuru
    Script_KhitanSmallScript
    Script_Yezidi

    // Unicode 14.0
    Script_CyproMinoan
    Script_OldUyghur
    Script_Tangsa
    Script_Toto
    Script_Vithkuqi

    // Unicode 15.0
    Script_Kawi
    Script_NagMundari

    // Unicode 16.0
    Script_Garay
    Script_GurungKhema
    Script_KiratRai
    Script_OlOnal
    Script_Sunuwar
    Script_Todhri
    Script_TuluTigalari

    // Unicode 17.0
    Script_Sidetic
    Script_TaiYo
    Script_TolongSiki
    Script_BeriaErfe

    ScriptCount
)
