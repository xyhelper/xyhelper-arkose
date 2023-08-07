package api

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"xyhelper-arkose/config"

	"gitee.com/baixudong/gospider/ja3"
	"gitee.com/baixudong/gospider/requests"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	challengeUrl = "https://client-api.arkoselabs.com/fc/gt2/public_key/35536E1E-65B4-4D96-9D97-6ADB7EFF8147"
	headers      = map[string]string{
		"Origin":          "https://client-api.arkoselabs.com",
		"Referer":         "https://client-api.arkoselabs.com/v2/1.5.4/enforcement.cd12da708fe6cbe6e068918c38de2ad9.html",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/116.0",
		"Content-Type":    "application/x-www-form-urlencoded; charset=UTF-8",
		"Accept":          "*/*",
		"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Site":  "same-origin",

		"Connection":     "keep-alive",
		"Sec-Fetch-Dest": "empty",
		"Cookie":         "gfsessionid=13yrbr4tjmbp40cujvuhyibfk8100im0; _dd_s=rum=0&expire=1691165643160; _account=1 ; _cfuvid=jkjdjfldfjkjkljlfdsjklfjk",
	}
)

func GetTokenByPayload(ctx g.Ctx, payload string, userAgent string) (string, error) {
	// g.Log().Info(ctx, "开始获取token", payload)
	// 以&分割转换为数组
	payloadArray := gstr.Split(payload, "&")
	// 移除最后一个元素
	payloadArray = payloadArray[:len(payloadArray)-1]
	// 将 rnd=0.3046791926621015 添加到数组最后

	payloadArray = append(payloadArray, "rnd="+strconv.FormatFloat(rand.Float64(), 'f', -1, 64))
	// 以&连接数组
	payload = strings.Join(payloadArray, "&")
	// g.Log().Info(ctx, "payload", payload)
	client := g.Client()
	client.SetHeaderMap(headers)
	if config.Proxy != "" {
		client.SetProxy(config.Proxy)
	}
	response, err := client.SetHeader("User-Agent", userAgent).Post(ctx, challengeUrl, payload)
	if err != nil {
		log.Panic(err)
	}
	defer response.Close()
	if response.StatusCode != 200 {
		return "", gerror.New("获取token失败" + response.Status)
	}
	// response.RawDump()
	resBodyStr := response.ReadAllString()
	token := gjson.New(resBodyStr).Get("token").String()
	if strings.Contains(token, "sup=1|rid=") {
		return token, nil
	}
	return "", gerror.New("获取token失败:" + resBodyStr)

}

// 使用ja3指纹获取token
func GetTokenByPayloadJa3(ctx g.Ctx, payload string, userAgent string) (string, error) {
	Ja3Spec, err := ja3.CreateSpecWithId(ja3.HelloFirefox_Auto) //根据id 生成指纹
	if err != nil {
		log.Panic(err)
	}
	reqCli, err := requests.NewClient(ctx, requests.ClientOption{
		Ja3Spec: Ja3Spec,
		H2Ja3:   true,
		Proxy:   config.Proxy,
		// Proxy:   "socks5://codespace.sltapp.cn:21000",
	})
	if err != nil {
		log.Panic(err)
	}
	defer reqCli.Close()
	headers["User-Agent"] = userAgent
	response, err := reqCli.Request(ctx, "post", challengeUrl, requests.RequestOption{
		Headers: headers,
		Data:    payload,
	})
	if err != nil {
		log.Panic(err)
	}
	jsonData, _ := response.Json()
	token := jsonData.Get("token").String()
	if strings.Contains(token, "sup=1|rid=") {
		return token, nil
	}
	return "", gerror.New("获取token失败:" + token)
	// return "", gerror.New("获取token失败")
}

func GPT4(r *ghttp.Request) {
	ctx := r.Context()

	Ja3Spec, err := ja3.CreateSpecWithId(ja3.HelloFirefox_Auto) //根据id 生成指纹
	if err != nil {
		log.Panic(err)
	}
	reqCli, err := requests.NewClient(ctx, requests.ClientOption{
		Ja3Spec: Ja3Spec,
		H2Ja3:   true,
		Proxy:   config.Proxy,
		// Proxy:   "socks5://codespace.sltapp.cn:21000",
	})
	if err != nil {
		log.Panic(err)
	}
	defer reqCli.Close()
	response, err := reqCli.Request(ctx, "post", challengeUrl, requests.RequestOption{
		Headers: r.Header,
		Data:    r.GetBodyString(),
	})
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	defer response.Close()
	text := response.Text()
	r.Response.Status = response.StatusCode()
	r.Response.WriteJsonExit(text)
}

func GetPayloadFromCache(ctx g.Ctx) (payload config.Payload, err error) {
	cache := config.Cache.MustGet(ctx, "payload")
	if cache.IsEmpty() {
		return payload, gerror.New("payload is empty")
	}
	err = gconv.Struct(cache, &payload)
	if err != nil {
		return payload, gerror.New("payload format error")
	}
	return payload, nil

}

func RefreshPayloadFromMaster(ctx g.Ctx) (err error) {
	if g.Cfg().MustGetWithEnv(ctx, "MASTER").String() == "" {
		res := g.Client().GetVar(ctx, "https://chatarkose.xyhelper.cn/payload")
		payloadStr := gjson.New(res).Get("payload").String()
		payloadStr = "bda=eyJjdCI6InNWM3ZKaUl5VnZRSm9oWU9ReHhxQnJRZC9INnlqK1VJaUh5cVZTL20vWFRTMGMxbS9rbDFOSEUrejJCc3ZEL1VrMmowT3hxaTJySUdZV1B5SmtXVlQxMmtZTEk0U055bVREQncrQVNmRWxMRUw3M2VmRjJ0ZFB3U2UyTE5CUURNSmpyN1VzTVRZcTNhc0Qvd1hQZllYUXd3YWNMVXpxZFNRZXdHdUo2V1VsUW51eTArckNwVU54TGx3U2UzRHZFRExOc2V5RTVrV2lnWDRhZ0FiM2k0ZkVINFE3VUxBZTVQNDh1a3hocW5FWlNhM3JEMnFmeGdrQmxtNzErQjljQ2E2NFY0MVF4a3VpSTZFbnRHVTVqRUNUdWJvMHRwVEQ0Zkc2eml3UnpDWHZWeHFveThiVFlKa2lFVFJqRTBvWmluUm9NazE0OGh1N3htT3VUZlZzY1VueGIvU29WbXNWUFBmRDlYc01rUXhQTUUzMkVXaXRuNlNtTVZkQ3ZBZnd5YktGcEF5aVI2dUJhT1YwaHAwU3kyeGlnUURtbk00b2Y4bmFiNEg0Vy9JSlpWVFJpQXNRS0I0YjRJMGJEenJJOG1PRVJ3TWhoazNDZTBwcGF1OUo3QS8vQWVWTUlRc3kyZ1lrcWVmV081T0Qwa2pvcWFQT1p2bUl2UU12SVMrQ0VzY0tlT2c4aDV5VUM1bTU0K1A5dWtWUURwU08yMUZyR1Nlei9veG9HTExMS1BtVDFjbktGWmc0NHBFS2ZKc2ZZc1RKQnIrU3BBUEJ0b0wwRVNGY0RuU1ZKYVBaK3pXaDllVW9RTi8vVTRqdU5PWUh5SlpHcHdWYmprcDdpWkF5Y01Ld2FLT1c0dzhaQ2IxWFJXaUtoMjVvNEdYM2VnM2pkeGFnRjNiRmIxSGw5MnFROXM5cWhDZGh5NGhpaEpDVEdSSDRTanQwc2wwMVB1Z1hYOENEenZ6ZlpJYkI2UFB2ZVEwUHArUVh0RUh1OWZ6RVVaZytmelV4YlIzcjBmd1FDVWxSWWU2OE44QU5IQVIxb3daWWltZTI1TGZMRk14ZGdIUzJoSUhSMVJKSmpKVEJ0K0pGNVZCbEhWTW1qc2JWMWQxUERlREpTRFlRaFkxbm8yaGl3MUZ6TmZpWmlHaTZlTFUxUnM2M0xEUmx0ZnhINW50eFR3Rk13SldTWU9tZU9oRDdrNWNlNG9tVzFScElCeDJ4RkovbjdTWXpKVkU3S0NVQVA5N3BTQWZNK2E4akZKSHMrd0tlWFFva2FjUVRlRWQ2UG1Hd1ZzZ1RUeFZEeTBvczZXbGdxQTlSY0FtMGpGWFREMERPRjRzQVRHZHNXZzNhdEJuNDNCTytwaHN5cTFlbDdqOGtXZEtRL0pEaEFNd3V2NENydUxPZHdDUzMzT1lnaDNmb2FKOGNkeVpRVG54Ni8yNVQyOThYQVhBUW1keiswZDF3WHV0U2tyTStHQU80VmxzY0ZRaVlJaUpEWCt2djVnMFNhbzNPWGNtZXJYVGJ5RjZaVC9vVTNFeWFUMVVVdExpR3czdDcrWWdyQmRsLzFOTUNXOTl6U0hON0h2UWE3Yit4YlFMQVdsR0R0N3NiYnR5aDlFSGY5YURacGpuK0VRRkFVL1RtT3JQR0RUU0FSK0lRVnhQdXZuS09KQU9LQ2FIenNjcW5uTmlqZEsrbHR3M0k5QXh6QWlNV2tieU5IcCtnVExjdWNwVmdHbEpRanpQS3c5L0p1dDZBcWtudU44NXo2VEJ5Vkg0ZTFGZ3FrN3ZPeU8xMHIwemdRdVRBTTlyWHk0QTk1ZERCMTFXelhZc21NT0JKUTdlT2RNeit0bUFqQlVybkN5bkdFQjRlSS9rQnpndUs5Q2F0UVE1azJTcFBNZGUreDM4OGpLQk9CMWJqdXdFdVZrSmFaQ0pCOEVFYXhyUUVNa0F2Y2xLTUlsWSt4UHZvVGlHdVE5VUJRQ3hXMVNFOEtRSzYxQVhLTEJOdXJZckdVZzNZTHVzSFJkZ3pUajB1aWt4aFJCMllMeFJkVXBneSthR0p0VlpGdk9zV1I2YzdIUTZBajVPMmYyQWx5b21WaWR0YksxL2ViSVZqWmFRb1FsNzBoNjk1T1R2K1czb3FXY09XaTZraThOdUpid3VDcUlGSnlreXFxWnQzMFRMT2Q3bUo2dG0yR0lRajhTVUZ1VHVDVWt0d0lLUE43MEpXNlZJVDdmSHN3UzE5TVhmbFQ2VzB2aDFwUjlXWXl4b09ZMGl1QVBpajR1cUlkc0tpMHF5a25Pckhnc0VadUk0MWVEYitFdmNOOGxhandnZkpaMXZjbm52aHMxMGxZdmoyVUdkYWZ0V3Q0L3JSV1o5N0VFbkR4MEN5dEV2MzRUMDd4eGtIdENnclJSVXdPMnBJVnF1SlFPM3NRSlI0ckNYZ3dDektIRSs3OER0MHpUZSt0dmJ3ZGNxYTNCdGlDNWRPSGFzWkV5MVA0dnpBODRMMDQ0UFZEbUtybmJ4UUdqNkNFdGI4N29LQUNzcWhLVlNac1F2ODhCdVF6czJmaU9tMk90RFgrTmtqclNTdy8vTnFycTRZOHFlTEFSNktNdjZjaC9aa0JhVlluSHFwWXQ3YnF4TGhpN2ZaQXA3STFFbHRqTlFySG9IYTBYU1RnWUg3Y1pycE9SdzZSV2dkdGowU2xjSWJRaS85VTZEM2JwWmxWeUdWQk0wd25FQXBzV2d6Ty9HRURSdDZubWY3S1lIMnBKWkFvVWFMMEFidk41ckt6dTVObjZldGRFeUNRWjlOSFo5UzhSUmlvS05hOXlaZXVVTHUxR2llR0l6NlU2TjlXS3VPRW1DTCtaL3JpTWRXWUdSbXBLSTNKYWV1OHUvandsUC9OeXpQcGc5c2lYMlpSSlduRmdjSlRITzZxeGhYVlBuVDdZL2dCTE1Rcld5dmtKY1NyNEo4dTBDWVVXVWJtdmd5RDlEV2kya0kwaTZ4UTcyZ0NaTnBVN2VKMGYya2xxc0RYMTZQb0syWktwZXFMK1VFWWhXVTMrQjF5clB2MG9oM0hka1hRR2xUZm11cGVXMmZpSS9RbTk5NWZhbkZJQ1dxUjl0VjVMN0dMdjFzMllhWmlycFBEUjM3S2wrejIxbXdtd3NyM1R5R25GbFJiM3RNRlVSZHpaTlFxcDlrakdQaytJeTBsWFJxVUl2TkVYL0cwRE9jVWNVcUx0d3ZKU21sVXE0amM2elpzYkFUSDg1SnhGYUQxRzNaczBtUGJndlpidzdPK01qM3EyWjhtYnZCTHpFN2pNQytFczFTdFV2aGJDVUg0ZHJwVHBzdmJnOWtwazQ2bjdpR3NzOVFEOGZEWE5VWmtsWnhnUFZDdFRlWnFGcElpK2FRdmk5QzVZK05YamVDL3h2YUZFODNNbzNMSkdDUXVtcUpYU1pVUCt4Smcyb3NMbXB3VmhLc0E0VHVOSkE2YlRoR0Q2TGlnY0M0RTFYZ2lyb1VlcS9peWU1WmhPY1RZRFpDbWRxcElKR0xjSnBEbzQzSGdUOVNaTEJJd2RVa3FKVUtpZTBjOHJ1cWltYzJ1YVRCQzFQeVc1bktXMGhLOWRSQWNzSmdSaTFuZGY5WVozYnliMzBpclRobnZQS0cyMG8wT0VRUElFQnNXSDE1eXlUSFFuVGJETjMyaVRneEJZOStUeHVEd25CTTJkT0tObGJncGN6cVdEc1hTQ3FKYzNEcnR4dHhKM01RTkd5NFZrMVRORjM4eDNMZjZ1eENLc0sxcjF1Y3YwcnJJZ3NKSWI0eCtLRHZxOVNCWWpHelovL0lKakdmYzJlK0duVTZUVWl0UUZLSlZWNXBodkQzVENKQmIrWjFHbUZZRzZ2Z0RPdTkzVE81V0Rzay95YVdUMWpsa1FJUFBGYU52eGxKdlJmaXBQWFgrcWw1NmllRTdkR0N5cWJyb0NhTkgwQnlDYzhzRVZwRnBnMTF6ZEM1UGl3Y0xjSkdRSVRXaDBUbjFkVWE3YkRYNGptV2RpcURNTG41a2ZZV3ZFRy9CNXpMWldQZi9BUEtaSkROL0ZiVk9mdGsrcHcva3VkWktwKzFnSGFpUzFQSXI2TVdMY2dWVnN6TFc2SEZxbTZRSnBXNVJkdGJRcFhXSXRmRXRxRDNUTDdPc1JPSlgrLzlCdXZRZGxwU25LMlB0QzdUc0YwQ3R2RW83YzF4Z1R1R0JLbHk2T2NFTzgyTXA0TENPNzdqZTJEcnY4bzg5Y2RhSGx0dndjRE5kSTQzS09RMkUyUEJzQjYxbVp0QXdYS0NwVExiQ0JRZ3VsV09TbWVpakIyVDh4ZDJBczh2WkE1YU1VSXRudnpNdGVZTHB6ZkpqM2FJc3ZCM2lDYmJqOXNBNmxJRnNsanBjdHZtL1EzSTdtNERUQWcyVjVCbEVaMXR5dGVpeURzaHk1UlEraG1WSzN2dmphMDRXRnFwLzdkQUk1SkFrUGZRNERYYm5sUWVydE54Q1Y5NkpObllDMVRPdm8wRytpd2JhbVZQTUw3NEI1WmhVemN3azlyNTBOWFRvZ2NwdVhQVnFIbTN1eW02K0dOKytkbnBQWnA0Q2FVVEpUaE9OVUxGbjROdEpkYjBQRlo5Y3VNMnZWR25NZGVJUktmNXVyK3hHSUtTRWcvajFkYm5QWU9YYlhpMHNZdi9YeDM0bXQ3Zy9SMFR0OEJwZWdnWkFaeWNaWDY2MXZrTmdSTjFxbm9CUnRmVHFQZmNaUzF1a1pGSmZqdG84cmtqMmF4QzVOczFndWpGbmJnaEZxU0pXOXJDZ21BdUlXTzh0RGZCRTA2RU1OUE9vNW9ZNVAwWHYxM1FJR0RQVHVERHNJWkg3YkNVakE5dFYrNjRpVVVGeDY3U2kxakJpMnN3VENRQytEVHhpNVNnVmJWQjByNEhaU09oZVBLOGJQcFRLS2JPeUpwUWJMMGFlaXBXYUhudXFDbnFzZVBSYkpNS3JoanQ1OFJnUDlJVlEwR1l1cm1iZnRHYThwU3RhTzRycWx4Q3JGZE1lOFpHWkEvYWVJZlplb0txQURjRnMvWmxZdkJVeUI0NUhCNHJKRzhqNnZadnhuYS9JMFZ5cDdua3NPNDJRcmU4NU9Wck5ENlNqZHdXa1V3Sk1jaXFMTGkyRDIyb0xWVm1yTGhWYjdzbTJhTDB3RVJ5ZFZQT0pPeE4zQnlVRko4VmRuTnZxdnlUcWtoRXFZeVRYTlROUzlEMnR6RS8vM3I4N2c0MFU5aXQrMFFCbWo3Ui9kUkc0eitaS3p1R3JhaFlJV3lQRHRaMFA0R1Y4VS9pMWhYYXFHdlNOcXh2S2dodW5GSjhaTVFFc1cyUWlmNC8veDNtY3dtZGJiYlREYk1SMjZqRnF0VmduRHY4VFhkNVp3Tkt4ZTU0M1ljNVZOMXhoZEMzamFQTUZueFRvVmFmQUtSQmFNQlQyeFp1eFoyZzhKa2hxYXd0REtwWUJiZHo5ZmhYUkhpdU1rNlYvUk45SzhlKzB1d0NRa29iYzU1S0c0L3pER1JoT2ZZT3dQYmJmQ2hLbjlFYTRkckxKUG9DUk1ya1Vsb1M5ZzN4ZnAxcVlIZXpPcjBnYlRJRFpOKzAvczhQaXpHQzhTVmtPTFVHcmtsMTJFaXh1eXRlQktnR1lqNmlIVk8vR3hWeC9IU2JRNTd2WWNNeFNGRlFERkMyNE50V0NuL0pKU2wxWU0vV3Jkc2NURHpPTmhDM3BnaFFFZTZCaEc4OXdPVU93dkxSQ2EzYWFSL25abXF3UWtBL2xRdHoyaTdodDUyaCs5Ly85Nkp6cWVyKzgvc0M2ZDZrSEZ5MUZ5RElUQjlmbVlxUm9ra3I0WGN3OGh6MTNiUU9WcDRTYWhRbVpxWFRLNE40TnFkVWVEVm83a3JCaVFrY1VCbTBPRHdMb0pnNWVEZHdNZ0FwcCtNOG9CNEU2ek9yWmRKa2lzczJiR3VyazAxdUtFMzUzbzkwdENUSVdMNFowQjZCZFNMYWRhQ2JnckFXTG5DNDJDTlFJdWFITEMwaWt5S0hKbER6OWphUEtkRXFRRE5veUNqdFJ3cDBCeUtoS3dyczlHYWUxS2NmMFp5ZTN5eTlRaUU4MkUwU3kzcVU1ZzRuNkVjQlpEaWJVK2FuVUxDZFBsU1hmQkhKYnhCZ2lyTXJUdEw3SjNxY3ZuUlltSUxhZVBXdlhzZFdMM0Ftc3Q2YlBQcTY3dmV1THFOaWxRaS92NngwMFJ1S1ZDelo2M3VtamE0alMvYlpWVGduWS9yTElkc2pXcWx5QnhHT2N6b01KUERZbGs0c1l1RjNPMGszelAwOHByUVN0MUcxRkh2S1Z2SVQreERLdWEwK2ltY1VQNUx6LzdES01hcmdwUTNGdEJnd2kwaEdyM3VoT01WRHU1cVVmYWZRblR1L0pCaUZCWGZUczdGK1l0clFBbWxPUENxakkvSk5xZmQxSFJqZi9rWG1SNCtBK0UwU0tGcS9EWGJOemo2bWRud29yMkQ2eVh6WGkvSWtORUVFVWJ3UTg5ZFoveEptNkQ0MmpqWDIrdCtBaTRkZVlpSU9IN000aTN6N0Uxb2xyRUJsU0JrS0Zzb3ZpSCsvQm9sNHVOcjMvTUh0S1NkeS9rbFNGMnYvZVRNZU9wMStLc3N0TjZQWGVhVU5YUlJadC8wZ3N0TU9Vb0lTZndvWlBRN2hNdnFheitkcjJiNk1QT3V0S1RYNnJLajVkY2pHZ0FpeEx1RS9SUUdIYTBwbDdhazJPYXJzMFBSVHZpc3JPempwZnlaYkxOWm43SklleUc5eHppV21yQjExd1A0RXUvZ0lBd08yNHNqVGNrWlQwRXpER1l3RjN1eHEyanRmOFIydkZHY0djbFVzUUkyQndDZ0dCU3ZhWFBOaCt4WUo3UGFkNWEwVFZNcEJNZTMxQUNpaVV4bGhCL0JDaGlRbEdsSUVYN1BROG0vT0Y1MkJhekN6YkJ1QkNsRkxYSTJqNkxueGF5SDJzdk81emJhNk1OclF1YzMwWXBXZGozV0VDelNySlk0Q2Zzdkxod0dWaksyaFBFTTYycWxFeUdXaHFBUHB5QVFHUmFrV0xRczNyZEkvVmMzK0tqajRsdXlRSjRlM1ovM05yNnE5NGJ0ZEhldXB6eHJTVjcyakJncXZ4MlpyQmFJRkRoQUd5MDlJczVTZVpzM01mQkxEbzY2VkpsTnYyejEwNmFzcHBhN3JsZ2NrZUNEN1ZaUS9XMkVTZXY2eFhKQTlhWlRRYklJbnh4OU51clpOSmx0ZkdoSEYyd0lDVEdjemNsbnNMWmMrR1ZvMTg5NGc5NzQ0TUk2MTJSc3RLb3UzeEp4VHVWTU1GM1lDY01LM2lVNGNxdkxvSzIxQUhKTDFSU2tTb01yMnY1Z3MwUVF2bXUzamJmeHRodHU0emVkcDgya2FVSmVZUVlQeGQwV29WVEZkSlNoSlB3c2h0amh0eEJEYWhqdVpEV2RGS2xwVjFyRFpNcHlLR08rbG1Ua1V2WG93dHk4cHlTeFg4SThlam5RM1hONnM2MXMxNlZnWkFWTi9DaVhxWXdrelpFRUd3TkczSFQ0Z1FGdEs4eWhzMDlWdWxYUzBkY3JkUU5UTFpsMFl0Z3ZxS0RFMVlRZERGRC90eWl0ZlRLbjRVVkpkeGloZitWazBMd1hDRmFlK0RseWhPV3grbWZHSlFubWxsaWtYSWJRNlVRL0cxWXlDQWFibHhHYi8xMDh0L1Z4WVlxNkVJUVFMQURxMTJ3QTM3ampTQUVRQTVFZncvbmlzY3ZSby9IYVZBMmJWajBWVjY0SmJaTXU2YklkRSs1ZXJVWmFTeisreXRvbXBpd25kVytiYllac3I4REdNdndvT05JdEdSOVJiTG9rb0owMFM0dlFxZGMzSmlXRHRnTDVVNmo2T1VxV1kvdGVveStnR2xOM2V5N2NIYWZBTmFKNzBkTi9nUzQwdVkvUDBiMDBqdzUyeWhpdlQ0aFR6UWNuc1kzMWJZMXc1QXJpcW5DV0VHSThnTVVPbGdjNy95QXJVd1JZOUdnMkkzNi9MazBSUytjQTRCaWpJZG1mbTV0dGozcVpBUTFGd2d6ZzZIeU9iaHBGR3hOVldjQXVQdVNSbUp0eldCR1FzTWd6OWZrRFFoL1c4T24raEhyUWdrWTJOT1lFODRxQS94bmVFZkVLM1hzTzgveHg0VkU4ZlY0WjBoa3V5ZHhCQjZWNGh1Q0lIeFQ4TlphTUNUbHA3aTdhRFJ2MEo4a3AwdGZuZzduY3puRmllN21Na3kxcGRqZGpxd2ZxYkQwR3pwK2RKZEJjQkszL1JpQ0FEbHlMQkp2TUtiN0JTWFRLejlLU3NxY2t4aFhtYWo1Y1V1TEZVZUF3MlJrN3FFeklHVkFJQT09IiwiaXYiOiI0YjUyZmY2YjM2OTMzZmI1NGMzMTdmY2UwNDYwNTRiOCIsInMiOiI5OGMyZWQ4NDM0YzAzZGU1In0%3D&public_key=35536E1E-65B4-4D96-9D97-6ADB7EFF8147&site=http%3A%2F%2Flocalhost%3A3000&userbrowser=Mozilla%2F5.0%20(Macintosh%3B%20Intel%20Mac%20OS%20X%2010.15%3B%20rv%3A109.0)%20Gecko%2F20100101%20Firefox%2F116.0&capi_version=1.5.4&capi_mode=lightbox&style_theme=default&rnd=0.7598961588666363"
		if payloadStr != "" {
			payload := &config.Payload{
				Payload:   payloadStr,
				UserAgent: gjson.New(res).Get("user_agent").String(),
				Created:   gtime.Now().Unix(),
			}
			config.Cache.Set(ctx, "payload", payload, 0)
			g.Log().Info(ctx, "从主节点获取payload成功")
			g.Dump(config.Cache.MustGet(ctx, "payload"))
			return
		} else {
			return gerror.New("从主节点获取payload失败")
		}
	} else {
		payloadStr := g.Cfg().MustGetWithEnv(ctx, "PAYLOAD").String()
		userAgent := g.Cfg().MustGetWithEnv(ctx, "USER_AGENT").String()
		if payloadStr != "" && userAgent != "" {
			payload := &config.Payload{
				Payload:   payloadStr,
				UserAgent: userAgent,
				Created:   gtime.Now().Unix(),
			}
			config.Cache.Set(ctx, "payload", payload, 0)
			g.Log().Info(ctx, "从配置文件获取payload成功")
			g.Dump(config.Cache.MustGet(ctx, "payload"))
			return
		}

	}
	return nil
}
