package keystorage

// KeyStorage is responsible for managing daily tracking key data
type KeyStorage interface {
	// AddKeyRecords adds key records for a given authorisation code
	AddKeyRecords(string, []KeyRecord) error

	// PurgeRecords removes all records for a given authorisation code
	PurgeRecords(string) error

	// ListRecords streams all known key records to a given channel
	ListRecords(records chan RawKeyRecord, errors chan interface{})
}
