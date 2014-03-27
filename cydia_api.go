/*
The MIT License (MIT)

Copyright (c) 2014 Evan Swick

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cydia;

import(
	"net/http"
	"time"
	"fmt"
	"crypto/hmac"
	"crypto/sha1"
	"io/ioutil"
	"io"
	"encoding/base64"
	"strings"
	"net/url"
	"errors"
)

/* === Struct to hold API response info === */
type CydiaPurchaseInfo struct{
	Response url.Values;
}

func (r CydiaPurchaseInfo) PurchaseComplete() bool{
	if(r.Response.Get("state") == "completed"){
		return true;
	}
	return false;
}

/* === === */

func urlsafe_b64encode(b64 string) string{
	result := strings.Replace(b64, "=", "", -1);
	result = strings.Replace(result, "/", "_", -1);
	result = strings.Replace(result, "+", "-", -1);

	return result;
}

func get_hmac(query string, key string) string{

	mac := hmac.New(sha1.New, []byte(key));
	io.WriteString(mac, query);

	signature := urlsafe_b64encode(base64.StdEncoding.EncodeToString(mac.Sum(nil)));

	return signature;
}

func buildQuery(udid string, package_id string, dev string, key string) string{
	/* This must be in alphabetical order */
	query := fmt.Sprintf("api=store-0.9&device=%s&mode=local&nonce=%d&package=%s&timestamp=%d&vendor=%s", udid, time.Now().Unix(), package_id, time.Now().Unix(), dev);

	finalQuery := fmt.Sprintf("%s&signature=%s", query, get_hmac(query, key));

	return finalQuery;
}

func check_response_signature(response url.Values, key string) bool{
	/* Remove 'signature' value from response.
	   The response appears to already be in the correct order to be hashed, no need to mess with it. */
	signature := response.Get("signature");
	response.Del("signature");

	hmac := get_hmac(response.Encode(), key);

	/* Add 'signature' value back to the response because...it just seems right. */
	response.Set("signature", signature);

	return (urlsafe_b64encode(hmac) == signature);
}


func CheckCydiaPurchase(udid string, package_id string, dev string, apikey string) (*CydiaPurchaseInfo, error){

	/* Build the query string */
	query := buildQuery(udid, package_id, dev, apikey);

	client := &http.Client{}

	/* Send the query to the Cydia API and get response */
	req, err := http.NewRequest("GET", fmt.Sprintf("http://cydia.saurik.com/api/check?%s", query), nil);
	req.URL = &url.URL{
		Scheme: "http",
		Host: "cydia.saurik.com",
		Opaque: fmt.Sprintf("//cydia.saurik.com/api/check?%s", query),
	}

	//TODO: Proper error handling.
	resp, err := client.Do(req);

	if(err != nil){
		return nil, err;
	}

	/* Get body text */
	body, err := ioutil.ReadAll(resp.Body);

	if(err != nil){
		return nil, err;
	}

	/* Parse reponse values */
	values, err := url.ParseQuery(string(body));

	if(err != nil){
		return nil, err;
	}

	if(!check_response_signature(values, apikey)){
		return nil, errors.New("Invalid response signature.");
	}

	responseInfo := CydiaPurchaseInfo{values};

	return &responseInfo, nil;
}