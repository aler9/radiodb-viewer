package shared

import (
    "fmt"
    "strings"
    "rdbviewer/back/defs"
)

func TourLabel(s *defs.RadioShow) (string) {
    switch s.Tour {
        // rh
        case "amsp":     return "A Moon Shaped Pool tour"
        case "tkol":     return "The King of Kimbs tour"
        case "irb":      return "In Rainbows tour"
        case "httt":     return "Hail to the Thief tour"
        case "amn":      return "Amnesiac tour"
        case "kida":     return "Kid A tour"
        case "okc":      return "OK Computer tour"
        case "bends":    return "The Bends tour"
        case "pbh":      return "Pablo Honey tour"
        case "oaf":      return "On a Friday"

        // jonny / junun
        case "varj":     return "Various Jonny"
        case "junun":    return "Junun tour"

        // thom / afp
        case "tmb":     return "Tomorrow's Modern Boxes tour"
        case "vart":    return "Various Thom"
        case "amok":    return "Amok tour"
        case "eras":    return "The Eraser tour"

        // phil
        case "weath":   return "The Weatherhouse tour"
        case "famil":   return "Familial tour"
    }
    panic("tour not recognized")
}

func ArtistLabel(s *defs.RadioShow) string {
    switch s.Artist {
    case "radiohead":   return "Radiohead"
    case "thom":        return "Thom Yorke"
    case "afp":         return "Atoms for Peace"
    case "jonny":       return "Jonny Greenwood"
    case "phil":        return "Phil Selway"
    case "junun":       return "Junun"
    }
    panic("artist not recognized")
}

func SetlistUrlLabel(url string) string {
    switch url {
    case "citizeninsane":   return "Citizen Insane"
    case "greenplastic":    return "Green Plastic Radiohead"
    case "setlistfm":       return "Setlist.fm"
    case "songkick":        return "Songkick"
    case "thomthomthom":    return "Thom thom thom"
    }
    panic("url not recognized")
}

func CountryLabel(s *defs.RadioShow) string {
    return countryMetaData[s.CountryCode].Name
}

func CountryCodeShort(s *defs.RadioShow) string {
    return countryMetaData[s.CountryCode].CodeShortLower
}

func AudioResolution(b *defs.RadioBootleg) string {
    if b.MinfoAudioRate == 0 {
        return "unknown"
    }
    ret := fmt.Sprintf("%.1fkhz ", float64(b.MinfoAudioRate)/1000.0)
    if b.MinfoAudioDepth != 0 {
        ret += fmt.Sprintf("%dbit", b.MinfoAudioDepth)
    } else {
        ret += "lossy"
    }
    return ret
}

func VideoResolution(b *defs.RadioBootleg) string {
    if b.MinfoVideoHeight == 0 {
        return "unknown"
    }
    return fmt.Sprintf("%dp", b.MinfoVideoHeight)
}

func ShortResolution(b *defs.RadioBootleg) string {
    if b.MinfoFormat == "" ||
        (b.Type == "audio" && b.MinfoAudioCodec == "") ||
        (b.Type == "video" && b.MinfoVideoCodec == "") {
        return "unknown"
    }
    switch b.Type {
    case "audio": return AudioResolution(b)
    case "video": return VideoResolution(b)
    }
    return ""
}

func FormatLabel(b *defs.RadioBootleg) string {
    if b.MinfoFormat == "" {
        return "unknown"
    }
    return b.MinfoFormat
}

func VideoCodecLabel(b *defs.RadioBootleg) string {
    if b.MinfoVideoCodec == "" {
        return "unknown"
    }
    return b.MinfoVideoCodec
}

func AudioCodecLabel(b *defs.RadioBootleg) string {
    if b.MinfoAudioCodec == "" {
        return "unknown"
    }
    return b.MinfoAudioCodec
}

func MediaTypeLabel(media string) string {
    return strings.Title(media)
}
