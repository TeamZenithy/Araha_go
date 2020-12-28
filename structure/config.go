package structure

//Config for discord bot.
type Config struct {
	Prefix                string
	Token                 string
	Owners                []string
	ShardLogChannel       string
	ShardStatusLogChannel string
	LavalinkHost          string
	LavalinkPort          string
	LavalinkPass          string
	RedisHost             string
	RedisPort             string
	RedisPass             string
}
