# VRCDN Network Test

This is a simple tool to help you test your connection to the VRCDNs network. This can be useful to determine if you are having issues connecting to the VRCDN network.

## How does it work?

This tool will run two tests to test your connectivity to the VRCDN network. The first is a ping test, which will ping
each of the VRCDN ingest nodes for your area and report the results. The second is traceroute, which will show you the
path your connection takes to reach the VRCDN network.

Trace route is a useful tool to determine if you are having issues connecting to the VRCDN network. If you see a lot of
packet loss or high latency in the trace route, this could indicate an issue with your connection to the VRCDN network.

## How do I use it?

You can download the latest release from the [releases page]() - or you can build it yourself.

Once you have the binary, you can simply double click it to run it. Results will be shown for 5 seconds, and then the
program will exit. We recommend opening a command prompt and running the program from there so you can see the results
after the program exits, or use one of our pre-built scripts to run the program for you.

A sample output from the program is shown below:

![Sample output](https://img.sticks.ovh/9Loapvy8u.png)

## Compiling from source

This project uses [Go](https://golang.org/) to compile. You can download Go from [here](https://golang.org/dl/). Once
you have Go installed, you can compile the project by running `go build` in the root directory of the project. This will
create a binary called `vrcdn-nettest` in the root directory of the project that is ready to run. Alternatively, you can
run `go run .` to run the program without compiling it first.
