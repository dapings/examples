package extensions

import (
	"fmt"
	"math/rand"
	"strings"
)

var (
	uaGens = []func() string{
		genFirefoxUA,
		genChromeUA,
		genEdgeUA,
		genOperaUA,
	}

	// nolint
	uaGensMobile = []func() string{
		genMobileUcwebUA,
		genMobileNexus10UA,
	}
)

func GenerateRandomUA() string {
	return uaGens[rand.Intn(len(uaGens))]()
}

// Generates Firefox Browser User-Agent (Desktop)
// e.g. "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:87.0) Gecko/20100101 Firefox/87.0"
func genFirefoxUA() string {
	version := ffVersions[rand.Intn(len(ffVersions))]
	os := osStrings[rand.Intn(len(osStrings))]

	return fmt.Sprintf("Mozilla/5.0 (%s; rv:%.1f) Gecko/20100101 Firefox/%.1f", os, version, version)
}

// Generates Chrome Browser User-Agent (Desktop)
// e.g. "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36"
func genChromeUA() string {
	version := chromeVersions[rand.Intn(len(chromeVersions))]
	os := osStrings[rand.Intn(len(osStrings))]

	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", os, version)
}

// Generates Microsoft Edge User-Agent (Desktop)
// e.g. "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36 Edg/90.0.818.39"
func genEdgeUA() string {
	version := edgeVersions[rand.Intn(len(edgeVersions))]
	splitVers := strings.Split(version, ",")
	chromeVersion := splitVers[0]
	edgeVersion := splitVers[1]
	os := osStrings[rand.Intn(len(osStrings))]

	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36 Edg/%s",
		os, chromeVersion, edgeVersion)
}

// Generates Opera Browser User-Agent (Desktop)
// e.g. "Opera/9.80 (X11; Linux x86_64; U; en) Presto/2.8.131 Version/11.11"
func genOperaUA() string {
	version := operaVersions[rand.Intn(len(operaVersions))]
	os := osStrings[rand.Intn(len(osStrings))]

	return fmt.Sprintf("Opera/9.80 (%s; U; en) Presto/%s", os, version)
}

// Generates UCWEB/Nokia203 Browser User-Agent (Mobile)
// e.g. "UCWEB/2.0 (Java; U; MIDP-2.0; Nokia203/20.37) U2/1.0.0 UCMini/10.9.8.1006 (SpeedMode; Proxy; Android 4.4.4; SM-J110H ) U2/1.0.0 Mobile"
func genMobileUcwebUA() string {
	device := ucwebDevices[rand.Intn(len(ucwebDevices))]
	version := ucwebVersions[rand.Intn(len(ucwebVersions))]
	android := androidVersions[rand.Intn(len(androidVersions))]

	return fmt.Sprintf("UCWEB/2.0 (Java; U; MIDP-2.0; Nokia203/20.37) U2/1.0.0 UCMini/%s (SpeedMode; Proxy; Android %s; %s ) U2/1.0.0 Mobile",
		version, android, device)
}

// Generates Nexus 10 Browser User-Agent (Mobile)
// e.g. "Mozilla/5.0 (Linux; Android 5.1.1; Nexus 10 Build/LMY48T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.91 Safari/537.36"
func genMobileNexus10UA() string {
	build := nexus10Builds[rand.Intn(len(nexus10Builds))]
	android := androidVersions[rand.Intn(len(androidVersions))]
	chrome := chromeVersions[rand.Intn(len(chromeVersions))]
	safari := nexus10Safari[rand.Intn(len(nexus10Safari))]

	return fmt.Sprintf("Mozilla/5.0 (Linux; Android %s; Nexus 10 Build/%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/%s", android, build, chrome, safari)
}
