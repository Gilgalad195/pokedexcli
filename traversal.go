package main

var WorldMap = map[string][]string{
	"littleroot-town-area": {"hoenn-route-101-area"},
	"hoenn-route-101":      {"littleroot-town-area", "oldale-town-area"},
	"oldale-town-area":     {"hoenn-route-101-area", "hoenn_route-102-area", "hoenn-route-103-area"},
	"hoenn-route-103-area": {"oldale-town-area"},
	"hoenn-route-102-area": {"petalburg-town-area", "oldale-town-area"},
}
