package forex

import (
	"errors"
	"fmt"

	"math"
	"sort"
	"strings"
	"time"

	"github.com/botlabs-gg/sgpdb/v2/bot/paginatedmessages"
	"github.com/botlabs-gg/sgpdb/v2/commands"
	"github.com/botlabs-gg/sgpdb/v2/common"
	"github.com/botlabs-gg/sgpdb/v2/lib/dcmd"
	"github.com/botlabs-gg/sgpdb/v2/lib/discordgo"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const currencyPerPage = 16

var Command = &commands.YAGCommand{
	CmdCategory:         commands.CategoryFun,
	Cooldown:            5,
	Name:                "Forex",
	Aliases:             []string{"Money"},
	Description:         "💱 convert value from one currency to another.",
	RunInDM:             true,
	DefaultEnabled:      true,
	SlashCommandEnabled: true,
	RequiredArgs:        3,
	Arguments: []*dcmd.ArgDef{
		{Name: "Amount", Type: dcmd.Float}, {Name: "From", Type: dcmd.String}, {Name: "To", Type: dcmd.String},
	},

	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		amount := data.Args[0].Float64()
		from := strings.ToUpper(data.Args[1].Str())
		to := strings.ToUpper(data.Args[2].Str())

		currenciesResult, exchangeRateResult, err := common.ForexConvert(amount, from, to)

		// Checks the max amount of pages by the number of symbols on each page
		maxPages := int(math.Ceil(float64(len(*currenciesResult)) / float64(currencyPerPage)))

		if err != nil {
			if errors.Is(err, common.ErrUnknownForexCurrencyError) {
				// If the currency isn't supported by API.
				return paginatedmessages.NewPaginatedResponse(
					data.GuildData.GS.ID, data.ChannelID, 1, maxPages, func(p *paginatedmessages.PaginatedMessage, page int) (*discordgo.MessageEmbed, error) {
						embed, err := errEmbed(*currenciesResult, page)
						if err != nil {
							return nil, err
						}
						return embed, nil
					}), nil
			} else if errors.Is(err, common.ErrFailedConversion) {
				return nil, commands.NewPublicError("Failed to convert, Please verify your input")
			}
			return nil, err
		}

		p := message.NewPrinter(language.English)
		embed := &discordgo.MessageEmbed{
			Title:       "💱Currency Exchange Rate",
			Description: p.Sprintf("\n%.2f **%s** (%s) is %.3f **%s** (%s).", amount, (*currenciesResult)[from], from, exchangeRateResult.Rates[to], (*currenciesResult)[to], to),
			Color:       0xAE27FF,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
		}
		return embed, nil
	},
}

func errEmbed(currenciesResult common.Currencies, page int) (*discordgo.MessageEmbed, error) {
	desc := "CODE | Description\n------------------"
	codes := make([]string, 0, len(currenciesResult))
	for k := range currenciesResult {
		codes = append(codes, k)
	}
	sort.Strings(codes)
	start := (page * currencyPerPage) - currencyPerPage
	end := page * currencyPerPage
	for i, c := range codes {
		if i < end && i >= start {
			desc = fmt.Sprintf("%s\n%s  | %s", desc, c, currenciesResult[c])
		}
	}
	embed := &discordgo.MessageEmbed{
		Title:       "Invalid currency code",
		URL:         common.CurrenciesAPIURL,
		Color:       0xAE27FF,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Description: fmt.Sprintf("Check out available codes on: %s ```\n%s```", common.CurrenciesAPIURL, desc),
	}
	return embed, nil
}
