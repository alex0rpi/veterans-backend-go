package constants

type MediaContext string

const (
	MediaContextSeasonSlider MediaContext = "season_slider"
	MediaContextHomeSlider   MediaContext = "home_slider"
)

type MediaVariant string

const (
	VariantOriginal MediaVariant = "original"
	VariantBlur     MediaVariant = "blur"
	VariantSmall    MediaVariant = "small"
	VariantMedium   MediaVariant = "medium"
	VariantLarge    MediaVariant = "large"
)

type MediaCategory string

const (
	ProCategory   MediaCategory = "pro"
	BasesCategory MediaCategory = "bases"
	VetsCategory  MediaCategory = "vets"
	BoardCategory MediaCategory = "board"
	OtherCategory MediaCategory = "other"
)
