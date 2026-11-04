package entity

// HotspotCategory defines the strictly typed category of an interactive point of interest.
type HotspotCategory string

const (
	CategoryCoffeeMachine HotspotCategory = "COFFEE_MACHINE"
	CategoryFridge        HotspotCategory = "FRIDGE"
	CategoryDispenser     HotspotCategory = "DISPENSER"
	CategoryProjector     HotspotCategory = "PROJECTOR"
	CategoryViewpoint     HotspotCategory = "VIEWPOINT"
	CategoryEasterEgg     HotspotCategory = "EASTER_EGG"
)

// IsValid checks whether the category belongs to the allowed domain set.
func (c HotspotCategory) IsValid() bool {
	switch c {
	case CategoryCoffeeMachine,
		CategoryFridge,
		CategoryDispenser,
		CategoryProjector,
		CategoryViewpoint,
		CategoryEasterEgg:
		return true
	default:
		return false
	}
}

// Hotspot represents an interactive object pin or hidden Easter egg on the map or inside a room.
type Hotspot struct {
	ID       string
	ZoneID   *string
	Category HotspotCategory
	XPercent float64
	YPercent float64
	Title    string
	Content  string
}

// IsEasterEgg returns true if the hotspot is a hidden Xsolla logo.
func (h *Hotspot) IsEasterEgg() bool {
	return h.Category == CategoryEasterEgg
}
