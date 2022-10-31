package hclient

import "github.com/HellasAIO/HellasAIO/internal/utils"

func (r *Request) SetDefaultHeadersAF() *Request {
	r.SetHeader("User-Agent", utils.UserAgent)
	r.SetHeader("Connection", "keep-alive")
	r.SetHeader("sec-ch-ua", `"Chromium";v="94", "Google Chrome";v="94", ";Not A Brand";v="99"`)
	r.SetHeader("Accept", "application/json, text/javascript, */*; q=0.01")
	r.SetHeader("Content-Type", "application/json; charset=UTF-8")
	r.SetHeader("X-Requested-With", "XMLHttpRequest")
	r.SetHeader("sec-ch-ua-mobile", "?0")
	r.SetHeader("sec-ch-ua-platform", `"Windows"`)
	r.SetHeader("Origin", "https://www.theathletesfoot.gr")
	r.SetHeader("Sec-Fetch-Site", "same-origin")
	r.SetHeader("Sec-Fetch-Mode", "cors")
	r.SetHeader("Sec-Fetch-Dest", "empty")
	r.SetHeader("Accept-Language", "en-US,en;q=0.9")

	return r
}

func (r *Request) SetDefaultHeadersFuel() *Request {
	r.SetHeader("Sec-Ch-Ua", `"-Not.A/Brand";v="8", "Chromium";v="102"`)
	r.SetHeader("Sec-Ch-Ua-Mobile", `?0`)
	r.SetHeader("Sec-Ch-Ua-Platform", `"Windows"`)
	r.SetHeader("Upgrade-Insecure-Requests", `1`)
	r.SetHeader("User-Agent", utils.UserAgent)
	r.SetHeader("Accept", `*/*`)
	r.SetHeader("Sec-Fetch-Site", `same-origin`)
	r.SetHeader("Sec-Fetch-Mode", `navigate`)
	r.SetHeader("Sec-Fetch-Dest", `document`)
	r.SetHeader("Accept-Encoding", `gzip, deflate, br`)
	r.SetHeader("Connection", `keep-alive`)
	r.SetHeader("Accept-Language", `en-US,en;q=0.9`)
	r.SetHeader("Origin", `https://www.fuel.com.gr`)
	r.SetHeader("Cache-Control", `max-age=0`)

	return r
}

func (r *Request) SetDefaultHeadersBuzz() *Request {
	r.SetHeader("Accept", "application/json, text/javascript, */*; q=0.01")
	r.SetHeader("Accept-Language", "en-US,en;q=0.9")
	r.SetHeader("Connection", "keep-alive")
	r.SetHeader("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	r.SetHeader("Origin", "https://www.buzzsneakers.gr")
	r.SetHeader("Sec-Fetch-Dest", "empty")
	r.SetHeader("Sec-Fetch-Mode", "cors")
	r.SetHeader("Sec-Fetch-Site", "same-origin")
	r.SetHeader("User-Agent", utils.UserAgent)
	r.SetHeader("X-Requested-With", "XMLHttpRequest")
	r.SetHeader("sec-ch-ua", "\".Not/A)Brand\";v=\"99\", \"Google Chrome\";v=\"103\", \"Chromium\";v=\"103\"")
	r.SetHeader("sec-ch-ua-mobile", "?0")
	r.SetHeader("sec-ch-ua-platform", `"Windows"`)
	r.SetHeader("Pragma", "no-cache")

	return r
}
