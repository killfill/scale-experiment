package main

import (
	"fmt"
	"os"

	"scale-experiment/appscaler"
	"scale-experiment/cf"
	"time"
)

func main() {

	apps := []string{os.Getenv("GUID")}

	api := cf.NewApi(os.Getenv("DOMAIN"), os.Getenv("USER"), os.Getenv("PASS"), false, false)
	scaler := appscaler.New(api)

	for {

		for _, app := range apps {

			fmt.Println("=> Checking app", app)

			summary, err := scaler.GetSummary(app)
			if err != nil {
				fmt.Println("   Could not get stats: " + err.Error())
				continue
			}
			// fmt.Printf("   summary: %+v\n", summary)

			newInstancesCount := scaler.ProposedInstances(summary)
			if newInstancesCount == summary.Instances {
				continue
			}

			err = scaler.ScaleAppTo(app, newInstancesCount)
			if err != nil {
				fmt.Println("  error when scaling." + err.Error())
			}
			fmt.Printf("   Scaled to: %d. Summary: %+v\n", newInstancesCount, summary)

		}

		time.Sleep(2 * time.Second)
	}
}
