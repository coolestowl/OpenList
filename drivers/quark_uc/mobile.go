package quark

import (
	"errors"
	"net/url"
	"strings"
)

const mobilePlayPath = "/1/clouddrive/file/v2/play/project"

// mobilePlayURL uses the device and authentication parameters from a captured
// Quark app request for the video playback endpoint.
func mobilePlayURL(raw string) (string, error) {
	if len(raw) > 32768 || strings.ContainsAny(raw, "\r\n\x00") {
		return "", errors.New("invalid Quark mobile request URL")
	}
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Scheme != "https" || u.Host != "drive-m.quark.cn" || u.User != nil || u.Fragment != "" || !strings.HasPrefix(u.Path, "/1/clouddrive/") {
		return "", errors.New("mobile request URL must be an HTTPS drive-m.quark.cn/1/clouddrive URL")
	}
	query, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", errors.New("invalid Quark mobile request parameters")
	}
	for _, key := range []string{"kps", "sign", "vcode", "ut"} {
		if len(query[key]) != 1 || strings.TrimSpace(query.Get(key)) == "" {
			return "", errors.New("mobile request URL requires kps, sign, vcode and ut")
		}
	}
	for _, key := range []string{"fid", "fids", "pdir_fid", "share_id", "share_token", "share_fid_token", "resolutions", "supports", "_page", "_size", "_fetch_total", "fetch_all_file", "fetch_risk_file_name", "_sort", "_t"} {
		query.Del(key)
	}
	if query.Get("pr") == "" {
		query.Set("pr", "ucpro")
	}
	if query.Get("fr") == "" {
		query.Set("fr", "android")
	}
	u.Path = mobilePlayPath
	u.RawPath = ""
	u.RawQuery = query.Encode()
	return u.String(), nil
}
