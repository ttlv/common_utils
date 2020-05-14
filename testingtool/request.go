package testingtool

import (
	"io/ioutil"
	"net/url"
	"testing"

	httpclient "github.com/ddliu/go-httpclient"
	"github.com/golang/mock/gomock"
	"gitlab.com/hashtech/common/services/coinx"
	"gitlab.com/hashtech/common/services/coinx/mock"
)

var TestServerDomain = ""

func Get(url string) (string, *httpclient.Response) {
	resp, _ := httpclient.NewHttpClient().Get(TestServerDomain + url)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), resp
}

func SignedGet(t *testing.T, coinxClient *mock_coinx.MockCoinxClient, url string, token string, memberInfo *coinx.GetMemberInfoResp) string {
	if memberInfo == nil {
		memberInfo = &coinx.GetMemberInfoResp{}
	}
	if token == "" {
		token = "imrobot"
	}
	coinxClient.EXPECT().GetMemberInfo(gomock.Any(), &coinx.GetMemberInfoReq{AccessToken: token}).Return(memberInfo, nil).AnyTimes()

	client := httpclient.NewHttpClient()
	client.Headers = make(map[string]string)
	client.Headers["Authorization"] = "Token token=" + token
	resp, _ := client.Get(TestServerDomain + url)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func Post(url string, values url.Values) (string, *httpclient.Response) {
	params := make(map[string]string)
	for k, value := range values {
		params[k] = value[0]
	}

	resp, _ := httpclient.Post(TestServerDomain+url, params)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), resp
}

func SignedPost(t *testing.T, coinxClient *mock_coinx.MockCoinxClient, url string, values url.Values, token string, memberInfo *coinx.GetMemberInfoResp) string {
	if memberInfo == nil {
		memberInfo = &coinx.GetMemberInfoResp{}
	}
	if token == "" {
		token = "imrobot"
	}
	coinxClient.EXPECT().GetMemberInfo(gomock.Any(), &coinx.GetMemberInfoReq{AccessToken: token}).Return(memberInfo, nil).AnyTimes()
	params := make(map[string]string)
	for k, value := range values {
		params[k] = value[0]
	}

	client := httpclient.NewHttpClient()
	client.Headers = make(map[string]string)
	client.Headers["Authorization"] = "Token token=" + token
	resp, _ := client.Post(TestServerDomain+url, params)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func Put(url string, values url.Values) (string, *httpclient.Response) {
	params := make(map[string]string)
	params["_method"] = "PUT"
	for k, value := range values {
		params[k] = value[0]
	}

	resp, _ := httpclient.Post(TestServerDomain+url, params)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), resp
}

func SignedPut(t *testing.T, coinxClient *mock_coinx.MockCoinxClient, url string, values url.Values, token string, memberInfo *coinx.GetMemberInfoResp) string {
	values["_method"] = []string{"PUT"}
	return SignedPost(t, coinxClient, url, values, token, memberInfo)
}
