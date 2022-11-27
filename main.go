package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"strings"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"github.com/cip8/autoname"
)

type structs struct {
	Captcha  interface{}
	solution interface{}
	Time     float64
	Dcfd     string
	Sdcfd    string
	Xprops   string
	Xconst   string
	Finger   string
	randusr  string
	Config struct {
		Proxy 	 string `json:"proxy"`
		ApiKey 	 string `json:"capkey"`
	} `json:"Config"`
}

var (
	c = "\033[36m"
	r = "\033[39m"
	username string
	proxy = config().Config.Proxy
	apikey = config().Config.ApiKey
	p, _ = url.Parse("http://" + proxy)
)



func generate(username string, invite string) error {
	for true {
		Client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					MaxVersion: tls.VersionTLS13,
				},
				Proxy: http.ProxyURL(p),
			},
		}
		if username == "!random" {
			username = names().randusr
		}
		xfingerprint := fingerprints().Finger
		payload := map[string]string{
			"captcha_key": "",
			"consent": "true",
			"fingerprint": xfingerprint,
			"gift_code_sku_id": "",
			"invite": invite,
			"username": username,

		}
		xp,_ := json.Marshal(payload)
		req,_ := http.NewRequest("POST", "https://discord.com/api/v9/auth/register", bytes.NewBuffer(xp))
		Cookie := Build_cookie()
		Cookies := "__dcfduid=" + Cookie.Dcfd + "; " + "__sdcfduid=" + Cookie.Sdcfd + "; "
		for x,o := range map[string]string{
			"accept":" */*",
			"accept-encoding": "gzip, deflate, br",
			"accept-language": "en-US,en;q=0.9",
			"content-type": "application/json",
			"cookie": Cookies,
			"origin": "https://discord.com",
			"referer": "https://discord.com/invite/"+invite+"",
			"sec-ch-ua": `Google Chrome";v="105", "Not)A;Brand";v="8", "Chromium";v="105`,
			"sec-ch-ua-mobile": "?0",
			"sec-ch-ua-platform": "Windows",
			"sec-fetch-dest": "empty",
			"sec-fetch-mode": "cors",
			"sec-fetch-site": "same-origin",
			"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
			"x-debug-options": "bugReporterEnabled",
			"x-discord-locale": "en-US",
			"x-fingerprint": xfingerprint,
			"x-super-properties": Build_Xheader().Xprops,
		} {
			req.Header.Set(x,o)
		}
		resp, _ := Client.Do(req)
		if resp.StatusCode == 400 {
			scp := solve_cap()
			solved := scp.Captcha.(string)
			payload := map[string]string{
				"captcha_key": solved,
				"consent": "true",
				"fingerprint": xfingerprint,
				"gift_code_sku_id": "",
				"invite": invite,
				"username": username,
		
			}
			xp,_ := json.Marshal(payload)
			req,err := http.NewRequest("POST", "https://discord.com/api/v9/auth/register", bytes.NewBuffer(xp))
			if err != nil {
				log.Fatal(err)
				continue
			}	
			for x,o := range map[string]string{
				"accept":" */*",
				"accept-encoding": "gzip, deflate, br",
				"accept-language": "en-US,en;q=0.9",
				"content-type": "application/json",
				"cookie": Cookies,
				"origin": "https://discord.com",
				"referer": "https://discord.com/invite/"+invite+"",
				"sec-ch-ua": `Google Chrome";v="105", "Not)A;Brand";v="8", "Chromium";v="105`,
				"sec-ch-ua-mobile": "?0",
				"sec-ch-ua-platform": "Windows",
				"sec-fetch-dest": "empty",
				"sec-fetch-mode": "cors",
				"sec-fetch-site": "same-origin",
				"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
				"x-debug-options": "bugReporterEnabled",
				"x-discord-locale": "en-US",
				"x-fingerprint": xfingerprint,
				"x-super-properties": Build_Xheader().Xprops,
			} {
				req.Header.Set(x,o)
			}
			resp, err := Client.Do(req)
			if err != nil {
				log.Fatal(err)
				continue
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
				continue
			}
			var data map[string]interface{}
			if !strings.Contains(string(body), "token") {
				return errors.New(string(body))
			}
			err = json.Unmarshal(body, &data)
			if err != nil {
				log.Fatal(err)
				continue
			}
			token := data["token"].(string)
			if resp.StatusCode == 201 {
				fmt.Println("("+c+"+"+r+") Token "+c+":"+r+" ", token)
				f, err := os.OpenFile("tokens.txt", os.O_RDWR|os.O_APPEND, 0660)
				if err != nil {
					log.Fatal(err)
					continue
				}
				defer f.Close()
				_, ers := f.WriteString(token + "\n")
				if ers != nil {
					log.Fatal(ers)
					continue
				}
			}
		} else {
			continue
		}

	}
	return nil
}





// func science(invite string, xfingerprint string) {
// 	Client := &http.Client{
// 		Transport: &http.Transport{
// 			TLSClientConfig: &tls.Config{
// 				MaxVersion: tls.VersionTLS13,
// 			},
// 			Proxy: http.ProxyURL(p),
// 		},
// 	}
// 	Cookie := Build_cookie()
// 	Cookies := "__dcfduid=" + Cookie.Dcfd + "; " + "__sdcfduid=" + Cookie.Sdcfd + "; "
// 	data := map[string]interface{}{
// 		"events": map[string]interface{}{
// 			"type": "network_action_user_register",
// 			"properties": map[string]interface{}{
// 				"client_track_timestamp": 1665494233902,
// 				"status_code": 400,
// 				"url": "/auth/register",
// 				"request_method": "post",
// 				"invite_code": "spammer",
// 				"client_performance_memory": 0,
// 				"accessibility_features": 256,
// 				"rendered_locale": "en-US",
// 				"accessibility_support_enabled": false,
// 				"client_uuid": "WgCEMYcZSQ6Ke0t43xkxx4MBAAAKAAAA",
// 				"client_send_timestamp": 1665494233922,
			
// 			"type": "impression_user_registration",
// 			"properties": map[string]interface{}{
// 				"client_track_timestamp": 1665494233912,
// 				"impression_type": "view",
// 				"impression_group": "user_registration_flow",
// 				"step": "captcha",
// 				"location_section": "impression_user_registration",
// 				"client_performance_memory": 0,
// 				"accessibility_features": 256,
// 				"rendered_locale": "en-US",
// 				"accessibility_support_enabled": false,
// 				"client_uuid": "WgCEMYcZSQ6Ke0t43xkxx4MBAAALAAAA",
// 				"client_send_timestamp": 1665494233922,
// 		  		},
// 			},
// 		},
// 	}
// 	xp,_ := json.Marshal(data)
// 	req, err := http.NewRequest("POST", "https://discord.com/api/v9/science", bytes.NewBuffer(xp))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for x,o := range map[string]string{
// 		"accept":" */*",
// 		"accept-encoding": "gzip, deflate, br",
// 		"accept-language": "en-US,en;q=0.9",
// 		"content-type": "application/json",
// 		"cookie": Cookies,
// 		"origin": "https://discord.com",
// 		"referer": "https://discord.com/invite/"+invite+"",
// 		"sec-ch-ua": `Google Chrome";v="105", "Not)A;Brand";v="8", "Chromium";v="105`,
// 		"sec-ch-ua-mobile": "?0",
// 		"sec-ch-ua-platform": "Windows",
// 		"sec-fetch-dest": "empty",
// 		"sec-fetch-mode": "cors",
// 		"sec-fetch-site": "same-origin",
// 		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
// 		"x-debug-options": "bugReporterEnabled",
// 		"x-discord-locale": "en-US",
// 		"x-fingerprint": xfingerprint,
// 		"x-super-properties": Build_Xheader().Xprops,
// 	} {
// 		req.Header.Set(x,o)
// 	}
// 	resp, err := Client.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if resp.StatusCode == 204 {

// 	} else {

// 	}


// }




 


func solve_cap() structs {
	Client := http.Client{}
	payload := map[string]interface{}{
        "clientKey": apikey,
        "task": map[string]interface{}{
			"type": "HCaptchaTaskProxyless",
			"userAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36 Edg/92.0.902.73",
			"websiteKey": "4c672d35-0701-42b2-88c3-78380b0db560",
			"websiteURL": "https://discord.com/",
		},
	}
	xp,_ := json.Marshal(payload)
	req,_ := http.NewRequest("POST", "https://api.capmonster.cloud/createTask", bytes.NewBuffer(xp))
	resp, err := Client.Do(req)
	if err != nil {
		log.Fatal(err)
	} 
	cap := structs{}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 200 {
		taskid := data["taskId"]
		xpayload := map[string]interface{}{
			"clientKey": apikey,
			"taskId": taskid,
		}
		xpy,_ := json.Marshal(xpayload)
		for true {
			req, err := http.NewRequest("POST", "https://api.capmonster.cloud/getTaskResult", bytes.NewBuffer(xpy))
			if err != nil {
				log.Fatal(err)
				continue
			}	
			resp, err := Client.Do(req)		
			if err != nil {
				log.Fatal(err)
				continue
			}		
			defer resp.Body.Close()
			responseBody := make(map[string]interface{})
			json.NewDecoder(resp.Body).Decode(&responseBody)
			status := responseBody["status"]
			if status == "ready" {	
				cap.Captcha = responseBody["solution"].(map[string]interface{})["gRecaptchaResponse"].(string)
				break
			} else if status == "processing" {
				continue
			} else {
				fmt.Println("[ERR] | ", data)
			}
		
		//fmt.Println("Solving Captcha..")
		}
	}
	return cap
}



func fingerprints() structs {
	Client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MaxVersion: tls.VersionTLS13,
			},
			Proxy: http.ProxyURL(p),
		},
	}
	xf := structs{}
	req,err := http.NewRequest("GET", "https://discord.com/api/v9/experiments", nil)
	req.Close = true
	if err != nil {
		log.Fatal(err)
	}	
	resp, err := Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}	
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)
	xf.Finger = data["fingerprint"].(string)

	return xf
}

func Build_Xheader() structs {
	Xheader := structs{}
	xconststr := `{"location":"Invite Button Embed","location_guild_id":null,"location_channel_id":"","location_channel_type":3,"location_message_id":""}`
	xpropsstr := `{"os":"Windows","browser":"Chrome","device":"","system_locale":"en-US","browser_user_agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36","browser_version":"105.0.0.0","os_version":"10","referrer":"","referring_domain":"","referrer_current":"","referring_domain_current":"","release_channel":"stable","client_build_number":151638,"client_event_source":null}`
	Xheader.Xconst = base64.StdEncoding.EncodeToString([]byte(xconststr))
	Xheader.Xprops = base64.StdEncoding.EncodeToString([]byte(xpropsstr))
	return Xheader
}

func names() structs {
	name := structs{}
	name.randusr = autoname.Generate("")
	return name
}

func config() structs {
	var config structs
	conf, err := os.Open("config.json")
	defer conf.Close()
	if err != nil {
		log.Fatal(err)
	}	
	xp := json.NewDecoder(conf)
	xp.Decode(&config)
	return config

}



func Build_cookie() structs {
	Client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MaxVersion: tls.VersionTLS13,
			},
			Proxy: http.ProxyURL(p),
		},
	}
	req, err := http.NewRequest("GET", "https://discord.com", nil)
	req.Close = true
	if err != nil {
		log.Fatal(err)
	}	
	for x,o := range map[string]string{
		"accept":" */*",
		"accept-encoding": "gzip, deflate, br",
		"accept-language": "en-US,en;q=0.9",
		"content-type": "application/json",
		"sec-ch-ua": `Google Chrome";v="105", "Not)A;Brand";v="8", "Chromium";v="105`,
		"sec-ch-ua-mobile": "?0",
		"sec-ch-ua-platform": "Windows",
		"sec-fetch-dest": "empty",
		"sec-fetch-mode": "cors",
		"sec-fetch-site": "same-origin",
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
	} {
		req.Header.Set(x,o)
	}
	resp, err := Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	Cookie := structs{}
	if resp.Cookies() != nil {
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "__dcfduid" {
				Cookie.Dcfd = cookie.Value
			}
			if cookie.Name == "__sdcfduid" {
				Cookie.Sdcfd = cookie.Value
			}
		}
	}
	return Cookie
}

func cls() {
	cmd := exec.Command("cmd", "/c", "cls") 
	cmd.Stdout = os.Stdout
	cmd.Run()
}









func main() {
	cls()
	var threads int
	var user, inv string
	logo := `
	____`+c+`_____`+c+`___`+r+`___     `+r+`____`+c+`_____`+r+`____`+c+`_____`+r+`___`+c+`_   __
	`+r+`__  `+c+`____/`+r+`_  `+c+`__ \    `+r+`__  `+c+`____/`+r+`__  `+c+`____/`+r+`__`+c+`  | / /
	`+r+`_  `+c+`/ __ `+r+`_  `+c+`/ / /    `+r+`_  `+c+`/ __ `+r+`__  `+c+`__/  `+r+`__   `+c+`|/ / 
	/ /_/ / / /_/ /     / /_/ / `+r+`_  `+c+`/___  `+r+`_`+c+`  /|  /  
	\____/  \____/      \____/  /_____/  /_/ |_/   

	[`+r+`GO MEMBER BOOSTER`+c+`]		 $`+r+`YABOI
	`
	fmt.Println(logo)
	fmt.Print("("+c+"-"+r+") Username: ")
	fmt.Scanln(&user)
	fmt.Print("("+c+"-"+r+") Invite: ")	
	fmt.Scanln(&inv)
	fmt.Print("("+c+"-"+r+") Threads: ")
	fmt.Scanln(&threads)
	for i := 0; i < threads; i++ {
		go func() {
			for {
				generate(user, inv)
			}
		}()
	}
	select{}
}


