package main

// import (
// 	"fmt"
// 	"os"
// 	"time"

// 	"scale-experiment/appscaler"
// 	"scale-experiment/cf"
// )

// var api = cf.NewApi(os.Getenv("DOMAIN"), os.Getenv("USER"), os.Getenv("PASS"), false, false)
// var scaler = appscaler.New(api)

// func main() {
// 	apps := []string{os.Getenv("GUID")}
// 	for {
// 		for _, app := range apps {
// 			checkApp(app)
// 		}
// 		time.Sleep(2 * time.Second)
// 	}
// }

// //Check if the app need to be scaled up or down
// func checkApp(app string) {

// 	fmt.Println("=> Checking app", app)

// 	summary, err := scaler.GetSummary(app)
// 	if err != nil {
// 		fmt.Println("   Could not get stats: " + err.Error())
// 		return
// 	}
// 	// fmt.Printf("   summary: %+v\n", summary)

// 	newInstancesCount := scaler.ProposedInstances(summary)
// 	if newInstancesCount == summary.Instances {
// 		return
// 	}

// 	err = scaler.ScaleAppTo(app, newInstancesCount)
// 	if err != nil {
// 		fmt.Println("  error when scaling." + err.Error())
// 	}
// 	fmt.Printf("   Scaled to: %d. Summary: %+v\n", newInstancesCount, summary)
// }
