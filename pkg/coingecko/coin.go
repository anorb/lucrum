package coingecko

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type CoinResponse struct {
	ID                           string              `json:"id"`
	Symbol                       string              `json:"symbol"`
	Name                         string              `json:"name"`
	AssetPlatformID              string              `json:"asset_platform_id"`
	BlockTimeInMinutes           int                 `json:"block_time_in_minutes"`
	HashingAlgorithm             string              `json:"hashing_algorithm"`
	Categories                   []string            `json:"categories"`
	PublicNotice                 string              `json:"public_notice"`
	AdditionalNotices            []string            `json:"additional_notices"`
	Localization                 Localization        `json:"localization"`
	Description                  Description         `json:"description"`
	Links                        Links               `json:"links"`
	Image                        Image               `json:"image"`
	CountryOrigin                string              `json:"country_origin"`
	GenesisDate                  string              `json:"genesis_date"`
	SentimentVotesUpPercentage   float64             `json:"sentiment_votes_up_percentage"`
	SentimentVotesDownPercentage float64             `json:"sentiment_votes_down_percentage"`
	MarketCapRank                int                 `json:"market_cap_rank"`
	CoingeckoRank                int                 `json:"coingecko_rank"`
	CoingeckoScore               float64             `json:"coingecko_score"`
	DeveloperScore               float64             `json:"developer_score"`
	CommunityScore               float64             `json:"community_score"`
	LiquidityScore               float64             `json:"liquidity_score"`
	PublicInterestScore          float64             `json:"public_interest_score"`
	MarketData                   MarketData          `json:"market_data"`
	PublicInterestStats          PublicInterestStats `json:"public_interest_stats"`
	StatusUpdates                []string            `json:"status_updates"`
	LastUpdated                  time.Time           `json:"last_updated"`
}

type Localization map[string]string

type Description map[string]string

type Links struct {
	Homepage                    []string `json:"homepage"`
	BlockchainSite              []string `json:"blockchain_site"`
	OfficialForumURL            []string `json:"official_forum_url"`
	ChatURL                     []string `json:"chat_url"`
	AnnouncementURL             []string `json:"announcement_url"`
	TwitterScreenName           string   `json:"twitter_screen_name"`
	FacebookUsername            string   `json:"facebook_username"`
	BitcointalkThreadIdentifier int      `json:"bitcointalk_thread_identifier"`
	TelegramChannelIdentifier   string   `json:"telegram_channel_identifier"`
	SubredditURL                string   `json:"subreddit_url"`
	ReposURL                    ReposURL `json:"repos_url"`
}

type Image struct {
	Thumb string `json:"thumb"`
	Small string `json:"small"`
	Large string `json:"large"`
}

type MarketData struct {
	CurrentPrice                           Currencies           `json:"current_price"`
	Ath                                    Currencies           `json:"ath"`
	AthChangePercentage                    Currencies           `json:"ath_change_percentage"`
	AthDate                                map[string]time.Time `json:"ath_date"`
	Atl                                    Currencies           `json:"atl"`
	AtlChangePercentage                    Currencies           `json:"atl_change_percentage"`
	AtlDate                                map[string]time.Time `json:"atl_date"`
	ROI                                    ROI                  `json:"roi"`
	MarketCap                              Currencies           `json:"market_cap"`
	MarketCapRank                          int                  `json:"market_cap_rank"`
	TotalVolume                            Currencies           `json:"total_volume"`
	High24                                 Currencies           `json:"high_24h"`
	Low24                                  Currencies           `json:"low_24h"`
	PriceChange24H                         float64              `json:"price_change_24h"`
	PriceChangePercentage24H               float64              `json:"price_change_percentage_24h"`
	PriceChangePercentage7D                float64              `json:"price_change_percentage_7d"`
	PriceChangePercentage14D               float64              `json:"price_change_percentage_14d"`
	PriceChangePercentage30D               float64              `json:"price_change_percentage_30d"`
	PriceChangePercentage60D               float64              `json:"price_change_percentage_60d"`
	PriceChangePercentage200D              float64              `json:"price_change_percentage_200d"`
	PriceChangePercentage1Y                float64              `json:"price_change_percentage_1y"`
	MarketCapChange24H                     float64              `json:"market_cap_change_24h"`
	MarketCapChangePercentage24H           float64              `json:"market_cap_change_percentage_24h"`
	PriceChange24hInCurrency               Currencies           `json:"price_change_24h_in_currency"`
	PriceChangePercentage1hInCurrency      Currencies           `json:"price_change_percentage_1h_in_currency"`
	PriceChangePercentage24hInCurrency     Currencies           `json:"price_change_percentage_24h_in_currency"`
	PriceChangePercentage7dInCurrency      Currencies           `json:"price_change_percentage_7d_in_currency"`
	PriceChangePercentage14dInCurrency     Currencies           `json:"price_change_percentage_14d_in_currency"`
	PriceChangePercentage30dInCurrency     Currencies           `json:"price_change_percentage_30d_in_currency"`
	PriceChangePercentage60dInCurrency     Currencies           `json:"price_change_percentage_60d_in_currency"`
	PriceChangePercentage200dInCurrency    Currencies           `json:"price_change_percentage_200d_in_currency"`
	PriceChangePercentage1yInCurrency      Currencies           `json:"price_change_percentage_1y_in_currency"`
	MarketCapChange24hInCurrency           Currencies           `json:"market_cap_change_24h_in_currency"`
	MarketCapChangePercentage24hInCurrency Currencies           `json:"market_cap_change_percentage_24h_in_currency"`
	FullyDilutedValuation                  Currencies           `json:"fully_diluted_valuation"`
	TotalSupply                            float64              `json:"total_supply"`
	MaxSupply                              float64              `json:"max_supply"`
	CirculatingSupply                      float64              `json:"circulating_supply"`
	LastUpdated                            time.Time            `json:"last_updated"`
}

type Currencies map[string]float64

type ROI struct {
	Times      float64 `json:"times"`
	Currency   string  `json:"currency"`
	Percentage float64 `json:"percentage"`
}

type ReposURL struct {
	Github    []string `json:"github"`
	Bitbucket []string `json:"bitbucket"`
}

type PublicInterestStats struct {
	AlexaRank   int `json:"alexa_rank"`
	BingMatches int `json:"bing_matches"`
}

func FetchCoin(coin string) (CoinResponse, error) {
	c := CoinResponse{}

	body, err := makeCall(fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s?tickers=false&market_data=true&community_data=false&developer_data=false&sparkline=false", coin))
	if err != nil {
		return c, err
	}

	if err = json.Unmarshal(body, &c); err != nil {
		return c, errors.New("Failed to unmarshal: " + err.Error())
	}

	return c, nil
}
