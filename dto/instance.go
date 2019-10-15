package dto

type Instance struct {
	InstanceId string `json:"instanceId"`
	HostName   string `json:"hostName"`
	App        string `json:"app"`
	IpAddr     string `json:"ipAddr"`
	Status     string `json:"status"`
	Port       Port   `json:"port"`
}
