package bootstrap

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"micro-gin/config"
	"micro-gin/global"
	"micro-gin/utils"
	"strconv"
	"strings"
)

type KongServicesData map[string][]Service
type KongRoutersData map[string][]Router
type KongUpstreamsData map[string][]Upstream
type KongTargetsData map[string][]Targets

type KongKong struct {
	Kong config.Kong
}

// http://127.0.0.1:8001/services
type Service struct {
	Host              string      `json:"host"`
	UpdatedAt         int         `json:"updated_at"`
	Retries           int         `json:"retries"`
	Enabled           bool        `json:"enabled"`
	WriteTimeout      int         `json:"write_timeout"`
	Path              []string    `json:"path"`
	Protocol          string      `json:"protocol"`
	TlsVerify         interface{} `json:"tls_verify"`
	ClientCertificate interface{} `json:"client_certificate"`
	TlsVerifyDepth    interface{} `json:"tls_verify_depth"`
	CaCertificates    interface{} `json:"ca_certificates"`
	Name              string      `json:"name"`
	Id                string      `json:"id"`
	Port              int         `json:"port"`
	CreatedAt         int         `json:"created_at"`
	ConnectTimeout    int         `json:"connect_timeout"`
	Tags              []string    `json:"tags"`
	ReadTimeout       int         `json:"read_timeout"`
}

// http://127.0.0.1:8001/services/57b7ff79-5ce2-4bec-85d7-4b3da18082b9/routes
type Router struct {
	Paths                   []string    `json:"paths"`
	Methods                 []string    `json:"methods"`
	RequestBuffering        bool        `json:"request_buffering"`
	ResponseBuffering       bool        `json:"response_buffering"`
	UpdatedAt               int         `json:"updated_at"`
	HttpsRedirectStatusCode int         `json:"https_redirect_status_code"`
	PreserveHost            bool        `json:"preserve_host"`
	PathHandling            string      `json:"path_handling"`
	Snis                    interface{} `json:"snis"`
	Tags                    []string    `json:"tags"`
	Service                 struct {
		Id string `json:"id"`
	} `json:"service"`
	Protocols     []string    `json:"protocols"`
	Hosts         interface{} `json:"hosts"`
	Destinations  interface{} `json:"destinations"`
	Sources       interface{} `json:"sources"`
	Id            string      `json:"id"`
	StripPath     bool        `json:"strip_path"`
	Headers       interface{} `json:"headers"`
	CreatedAt     int         `json:"created_at"`
	RegexPriority int         `json:"regex_priority"`
	Name          string      `json:"name"`
}

// http://127.0.0.1:8001/upstreams
type Upstream struct {
	HashOnUriCapture       interface{} `json:"hash_on_uri_capture"`
	CreatedAt              int         `json:"created_at"`
	HashOnCookiePath       string      `json:"hash_on_cookie_path"`
	HashOn                 string      `json:"hash_on"`
	HashFallbackHeader     interface{} `json:"hash_fallback_header"`
	HashFallbackQueryArg   interface{} `json:"hash_fallback_query_arg"`
	HashFallbackUriCapture interface{} `json:"hash_fallback_uri_capture"`
	HostHeader             interface{} `json:"host_header"`
	HashOnQueryArg         interface{} `json:"hash_on_query_arg"`
	Tags                   []string    `json:"tags"`
	Healthchecks           struct {
		Active struct {
			Type    string `json:"type"`
			Healthy struct {
				HttpStatuses []int `json:"http_statuses"`
				Successes    int   `json:"successes"`
				Interval     int   `json:"interval"`
			} `json:"healthy"`
			Unhealthy struct {
				HttpFailures int   `json:"http_failures"`
				Interval     int   `json:"interval"`
				HttpStatuses []int `json:"http_statuses"`
				TcpFailures  int   `json:"tcp_failures"`
				Timeouts     int   `json:"timeouts"`
			} `json:"unhealthy"`
			HttpPath               string      `json:"http_path"`
			HttpsSni               interface{} `json:"https_sni"`
			Headers                interface{} `json:"headers"`
			HttpsVerifyCertificate bool        `json:"https_verify_certificate"`
			Timeout                int         `json:"timeout"`
			Concurrency            int         `json:"concurrency"`
		} `json:"active"`
		Threshold int `json:"threshold"`
		Passive   struct {
			Healthy struct {
				HttpStatuses []int `json:"http_statuses"`
				Successes    int   `json:"successes"`
			} `json:"healthy"`
			Unhealthy struct {
				HttpFailures int   `json:"http_failures"`
				HttpStatuses []int `json:"http_statuses"`
				TcpFailures  int   `json:"tcp_failures"`
				Timeouts     int   `json:"timeouts"`
			} `json:"unhealthy"`
			Type string `json:"type"`
		} `json:"passive"`
	} `json:"healthchecks"`
	Algorithm         string      `json:"algorithm"`
	ClientCertificate interface{} `json:"client_certificate"`
	Id                string      `json:"id"`
	HashOnHeader      interface{} `json:"hash_on_header"`
	HashFallback      string      `json:"hash_fallback"`
	Name              string      `json:"name"`
	HashOnCookie      interface{} `json:"hash_on_cookie"`
	Slots             int         `json:"slots"`
}

// http://127.0.0.1:8001/upstreams/9a354478-6be7-4680-9742-d8f44dab2d29/targets
type Targets struct {
	Weight   int    `json:"weight"`
	Id       string `json:"id"`
	Upstream struct {
		Id string `json:"id"`
	} `json:"upstream"`
	CreatedAt float64  `json:"created_at"`
	Tags      []string `json:"tags"`
	Target    string   `json:"target"`
}

// TODO 完善下面方法
func RegistorToKong(r *gin.Engine) {
	k := &KongKong{
		Kong: global.App.Config.Kong,
	}
	services, err := k.GetAllServices()
	if err != nil {
		panic("获取网关服务失败")
	}
	for _, serv := range services {
		routers, err := k.GetRoutersOfServ(serv.Id)
		if err != nil {
			panic("获取网关路由失败")
		}
		for _, router := range routers {
			k.DelRouter(serv.Id, router.Id)
			fmt.Println(router.Id)
		}
		k.DelService(serv.Id)
	}
	upstreams, err := k.GetAllUpstreams()
	if err != nil {
		panic("获取网关服务失败")
	}

	for _, upm := range upstreams {
		tags, err := k.GetTargetssOfUpsteam(upm.Id)
		if err != nil {
			panic("获取网关路由失败")
		}
		for _, tag := range tags {
			fmt.Println(tag.Id)
		}
	}

	//处理注册逻辑
	upstreamName := global.App.Config.App.AppName
	k.DelUpstream(upstreamName)
	//注册upstream
	u, err := k.AddUpstream(map[string]interface{}{"name": upstreamName})
	if err != nil {
		panic("注册upstream失败")
	}
	targetsConf := global.App.Config.MicroServ

	for _, tc := range targetsConf.Host {
		targdata := map[string]string{
			"target": tc + ":" + targetsConf.Port,
		}
		_, err = k.AddTarget(u.Id, targdata)
		if err != nil {
			panic("注册target失败")
		}
	}

	for _, v := range r.Routes() {
		if strings.Contains(v.Path, "/api") {
			servicePath := v.Path
			routerPath := strings.ReplaceAll(v.Path[1:], "/", "_")
			serviceName := routerPath + "_service"
			routerName := routerPath + "_router"
			//创建service
			serv := map[string]interface{}{
				"name":     serviceName,
				"host":     upstreamName,
				"path":     servicePath,
				"protocol": "http",
			}
			s, err := k.AddServ(serv)
			if err != nil {
				panic("add service fail")
			}

			//创建router
			router := map[string]interface{}{
				"name":      routerName,
				"paths":     []string{"/" + routerPath},
				"methods":   []string{"GET", "POST"},
				"hosts":     []string{"test-go.com"},
				"protocols": []string{"http"},
				"service":   map[string]string{"id": s.Id},
			}
			k.AddRouter(router)
		}
	}
}

func (k *KongKong) url() string {
	return k.Kong.Protocol + "://" + k.Kong.Host + ":" + k.Kong.Port
}

func (k *KongKong) AddTarget(uid string, target map[string]string) (t Targets, err error) {
	url := k.url() + k.Kong.UpstreamsPath + "/" + uid + k.Kong.TargetsPath
	data, _ := json.Marshal(target)
	req, err := utils.HttpDo(url, "POST", nil, nil, data)
	if err != nil {
		return t, errors.New("000")
	}
	type kongRes struct {
		Code int
		Data Targets
	}
	resultStruct := new(kongRes)
	utils.TwoJson(resultStruct, req)
	if resultStruct.Code == 200 || resultStruct.Code == 201 {
		return resultStruct.Data, nil
	}
	return t, errors.New(strconv.Itoa(resultStruct.Code))
}

func (k *KongKong) DelRouter(idServ, idRouter string) error {
	url := k.url() + k.Kong.ServicesPath + "/" + idServ + k.Kong.RoutesPath + "/" + idRouter
	_, err := utils.HttpDo(url, "DELETE", nil, nil, nil)
	return err
}

func (k *KongKong) DelService(id string) error {
	url := k.url() + k.Kong.ServicesPath + "/" + id
	_, err := utils.HttpDo(url, "DELETE", nil, nil, nil)
	return err
}

func (k *KongKong) DelUpstream(upstreamName string) error {
	url := k.url() + k.Kong.UpstreamsPath + "/" + upstreamName
	_, err := utils.HttpDo(url, "DELETE", nil, nil, nil)
	return err
}

func (k *KongKong) AddUpstream(upstreamName map[string]interface{}) (u Upstream, err error) {
	url := k.url() + k.Kong.UpstreamsPath
	data, _ := json.Marshal(upstreamName)
	req, err := utils.HttpDo(url, "POST", nil, nil, data)
	if err != nil {
		return u, errors.New("000")
	}
	type kongRes struct {
		Code int
		Data Upstream
	}
	resultStruct := new(kongRes)
	utils.TwoJson(resultStruct, req)
	if resultStruct.Code == 200 || resultStruct.Code == 201 {
		return resultStruct.Data, nil
	}
	return u, errors.New(strconv.Itoa(resultStruct.Code))
}

func (k *KongKong) GetServ(serviceName string) (s Service, err error) {
	url := k.url() + k.Kong.ServicesPath + "/" + serviceName
	req, err := utils.HttpDo(url, "GET", nil, nil, nil)
	if err != nil {
		return s, errors.New("000")
	}
	type kongRes struct {
		Code int
		Data Service
	}
	resultStruct := new(kongRes)
	utils.TwoJson(resultStruct, req)
	if resultStruct.Code == 200 {
		return resultStruct.Data, nil
	}
	return s, errors.New("fail")
}

func (k *KongKong) AddServ(service map[string]interface{}) (s Service, err error) {
	url := k.url() + k.Kong.ServicesPath
	data, _ := json.Marshal(service)
	req, err := utils.HttpDo(url, "POST", nil, nil, data)
	if err != nil {
		return s, errors.New("000")
	}
	type kongRes struct {
		Code int
		Data Service
	}
	resultStruct := new(kongRes)
	utils.TwoJson(resultStruct, req)
	if resultStruct.Code == 200 || resultStruct.Code == 201 {
		return resultStruct.Data, nil
	}
	return s, errors.New(strconv.Itoa(resultStruct.Code))
}

func (k *KongKong) AddRouter(router map[string]interface{}) (r Router, err error) {
	url := k.url() + k.Kong.RoutesPath
	data, _ := json.Marshal(router)
	req, err := utils.HttpDo(url, "POST", nil, nil, data)
	if err != nil {
		return r, errors.New("000")
	}
	type kongRes struct {
		Code int
		Data Router
	}
	resultStruct := new(kongRes)
	utils.TwoJson(resultStruct, req)
	if resultStruct.Code == 200 || resultStruct.Code == 201 {
		return resultStruct.Data, nil
	}
	return r, errors.New(strconv.Itoa(resultStruct.Code))
}

func (k *KongKong) GetAllServices() ([]Service, error) {
	url := k.url() + k.Kong.ServicesPath
	httpRequestResult, err := utils.HttpDo(url, "GET", nil, nil, nil)
	if err != nil {
		return nil, errors.New("request kong service api error")
	}
	type kongRes struct {
		Code int
		Data KongServicesData
	}
	resultStruct := new(kongRes)
	utils.TwoJson(resultStruct, httpRequestResult)
	if resultStruct.Code == 200 {
		return resultStruct.Data["data"], nil
	}
	return nil, errors.New("request kong service api error")
}

func (k *KongKong) GetRoutersOfServ(servId string) ([]Router, error) {
	url := k.url() + k.Kong.ServicesPath + "/" + servId + k.Kong.RoutesPath
	httpRequestResult, err := utils.HttpDo(url, "GET", nil, nil, nil)
	if err != nil {
		return nil, errors.New("request kong router api error")
	}
	type kongRes struct {
		Code int
		Data KongRoutersData
	}
	resultStruct := new(kongRes)
	utils.TwoJson(resultStruct, httpRequestResult)
	if resultStruct.Code == 200 {
		return resultStruct.Data["data"], nil
	}
	return nil, errors.New("request kong routersapi error")
}

func (k *KongKong) GetAllUpstreams() ([]Upstream, error) {
	url := k.url() + k.Kong.UpstreamsPath
	httpRequestResult, err := utils.HttpDo(url, "GET", nil, nil, nil)
	if err != nil {
		return nil, errors.New("request kong upsteams api error")
	}
	type kongRes struct {
		Code int
		Data KongUpstreamsData
	}
	resultStruct := new(kongRes)
	utils.TwoJson(resultStruct, httpRequestResult)
	if resultStruct.Code == 200 {
		return resultStruct.Data["data"], nil
	}
	return nil, errors.New("request kong upsteams error")
}

func (k *KongKong) GetTargetssOfUpsteam(upsteamId string) ([]Targets, error) {
	url := k.url() + k.Kong.UpstreamsPath + "/" + upsteamId + k.Kong.TargetsPath
	httpRequestResult, err := utils.HttpDo(url, "GET", nil, nil, nil)
	if err != nil {
		return nil, errors.New("request kong router api error")
	}
	type kongRes struct {
		Code int
		Data KongTargetsData
	}
	resultStruct := new(kongRes)
	utils.TwoJson(resultStruct, httpRequestResult)
	if resultStruct.Code == 200 {
		return resultStruct.Data["data"], nil
	}
	return nil, errors.New("request kong routersapi error")
}
