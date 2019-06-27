package defs

type RadioOut struct {
	Stats    RadioOutStats
	Shows    map[string]*RadioShow
	Bootlegs map[string]*RadioBootleg
}
