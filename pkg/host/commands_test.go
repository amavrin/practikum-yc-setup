package host

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRedirectPort(t *testing.T) {
	u := `https://auth.yandex.cloud/oauth/authorize?client_id=yc.oauth.public-sdk&code_challenge=xxxxxxxx&code_challenge_method=S256&redirect_uri=http%3A%2F%2F127.0.0.1%3A40567%2Fauth%2Fcallback&response_type=code&scope=openid&state=zzzz-zzzz&yc_federation_hint=ffffffffffffffff`
	urlStr := "junk before " + u + " junk after"
	port, url, err := getRedirectPort(urlStr)
	require.NoError(t, err)
	require.Equal(t, 40567, port)
	require.Equal(t, u, url)
}
