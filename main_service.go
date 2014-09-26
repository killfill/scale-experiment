package main

import (
	"fmt"
	"os"
	"time"

	"scale-experiment/appscaler"
	"scale-experiment/cf"
	"scale-experiment/scalerservice"

	"go-bro/broker"
	"go-bro/config"
)

var api = cf.NewApi(os.Getenv("DOMAIN"), os.Getenv("USER"), os.Getenv("PASS"), false, false)
var scaler = appscaler.New(api)

func main() {
	allApps := map[string]string{}

	go startLoop(allApps)
	startBroker(allApps)
}

//Check all apps that wants to be auto-scaled
func startLoop(list map[string]string) {
	for {
		for _, app := range list {
			checkApp(app)
		}
		time.Sleep(1 * time.Second)
	}
}

func startBroker(allApps map[string]string) {

	config := config.FromJson("config.json")
	b := broker.New(config.Username, config.Password, config.Plans)

	for _, serviceConfig := range config.Services {
		fmt.Printf("Registering Service Broker %+v \n", serviceConfig)
		service := scalerservice.New(allApps)
		b.RegisterService(serviceConfig.Id, service)
	}

	b.Listen(getAddr())
}

func mainWithoutBroker() {
	apps := []string{os.Getenv("GUID")}
	for {
		for _, app := range apps {
			checkApp(app)
		}
		time.Sleep(2 * time.Second)
	}
}

//Check if the app need to be scaled up or down
func checkApp(app string) {

	fmt.Println("=> Checking app", app)

	summary, err := scaler.GetSummary(app)
	if err != nil {
		fmt.Println("   Could not get stats: " + err.Error())
		return
	}
	// fmt.Printf("   summary: %+v\n", summary)

	newInstancesCount := scaler.ProposedInstances(summary)
	if newInstancesCount == summary.Instances {
		return
	}

	err = scaler.ScaleAppTo(app, newInstancesCount)
	if err != nil {
		fmt.Println("  error when scaling." + err.Error())
	}
	fmt.Printf("   Scaled to: %d. Summary: %+v\n", newInstancesCount, summary)
}

func getAddr() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":3000"
	}
	return ":" + port
}
