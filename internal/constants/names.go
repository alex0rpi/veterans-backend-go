package constants

type MediaContext string

const (
	MediaContextSeason MediaContext = "season"
	MediaContextHome   MediaContext = "home"
	MediaContextSingle MediaContext = "single"
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
	ProCategory           MediaCategory = "pro"
	BasesCategory         MediaCategory = "bases"
	VetsCategory          MediaCategory = "vets"
	BoardCategory         MediaCategory = "board"
	HistoricBoardCategory MediaCategory = "historic_board"
)
