package accountsParser

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
	"main/pkg/global"
	"main/pkg/types"
	"main/pkg/util"
	"strings"
)

func getTestnetRewards(accountData types.AccountData) float64 {
	var err error

	for {
		client := util.GetClient()

		req := fasthttp.AcquireRequest()

		req.SetRequestURI(fmt.Sprintf("https://airdrop-api.initia.xyz/info/initia/%s",
			strings.ToLower(accountData.AccountAddress.String())))
		req.Header.Set("accept", "application/json, text/plain, */*")
		req.Header.Set("accept-language", "ru,en;q=0.9,vi;q=0.8,es;q=0.7,cy;q=0.6")
		req.Header.SetMethod("GET")
		req.Header.SetUserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 YaBrowser/24.10.0.0 Safari/537.36")

		resp := fasthttp.AcquireResponse()

		if err = client.Do(req, resp); err != nil {
			log.Printf("%s | Error When Doing Request When Getting Testnet Rewards: %s",
				accountData.AccountLogData, err)

			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			continue
		}

		message := gjson.Get(string(resp.Body()), "message")

		if message.Exists() && strings.Contains(message.String(), "Expected one matching account, but found none or multiple for") {
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)

			return 0
		}

		dropAmount := gjson.Get(string(resp.Body()), "amount")
		if !dropAmount.Exists() {
			log.Printf("%s | Wrong Response When Getting Testnet Rewards: %s",
				accountData.AccountLogData, string(resp.Body()))

			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			continue
		}

		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)

		return dropAmount.Float() / 1e6
	}
}

func getOnchainRewards(accountData types.AccountData) float64 {
	var err error

	for {
		client := util.GetClient()

		req := fasthttp.AcquireRequest()

		req.SetRequestURI(fmt.Sprintf("https://airdrop-api.initia.xyz/info/onchain/%s",
			strings.ToLower(accountData.AccountAddress.String())))
		req.Header.Set("accept", "application/json, text/plain, */*")
		req.Header.Set("accept-language", "ru,en;q=0.9,vi;q=0.8,es;q=0.7,cy;q=0.6")
		req.Header.SetMethod("GET")
		req.Header.SetUserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 YaBrowser/24.10.0.0 Safari/537.36")

		resp := fasthttp.AcquireResponse()

		if err = client.Do(req, resp); err != nil {
			log.Printf("%s | Error When Doing Request When Getting Onchain Rewards: %s",
				accountData.AccountLogData, err)

			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			continue
		}

		message := gjson.Get(string(resp.Body()), "message")

		if message.Exists() && strings.Contains(message.String(), "No onchain airdrop info found for") {
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)

			return 0
		}

		dropAmount := gjson.Get(string(resp.Body()), "amount")
		if !dropAmount.Exists() {
			log.Printf("%s | Wrong Response When Getting Onchain Rewards: %s",
				accountData.AccountLogData, string(resp.Body()))

			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			continue
		}

		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)

		return dropAmount.Float() / 1e6
	}
}

func ParseAccount(accountData types.AccountData) {
	testnetReward := getTestnetRewards(accountData)
	onchainReward := getOnchainRewards(accountData)

	global.CurrentProgress += 1

	log.Printf("[%d/%d] | %s | %g $INITIA (testnet) | %g $INITIA (onchain)",
		global.CurrentProgress, global.TargetProgress, accountData.AccountLogData, testnetReward, onchainReward)

	if testnetReward > 0 || onchainReward > 0 {
		util.AppendFile("with_balances.txt",
			fmt.Sprintf("%s | %g $INITIA (testnet) | %g $INITIA (onchain)\n",
				accountData.AccountLogData, testnetReward, onchainReward))
	}
}
