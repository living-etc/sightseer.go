package linux

type Service struct {
	Description    string
	LoadState      string
	UnitFileState  string
	UnitFilePreset string
	ActiveState    string
}
