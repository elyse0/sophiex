package sites_extractor

import (
	"encoding/json"
	"os"
	"os/exec"
	"sophiex/internal/logger"
)

type ytDlpDownloadableFormat struct {
	Id          string            `json:"format_id"`
	Url         string            `json:"url"`
	Protocol    string            `json:"protocol"`
	HttpHeaders map[string]string `json:"http_headers"`
}

type InfoDict struct {
	Id               string                    `json:"id"`
	Title            string                    `json:"title"`
	RequestedFormats []ytDlpDownloadableFormat `json:"requested_formats"`
	Url              string                    `json:"url"`
	FormatId         string                    `json:"format_id"`
	Protocol         string                    `json:"protocol"`
}

func getProtocol(protocol string) DownloadProtocol {
	if protocol == "https" || protocol == "http" {
		return Http
	}
	if protocol == "m3u8" || protocol == "m3u8_native" {
		return Hls
	}

	panic(protocol)
}

func GetDownloadableFormats(url string) ([]DownloadableFormat, error) {
	ytDlp := exec.Command("yt-dlp", url, "--skip-download", "-S", "proto", "--print-json")
	logger.Log.Debug("%v", ytDlp.Args)

	ytDlp.Stderr = os.Stderr

	ytDlpOutput, err := ytDlp.Output()
	if err != nil {
		return nil, err
	}

	var infoDict InfoDict
	err = json.Unmarshal(ytDlpOutput, &infoDict)
	if err != nil {
		return nil, err
	}

	if infoDict.Url != "" {
		return []DownloadableFormat{
			{
				Id:       infoDict.FormatId,
				Url:      infoDict.Url,
				Protocol: getProtocol(infoDict.Protocol),
			},
		}, nil
	}

	var downloadableFormats []DownloadableFormat
	for _, format := range infoDict.RequestedFormats {
		downloadableFormats = append(downloadableFormats, DownloadableFormat{
			Id:          format.Id,
			Url:         format.Url,
			Protocol:    getProtocol(format.Protocol),
			HttpHeaders: format.HttpHeaders,
		})
	}

	return downloadableFormats, nil
}
