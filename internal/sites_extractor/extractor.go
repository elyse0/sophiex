package sites_extractor

type DownloadProtocol string

const (
	Http DownloadProtocol = "https"
	Hls                   = "hls"
)

type DownloadableFormat struct {
	Id          string            `json:"id"`
	Url         string            `json:"url"`
	Protocol    DownloadProtocol  `json:"protocol"`
	HttpHeaders map[string]string `json:"httpHeaders"`
}
