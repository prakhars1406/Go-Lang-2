package Domain

type PacksAndChannels struct {
	Packs    []PacksDetails
	Channels []ChannelsDetails
}
type Packs struct {
	Packs []PacksDetails
}
type PacksDetails struct {
	ID        int
	ChannelId []int
	Name      string
	Price     int
}
type Channels struct {
	Channels []ChannelsDetails
}
type ChannelsDetails struct {
	ID    int
	Name  string
	Price int
}
