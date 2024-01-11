package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

// Setup Color Palette
const (
	Black   = "\033[1;30m"
	Red     = "\033[1;31m"
	Green   = "\033[1;32m"
	Yellow  = "\033[1;33m"
	Blue    = "\033[1;34m"
	Magenta = "\033[1;35m"
	Cyan    = "\033[1;36m"
	White   = "\033[1;37m"
	Reset   = "\033[0m"
)

// Regions
var regions = []string{
	"Europe",
	"Americas",
	"Asia",
	"Australia",
}

// URLs for each region. Key is the region name, value is a array of objects with the following keys:
// - url: URL to test
// - url_friendlyname: Friendly name for the URL. Used in the output.

var urls = map[string][]map[string]string{
	"Europe": {
		{
			"url":              "uk.ingest.vrcdn.live",
			"url_friendlyname": "UK Ingest (England)",
		},
		{
			"url":              "de.ingest.vrcdn.live",
			"url_friendlyname": "DE Ingest (Germany)",
		},
	},
	"Americas": {
		{
			"url":              "use.ingest.vrcdn.live",
			"url_friendlyname": "USE Ingest (United States East)",
		},
		{
			"url":              "usc.ingest.vrcdn.live",
			"url_friendlyname": "USC Ingest (United States Central)",
		},
		{
			"url":              "usw.ingest.vrcdn.live",
			"url_friendlyname": "USW Ingest (United States West)",
		},
	},
	"Asia": {
		{
			"url":              "jpe.ingest.vrcdn.live",
			"url_friendlyname": "JPE Ingest (Japan East)",
		},
		{
			"url":              "jpw.ingest.vrcdn.live",
			"url_friendlyname": "JPW Ingest (Japan West)",
		},
	},
	"Australia": {
		{
			"url":              "au.ingest.vrcdn.live",
			"url_friendlyname": "AUS Ingest (Sydney)",
		},
	},
}

// Define our structs
type dynamicIngestResult struct {
	packetsTransmitted int
	packetsReceived    int
	packetLoss         float64
	minPing            float64
	maxPing            float64
	avgPing            float64
}

type regionResult struct {
	region             string
	areaName           string
	packetsTransmitted int
	packetsReceived    int
	packetLoss         float64
	minPing            float64
	maxPing            float64
	avgPing            float64
}

func testRegionPing(region string) []regionResult {
	// Get URLs for the region.
	regionUrls := urls[region]

	// Make the results array for the regions, using the regionRequest struct.
	var regionResults []regionResult = make([]regionResult, len(regionUrls))

	// Print region name.
	fmt.Printf("%s[ping_test] Testing region %s...%s\n", Yellow, region, Reset)

	// Test each URL.
	for _, url := range regionUrls {
		fmt.Printf("[ping_test] Testing endpoint %s in region %s with sample rate 10\n", url["url_friendlyname"], region)

		pinger, err := probing.NewPinger(url["url"])
		if err != nil {
			panic(err)
		}

		// Listen for Ctrl-C.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for range c {
				pinger.Stop()
				fmt.Printf("\n%s[ping_test] Stopped.%s\n", Red, Reset)
				os.Exit(0)
			}
		}()

		pinger.Count = 10
		pinger.Interval = 100 * time.Millisecond
		pinger.Timeout = 1 * time.Second

		pinger.OnSend = func(pkt *probing.Packet) {
			// Increment the packetsTransmitted counter.
			regionResults[len(regionResults)-1].packetsTransmitted++
		}

		pinger.OnRecv = func(pkt *probing.Packet) {
			// Increment the packetsReceived counter.
			regionResults[len(regionResults)-1].packetsReceived++
		}

		pinger.OnFinish = func(stats *probing.Statistics) {
			// Add our results to the regionResults array.
			regionResults = append(regionResults, regionResult{
				region:             region,
				areaName:           url["url_friendlyname"],
				packetsTransmitted: regionResults[len(regionResults)-1].packetsTransmitted,
				packetsReceived:    regionResults[len(regionResults)-1].packetsReceived,
				packetLoss:         stats.PacketLoss,
				minPing:            stats.MinRtt.Seconds() * 1000,
				maxPing:            stats.MaxRtt.Seconds() * 1000,
				avgPing:            stats.AvgRtt.Seconds() * 1000,
			})
		}

		// If we are in windows, we need to call pinger.SetPrivileged(true) to use ICMP.
		if runtime.GOOS == "windows" {
			pinger.SetPrivileged(true)
		}

		err = pinger.Run()
		if err != nil {
			panic(err)
		}
	}

	return regionResults
}

func testDynamicIngestPing() dynamicIngestResult {
	// Make a new pinger for ingest.vrcdn.live.
	pinger, err := probing.NewPinger("ingest.vrcdn.live")

	if err != nil {
		panic(err)
	}

	// Listen for Ctrl-C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			pinger.Stop()
			fmt.Printf("\n%s[ping_test:dynamicIngestTest] Stopped.%s\n", Red, Reset)
			os.Exit(0)
		}
	}()

	pinger.Count = 10
	pinger.Interval = 100 * time.Millisecond
	pinger.Timeout = 1 * time.Second

	// Create a empty dynamicIngestResult struct.
	var dynamicIngestResultData dynamicIngestResult

	pinger.OnSend = func(pkt *probing.Packet) {
		// Increment the packetsTransmitted counter.
		dynamicIngestResultData.packetsTransmitted++
	}

	pinger.OnRecv = func(pkt *probing.Packet) {
		// Increment the packetsReceived counter.
		dynamicIngestResultData.packetsReceived++
	}

	pinger.OnFinish = func(stats *probing.Statistics) {
		// Print results.
		fmt.Printf("%s[ping_test:dynamicIngestTest] %d/%d packets received (%.2f%% loss), min/avg/max = %.2f/%.2f/%.2f ms.%s\n", Cyan, stats.PacketsRecv, stats.PacketsSent, stats.PacketLoss, stats.MinRtt.Seconds()*1000, stats.AvgRtt.Seconds()*1000, stats.MaxRtt.Seconds()*1000, Reset)

		// Add our results to the regionResults array.
		dynamicIngestResultData = dynamicIngestResult{
			packetsTransmitted: dynamicIngestResultData.packetsTransmitted,
			packetsReceived:    dynamicIngestResultData.packetsReceived,
			packetLoss:         stats.PacketLoss,
			minPing:            stats.MinRtt.Seconds() * 1000,
			maxPing:            stats.MaxRtt.Seconds() * 1000,
			avgPing:            stats.AvgRtt.Seconds() * 1000,
		}
	}

	// If we are in windows, we need to call pinger.SetPrivileged(true) to use ICMP.
	if runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}

	err = pinger.Run()
	if err != nil {
		panic(err)
	}

	return dynamicIngestResultData
}

func Traceroute(mode string, region string) []string {
	// Figure out the binary we need to use based on the OS.
	var binary string
	if runtime.GOOS == "windows" {
		binary = "tracert"
	} else {
		binary = "traceroute"
	}

	// Ensure that the binary exists.
	if _, err := exec.LookPath(binary); err != nil {
		fmt.Printf("%s[!] Test Failed: %s %s is not installed. Please install it and try again.\n", Red, Reset, binary)
		os.Exit(1)
	}

	switch mode {
	case "region":
		// Get URLs for the region.
		regionUrls := urls[region]

		// Make the results array for the regions, using the regionRequest struct.
		var regionResults []string = make([]string, len(regionUrls))

		// Print region name.
		fmt.Printf("%s[traceroute_test] Testing region %s...%s\n", Yellow, region, Reset)

		// Test each URL.
		for _, url := range regionUrls {
			fmt.Printf("[traceroute_test] Testing endpoint %s in region %s\n", url["url_friendlyname"], region)

			// Run the traceroute command.
			out, err := exec.Command(binary, url["url"]).Output()

			if err != nil {
				fmt.Printf("%s[!] Test Failed: %s %s failed to run. Please install it and try again.\n", Red, Reset, binary)
				os.Exit(1)
			}

			// Add our results to the regionResults array.
			regionResults = append(regionResults, string(out))
		}

		return regionResults
	case "dynamicIngest":
		// Run the traceroute command.
		out, err := exec.Command(binary, "ingest.vrcdn.live").Output()

		if err != nil {
			fmt.Printf("%s[!] Test Failed: %s %s failed to run. Please install it and try again.\n", Red, Reset, binary)
			os.Exit(1)
		}

		return []string{string(out)}
	}

	return []string{}
}

func printResults(regionResults []regionResult, dynamicIngestResultData dynamicIngestResult) {
	// Find the best area in our selected region.
	var bestArea regionResult = regionResult{
		avgPing: 9999999999,
	}

	for _, regionResult := range regionResults {
		if regionResult.areaName == "" {
			continue
		}

		// Check if the regionResult is better than the best area.
		if regionResult.avgPing < bestArea.avgPing {
			bestArea = regionResult
		} else if regionResult.avgPing == bestArea.avgPing {
			// If the ping is the same, check if the packet loss is better.
			if regionResult.packetLoss < bestArea.packetLoss {
				bestArea = regionResult
			}
		} else {
			// If the ping is worse, skip it.
			continue
		}
	}

	// Check if dynamic is better than the best area.
	if dynamicIngestResultData.avgPing < bestArea.avgPing {
		bestArea = regionResult{ // Reconstruct the bestArea struct with the dynamicIngestResultData.
			region:             "dynamicIngest",
			areaName:           "Dynamic Ingest",
			packetsTransmitted: dynamicIngestResultData.packetsTransmitted,
			packetsReceived:    dynamicIngestResultData.packetsReceived,
			packetLoss:         dynamicIngestResultData.packetLoss,
			minPing:            dynamicIngestResultData.minPing,
			maxPing:            dynamicIngestResultData.maxPing,
			avgPing:            dynamicIngestResultData.avgPing,
		}
	}

	// Print the results.
	// ping less then 100ms = green
	// ping between 100ms and 200ms = yellow
	// ping more then 200ms = red

	bestAreaPingColor := Green
	if bestArea.avgPing > 100 && bestArea.avgPing < 200 {
		bestAreaPingColor = Yellow
	} else if bestArea.avgPing > 200 {
		bestAreaPingColor = Red
	}

	fmt.Printf("Your best area in your region is %s%s%s with a average ping of %.2fms.%s\n", Green, bestArea.areaName, bestAreaPingColor, bestArea.avgPing, Reset)

	// Print the results for each area in our region.
	for _, regionResult := range regionResults {
		// ping less then 100ms = green
		// ping between 100ms and 200ms = yellow
		// ping more then 200ms = red
		regionResultPingColor := Green
		if regionResult.avgPing > 100 && regionResult.avgPing < 200 {
			regionResultPingColor = Yellow
		} else if regionResult.avgPing > 200 {
			regionResultPingColor = Red
		}

		// If no area name, continue.
		if regionResult.areaName == "" {
			continue
		}

		fmt.Printf("%s%s%s: %s%d/%d%s packets received (%s%.2f%%%s loss), min/avg/max = %s%.2f%s/%s%.2f%s/%s%.2f%s ms.%s\n", Cyan, regionResult.areaName, Reset, Cyan, regionResult.packetsReceived, regionResult.packetsTransmitted, Reset, Cyan, regionResult.packetLoss, Reset, Cyan, regionResult.minPing, Reset, regionResultPingColor, regionResult.avgPing, Reset, Cyan, regionResult.maxPing, Reset, Reset)
	}

	// Print the results for dynamic ingest.
	// ping less then 100ms = green
	// ping between 100ms and 200ms = yellow
	// ping more then 200ms = red
	dynamicIngestResultPingColor := Green
	if dynamicIngestResultData.avgPing > 100 && dynamicIngestResultData.avgPing < 200 {
		dynamicIngestResultPingColor = Yellow
	} else if dynamicIngestResultData.avgPing > 200 {
		dynamicIngestResultPingColor = Red
	}

	fmt.Printf("%sDynamic Ingest%s: %s%d/%d%s packets received (%s%.2f%%%s loss), min/avg/max = %s%.2f%s/%s%.2f%s/%s%.2f%s ms.%s\n", Cyan, Reset, Cyan, dynamicIngestResultData.packetsReceived, dynamicIngestResultData.packetsTransmitted, Reset, Cyan, dynamicIngestResultData.packetLoss, Reset, Cyan, dynamicIngestResultData.minPing, Reset, dynamicIngestResultPingColor, dynamicIngestResultData.avgPing, Reset, Cyan, dynamicIngestResultData.maxPing, Reset, Reset)

	// Print the traceroute results.
	fmt.Printf("\nTraceroute results have been saved to traceroute_region.txt and traceroute_dynamicIngest.txt\n")

	// Print the disclaimer.
	fmt.Printf("\nDisclaimer: This tool should not be used as a 100%% accurate way to determine your best region or a possible network issue. Please work with official VRCDN support to determine the cause of any issues you may be having.\n\n")

	// Print warning about traceroute.
	fmt.Printf("%s%s%s\n", Yellow, "WARNING: Traceroute results contain SENSITIVE data. Please share with care!", Reset)
}

func main() {
	fmt.Printf("VRCDN Network Test v1.0, built with %s\n", runtime.Version())
	fmt.Print("Created by sticksdev. This tool is not affiliated with VRCDN and is community-made. See LICENSE for more information.\n\n")

	fmt.Print("This tool will test your connection to the VRCDN network. Please select your closest region.\n")

	// Print regions.
	for i, region := range regions {
		fmt.Printf("%d. %s\n", i+1, region)
	}

	// Get user input.
	var region int
	var errorInput error

	// Keep asking for input until a valid region is selected.
	for {
		fmt.Print("Please enter the number of your region: ")
		fmt.Scan(&region, &errorInput)

		// If the input is valid, break out of the loop.
		if region > 0 && region <= len(regions) {
			break
		} else {
			fmt.Print("\nInvalid region selected. Please try again.")

			// Sleep for .5 seconds
			time.Sleep(500 * time.Millisecond)
		}

		if errorInput != nil {
			fmt.Print(Red + "[!]" + Reset + " a error occurred while reading your input. Please try again.\n")
			fmt.Print(errorInput.Error() + "\n")

			// Sleep for .5 seconds
			time.Sleep(500 * time.Millisecond)
		}
	}

	// Test the region.
	var regionResults []regionResult = testRegionPing(regions[region-1])

	if len(regionResults) == 0 {
		fmt.Printf("%s[!]%s No results were returned from the previous test. Please try again.\n", Red, Reset)
		os.Exit(1)
	}

	fmt.Printf("[ping_test] Running dynamic ingest test... (ingest.vrcdn.live)\n")
	var dynamicIngestResultData dynamicIngestResult = testDynamicIngestPing()

	if dynamicIngestResultData == (dynamicIngestResult{}) {
		fmt.Printf("%s[!]%s No results were returned from the previous test. Please try again.\n", Red, Reset)
		os.Exit(1)
	}

	var traceRegionResults []string = Traceroute("region", regions[region-1])

	if len(traceRegionResults) == 0 {
		fmt.Printf("%s[!]%s No results were returned from the previous test. Please try again.\n", Red, Reset)
		os.Exit(1)
	}

	fmt.Printf("[traceroute_test] Running dynamic ingest test... (ingest.vrcdn.live)\n")
	var traceDynamicIngestResults []string = Traceroute("dynamicIngest", "")

	if len(traceDynamicIngestResults) == 0 {
		fmt.Printf("%s[!]%s No results were returned from the previous test. Please try again.\n", Red, Reset)
		os.Exit(1)
	}

	fmt.Printf("%s[...]%s Saving traceroute results to traceroute_region.txt and traceroute_dynamicIngest.txt\n", Yellow, Reset)

	// Save the traceroute results to a file.
	var traceRegionFile, traceDynamicIngestFile *os.File

	// If we have old results, delete them.
	if _, err := os.Stat("traceroute_region.txt"); err == nil {
		os.Remove("traceroute_region.txt")
	}

	if _, err := os.Stat("traceroute_dynamicIngest.txt"); err == nil {
		os.Remove("traceroute_dynamicIngest.txt")
	}

	// Define error.
	var err error

	traceRegionFile, err = os.OpenFile("traceroute_region.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("%s[!]%s Failed to open traceroute_region.txt. Please try again.\n", Red, Reset)
		os.Exit(1)
	}

	traceDynamicIngestFile, err = os.OpenFile("traceroute_dynamicIngest.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("%s[!]%s Failed to open traceroute_dynamicIngest.txt. Please try again.\n", Red, Reset)
		os.Exit(1)
	}

	// Write the results to the files.
	for _, traceRegionResult := range traceRegionResults {
		traceRegionFile.WriteString(traceRegionResult)
	}

	for _, traceDynamicIngestResult := range traceDynamicIngestResults {
		traceDynamicIngestFile.WriteString(traceDynamicIngestResult)
	}

	// Close the files.
	traceRegionFile.Close()
	traceDynamicIngestFile.Close()

	fmt.Printf("%s[i]%s Save complete. Showing results!\n", Green, Reset)
	printResults(regionResults, dynamicIngestResultData)

	// Sleep for 5 seconds.
	time.Sleep(5 * time.Second)

	// Exit.
	os.Exit(0)
}
