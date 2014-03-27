Golang Cydia API
============

The Golang Cydia API allows you to check the status of a purchased package for an iOS device. It is intended for use in anti-piracy server backends.

##Installation

Install the API to your system using the following command:

    $ go get github.com/eswick/cydia_go_api

##Usage
First, import the package.

    import "github.com/eswick/cydia_go_api"

The only public function in the API is CheckCydiaPurchase, as follows:

    func CheckCydiaPurchase(udid string, package_id string, dev string, apikey string) (*CydiaPurchaseInfo, error);

### Arguments

&nbsp;**udid**

&nbsp;&nbsp; The UDID of the device for which purchase information is to be retrieved.

&nbsp;**package_id**

&nbsp;&nbsp; The package identifier of the package to check.

&nbsp;**dev**

&nbsp;&nbsp; Your vendor ID. (obtained from Cydia's web interface)

&nbsp;**apikey**

&nbsp;&nbsp; Your API key. (also obtained from Cydia's web interface)

###Return Value

`CheckCydiaPurchase` returns a `CydiaPurchaseInfo` struct containing data from the Cydia API response, or nil + an error if an error occurred.

####Example

    info, err := cydia.CheckCydiaPurchase("udid_here", "us.kanyon.beacon", "eswick", "api_key_here");

    if(err != nil){
    	fmt.Println("Error checking Cydia API.");
    	return;
    }

    if(info.PurchaseComplete()){
        fmt.Println("Purchase complete!");
    }else{
        fmt.Println("Purchase incomplete.");
    }
