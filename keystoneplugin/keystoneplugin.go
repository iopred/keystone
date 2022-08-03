package keystoneplugin

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/iopred/bruxism"
)

type keystoneDungeonID int

const (
	BRH keystoneDungeonID = iota
	COS
	DHT
	EOA
	HOV
	MOS
	NL
	ARC
	VOW
	LKARA
	UKARA
	COEN
	SEAT
	AD
	FH
	KR
	SOTS
	SOB
	TOS
	TM
	UR
	TD
	WM
	OMJ
	OMW
	DOS
	MOTS
	HOA
	SD
	PF
	SOA
	TNW
	TOP
	TSW
	TSG
	GRD
	ID
)

type region int

const (
	US = iota
	EU
)

type keystoneDungeon struct {
	Name    string
	Aliases []string
}

var dungeons map[keystoneDungeonID]*keystoneDungeon = map[keystoneDungeonID]*keystoneDungeon{
	BRH: &keystoneDungeon{
		Name:    "Black Rook Hold",
		Aliases: []string{"black rook hold", "brh", "black rook", "hold"},
	},
	COS: &keystoneDungeon{
		Name:    "Court of Stars",
		Aliases: []string{"court of stars", "cos", "court"},
	},
	DHT: &keystoneDungeon{
		Name:    "Darkheart Thicket",
		Aliases: []string{"darkheart thicket", "dht", "thicket", "darkheart"},
	},
	EOA: &keystoneDungeon{
		Name:    "Eye of Azshara",
		Aliases: []string{"eye of azshara", "eoa", "eye", "azshara"},
	},
	HOV: &keystoneDungeon{
		Name:    "Halls of Valor",
		Aliases: []string{"halls of valor", "hall of valor", "hov"},
	},
	MOS: &keystoneDungeon{
		Name:    "Maw of Souls",
		Aliases: []string{"maw of souls", "mos", "maw"},
	},
	NL: &keystoneDungeon{
		Name:    "Neltharion's Lair",
		Aliases: []string{"neltharion's lair", "nl", "neltharions lair", "nel", "nelth", "lair"},
	},
	ARC: &keystoneDungeon{
		Name:    "The Arcway",
		Aliases: []string{"the arcway", "arc", "arcway"},
	},
	VOW: &keystoneDungeon{
		Name:    "Vault of the Wardens",
		Aliases: []string{"vault of the wardens", "vow", "vault", "warden", "wardens"},
	},
	LKARA: &keystoneDungeon{
		Name:    "Lower Karazhan",
		Aliases: []string{"lower karazhan", "lower kara", "lk", "lkara", "lower"},
	},
	UKARA: &keystoneDungeon{
		Name:    "Upper Karazhan",
		Aliases: []string{"upper karazhan", "upper kara", "uk", "ukara", "upper"},
	},
	COEN: &keystoneDungeon{
		Name:    "Cathedral of Eternal Night",
		Aliases: []string{"cathedral of eternal night", "coen", "cen", "cathedral", "cathedral of night", "cathedral eternal night", "eternal night"},
	},
	SEAT: &keystoneDungeon{
		Name:    "Seat of the Triumvirate",
		Aliases: []string{"seat of the triumvirate", "seat", "sott", "triumvirate", "seat of triumvirate", "seat the triumvirate"},
	},
	AD: &keystoneDungeon {
		Name:    "Atal'Dazar",
		Aliases: []string{"ad", "atal", "atal'dazar", "ataldazar"},
	},
	FH: &keystoneDungeon {
		Name:    "Freehold",
		Aliases: []string{"fh", "freehold"},
	},
	KR: &keystoneDungeon {
		Name:    "King's Rest",
		Aliases: []string{"kr", "kings", "kings rest", "king's rest"},
	},
	SOTS: &keystoneDungeon {
		Name:    "Shrine of the Storm",
		Aliases: []string{"shrine", "sots", "shrine of the storm"},
	},
	SOB: &keystoneDungeon {
		Name:    "Siege of Boralus",
		Aliases: []string{"sob", "siege", "boralus", "siege of boralus"},
	},
	TOS: &keystoneDungeon {
		Name:    "Temple of Sethraliss",
		Aliases: []string{"tos", "temple", "temple of sethraliss", "sethraliss"},
	},
	TM: &keystoneDungeon {
		Name:    "The MOTHERLODE!!",
		Aliases: []string{"tm", "motherlode", "mother", "the motherlode", "ml"},
	},
	UR: &keystoneDungeon {
		Name:    "The Underrot",
		Aliases: []string{"ur", "underrot", "the underrot"},
	},
	TD: &keystoneDungeon {
		Name:    "Tol Dagor",
		Aliases: []string{"td", "tol", "dagor", "tol dagor"},
	},
	WM: &keystoneDungeon {
		Name:    "Waycrest Manor",
		Aliases: []string{"wm", "waycrest", "manor", "waycrest manor"},
	},
	OMJ: &keystoneDungeon {
		Name:    "Operation: Mechagon - Junkyard",
		Aliases: []string{"omj", "operation junkyard", "mechagon junkyard", "lower mecha", "junk", "yard", "junkyard"},
	},
	OMW: &keystoneDungeon {
		Name:    "Operation: Mechagon - Workshop",
		Aliases: []string{"omw", "operation workshop", "mechagon workshop", "upper mecha", "work", "workshop"},
	},
	DOS: &keystoneDungeon {
		Name:    "De Other Side",
		Aliases: []string{"dos", "de other side", "the other side", "other side"},
	},
	MOTS: &keystoneDungeon {
		Name:    "Mists of Tirna Scithe",
		Aliases: []string{"mots", "mts", "mists", "mists of tirna scithe", "tirna", "scithe", "tirna scithe"},
	},
	HOA: &keystoneDungeon {
		Name:    "Halls of Atonement",
		Aliases: []string{"hoa", "halls", "halls of atonement", "atonement"},
	},
	SD: &keystoneDungeon {
		Name:    "Sanguine Depths",
		Aliases: []string{"sd", "sanguine depths"},
	},
	PF: &keystoneDungeon {
		Name:    "Plaguefall",
		Aliases: []string{"pf", "plaguefall", "plague"},
	},
	SOA: &keystoneDungeon {
		Name:    "Spires of Ascension",
		Aliases: []string{"soa", "spires", "ascension", "spires of ascension"},
	},
	TNW: &keystoneDungeon {
		Name:    "The Necrotic Wake",
		Aliases: []string{"tnw", "nw", "necrotic wake", "the necrotic wake"},
	},
	TOP: &keystoneDungeon {
		Name:    "Theatre of Pain",
		Aliases: []string{"top", "theatre", "theater", "theatre of pain", "theater of pain"},
	},
	TSW: &keystoneDungeon {
		Name:    "Tazavesh: Streets of Wonder",
		Aliases: []string{"tsw", "tsow", "Streets of Wonder", "streets"},
	},
	TSG: &keystoneDungeon {
		Name:    "Tazavesh: Solheah's Gambit",
		Aliases: []string{"tsg", "soleah", "Solheah's Gambit", "gambit"},
	},
	GRD: &keystoneDungeon {
		Name:    "Grimrail Depot",
		Aliases: []string{"gd", "grimrail depot", "grimrail", "depot", "train"},
	},
	ID: &keystoneDungeon {
		Name:    "Iron Docks",
		Aliases: []string{"id", "iron docks", "docks", "boat"},
	},
}

type keystone struct {
	User      string
	Alt       string
	Dungeon   keystoneDungeonID
	Level     int
	Depleted  bool
	Modifiers []string
}

func (k *keystone) String() string {
	str := fmt.Sprintf("Level %d **%s**", k.Level, dungeons[k.Dungeon].Name)
	if len(k.Modifiers) != 0 {
		str += " *(" + strings.Join(k.Modifiers, ", ") + ")*"
	}
	if k.Depleted {
		str += " - Depleted"
	}
	return str
}

type keystoneChannel struct {
	Users        map[string]*keystone
	Region       region
	LastModified time.Time
}

var location *time.Location

func init() {
	var err error
	location, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		location = time.Now().Location()
	}
}

func lastTuesday(t time.Time) time.Time {
	year, month, day := t.Date()
	t = time.Date(year, month, day, 0, 0, 0, 0, location)
	for t.Weekday() != time.Tuesday {
		t = t.Add(-24 * time.Hour)
	}

	return t
}

func lastWednesday(t time.Time) time.Time {
	year, month, day := t.Date()
	t = time.Date(year, month, day, 0, 0, 0, 0, location)
	for t.Weekday() != time.Wednesday {
		t = t.Add(-24 * time.Hour)
	}

	return t
}

func (c *keystoneChannel) check() {
	var lastReset time.Time
	if c.Region == EU {
		lastReset = lastWednesday(c.LastModified)
	} else {
		lastReset = lastTuesday(c.LastModified)
	}
	if time.Now().After(lastReset.Add(24 * 7 * time.Hour)) {
		c.Users = map[string]*keystone{}
	}
}

func (c *keystoneChannel) add(bot *bruxism.Bot, service bruxism.Service, message bruxism.Message, userID, alt, query string) bool {
	query = strings.ToLower(query)
	for dungeonID, dungeon := range dungeons {
		for _, alias := range dungeon.Aliases {
			if strings.Index(query, alias+" ") == 0 {
				query = query[len(alias)+1:]
				if len(query) == 0 {
					return false
				}

				parts := strings.Split(query, " ")

				level, err := strconv.Atoi(parts[0])
				if err != nil {
					return false
				}

				parts = parts[1:]

				depleted := false
				modifiers := []string{}

				for _, part := range parts {
					if strings.Index(part, "deplete") == 0 {
						depleted = true
					} else {
						modifiers = append(modifiers, part)
					}
				}

				c.Users[userID] = &keystone{
					User:      message.UserName(),
					Alt:       alt,
					Dungeon:   dungeonID,
					Level:     level,
					Depleted:  depleted,
					Modifiers: modifiers,
				}
				c.LastModified = time.Now()

				return true
			}
		}
	}
	return false
}

type keystoneList []*keystone

// Len is part of sort.Interface.
func (s keystoneList) Len() int {
	return len(s)
}

// Swap is part of sort.Interface.
func (s keystoneList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s keystoneList) Less(i, j int) bool {
	if s[i].Level == s[j].Level {
		return dungeons[s[i].Dungeon].Name > dungeons[s[j].Dungeon].Name
	}
	return s[i].Level > s[j].Level
}

func (c *keystoneChannel) list(bot *bruxism.Bot, service bruxism.Service, message bruxism.Message) {
	if len(c.Users) == 0 {
		service.SendMessage(message.Channel(), "No keystones have been set this week.")
		return
	}

	keystones := keystoneList{}
	for _, keystone := range c.Users {
		keystones = append(keystones, keystone)
	}

	sort.Sort(keystones)

	content := ""
	for _, keystone := range keystones {
		if len(content) != 0 {
			content += "\n"
		}
		content += keystone.String()

		if keystone.Alt != "" {
			content += " - " + keystone.Alt + " *(" + keystone.User + ")*"
		} else {
			content += " - " + keystone.User
		}
	}
	service.SendMessage(message.Channel(), content)
}

type keystonePlugin struct {
	sync.RWMutex
	Channels map[string]*keystoneChannel
}

// Name returns the name of the plugin.
func (p *keystonePlugin) Name() string {
	return "Keystone"
}

// Load will load plugin state from a byte array.
func (p *keystonePlugin) Load(bot *bruxism.Bot, service bruxism.Service, data []byte) error {
	if data != nil {
		if err := json.Unmarshal(data, p); err != nil {
			log.Println("Error loading data", err)
		}
	}

	if p.Channels == nil {
		p.Channels = make(map[string]*keystoneChannel)
	}

	return nil
}

// Save will save plugin state to a byte array.
func (p *keystonePlugin) Save() ([]byte, error) {
	return json.Marshal(p)
}

// Help returns a list of help strings that are printed when the user requests them.
func (p *keystonePlugin) Help(bot *bruxism.Bot, service bruxism.Service, message bruxism.Message, detailed bool) []string {
	help := []string{}

	if service.IsBotOwner(message) || service.IsModerator(message) {
		if p.Channels[message.Channel()] == nil {
			help = append(help, bruxism.CommandHelp(service, "start", "", "Starts keystone tracking in this channel.")[0])
		} else {
			help = append(help, bruxism.CommandHelp(service, "stop", "", "Stops keystone tracking in this channel.")[0])
		}
		help = append(help, bruxism.CommandHelp(service, "region", "<US|EU>", "Sets your region (default US)")[0])
	}

	ticks := ""
	if service.Name() == bruxism.DiscordServiceName {
		ticks = "`"
	}

	if p.Channels[message.Channel()] != nil {
		help = append(help, []string{
			bruxism.CommandHelp(service, "alt", "<alt name> <any other command>", fmt.Sprintf("Executes a command for an alt. Eg: %s%salt iopred set soa 2%s", ticks, service.CommandPrefix(), ticks))[0],
			bruxism.CommandHelp(service, "set", "<dungeon> <level> [modifiers]", fmt.Sprintf("Sets a keystone. Eg: %s%sset hoa 5 teeming%s", ticks, service.CommandPrefix(), ticks))[0],
			bruxism.CommandHelp(service, "list", "", "Lists all this weeks keystones.")[0],
			bruxism.CommandHelp(service, "deplete", "", "Depletes your keystone")[0],
			bruxism.CommandHelp(service, "undeplete", "", "Undepletes your keystone")[0],
			bruxism.CommandHelp(service, "unset", "", "Unsets your keystone")[0],
		}...)
	}

	if detailed {
		help = append(help, []string{
			"Examples:",
			fmt.Sprintf("%s%sset hoa 5 teeming%s - Adds a Level 5 Halls of Atonement keystone with teeming.", ticks, service.CommandPrefix(), ticks),
			fmt.Sprintf("%s%sset de other side 2 depleted%s - Adds a depleted Level 2 De Other Side keystone.", ticks, service.CommandPrefix(), ticks),
			fmt.Sprintf("%s%sregion EU%s - Sets the region to EU.", ticks, service.CommandPrefix(), ticks),
		}...)
	}

	if len(help) == 0 {
		return nil
	}
	return help
}

// Message handler.
func (p *keystonePlugin) Message(bot *bruxism.Bot, service bruxism.Service, message bruxism.Message) {
	defer bruxism.MessageRecover()
	if !service.IsMe(message) {
		messageChannel := message.Channel()

		if bruxism.MatchesCommand(service, "start", message) || bruxism.MatchesCommand(service, "stop", message) {
			if !service.IsBotOwner(message) && !service.IsModerator(message) {
				service.SendMessage(messageChannel, "You must be a server admin to start tracking mythic keystones.")
				return
			}

			p.Lock()
			defer p.Unlock()

			if bruxism.MatchesCommand(service, "start", message) {
				p.Channels[messageChannel] = &keystoneChannel{
					Users: map[string]*keystone{},
				}
				service.SendMessage(messageChannel, "This channel is now tracking mythic keystones.")
			} else {
				delete(p.Channels, messageChannel)
				service.SendMessage(messageChannel, "This channel is no longer tracking mythic keystones.")
			}
		} else if channel, ok := p.Channels[messageChannel]; ok {
			channel.check()

			alt := ""

			messageMessage := strings.TrimSpace(message.Message())

			lowerMessage := strings.ToLower(messageMessage)
			lowerPrefix := strings.ToLower(service.CommandPrefix())

			if !strings.HasPrefix(lowerMessage, lowerPrefix) {
				return
			}

			messageMessage = messageMessage[len(lowerPrefix):]
			parts := strings.Fields(messageMessage)

			if len(parts) == 0 {
				return
			}

			ticks := ""
			if service.Name() == bruxism.DiscordServiceName {
				ticks = "`"
			}

			command := strings.ToLower(parts[0])

			if command == "region" {
				if !service.IsBotOwner(message) && !service.IsModerator(message) {
					service.SendMessage(messageChannel, "You must be a server admin to change regions.")
					return
				}

				if len(parts) > 1 && strings.ToLower(parts[1]) == "eu" {
					channel.Region = EU
					service.SendMessage(messageChannel, "Your region is now set to EU. Keystones will clear midnight Wednesday.")
				} else {
					channel.Region = US
					service.SendMessage(messageChannel, "Your region is now set to US. Keystones will clear midnight Tuesday.")
				}

				return
			}

			userID := message.UserID()

			if command == "alt" {
				if len(parts) <= 1 {
					service.SendMessage(messageChannel, fmt.Sprintf("Invalid alt command. Eg: %s%salt iopred set mists of tirna scithe 9 depleted%s", ticks, service.CommandPrefix(), ticks))
					return
				} else {
					alt = parts[1]
					userID = strings.ToLower(alt) + "__" + userID
					parts = parts[2:]
					command = strings.ToLower(parts[0])
				}
			}

			keystone := channel.Users[userID]

			if command == "set" {
				if len(parts) > 2 && channel.add(bot, service, message, userID, alt, strings.Join(parts[1:], " ")) {
					service.SendMessage(messageChannel, "Keystone set.")
					channel.list(bot, service, message)
				} else {
					service.SendMessage(messageChannel, fmt.Sprintf("Invalid keystone. Eg: %s%sset mists of tirna scithe 3 sanguine%s", ticks, service.CommandPrefix(), ticks))
				}
			} else if command == "unset" {
				if keystone == nil {
					service.SendMessage(messageChannel, "You haven't set a keystone this week.")
				} else {
					delete(channel.Users, userID)
					service.SendMessage(messageChannel, "Keystone unset.")
					channel.list(bot, service, message)
				}
			} else if command == "list" {
				channel.list(bot, service, message)
			} else if command == "deplete" {
				if keystone == nil {
					service.SendMessage(messageChannel, "You haven't set a keystone this week.")
				} else {
					keystone.Depleted = true
					keystone.User = message.UserName()
					service.SendMessage(messageChannel, "Keystone depleted.")
					channel.list(bot, service, message)
				}
			} else if command == "undeplete" {
				if keystone == nil {
					service.SendMessage(messageChannel, "You haven't set a keystone this week.")
				} else {
					keystone.Depleted = false
					keystone.User = message.UserName()
					service.SendMessage(messageChannel, "Keystone undepleted.")
					channel.list(bot, service, message)
				}
			}
		}
	}
}

// Stats will return the stats for a plugin.
func (p *keystonePlugin) Stats(bot *bruxism.Bot, service bruxism.Service, message bruxism.Message) []string {
	return nil
}

// New will create a new wormhole plugin.
func New() bruxism.Plugin {
	return &keystonePlugin{}
}
