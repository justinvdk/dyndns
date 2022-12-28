package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/libdns/leaseweb"
	"github.com/libdns/libdns"
)

var (
	username       string
	password       string
	dyndnsZone     string
	leasewebApiKey string
	port           string
)

func init() {
	flag.StringVar(&username, "username", getEnv("DYNDNS_USERNAME", ""), "give me a username")
	flag.StringVar(&password, "password", getEnv("DYNDNS_PASSWORD", ""), "give me a password")
	flag.StringVar(&dyndnsZone, "dyndnsZone", getEnv("DYNDNS_ZONE", ""), "give me a dyndnsZone")
	flag.StringVar(&leasewebApiKey, "leasewebApiKey", getEnv("DYNDNS_LEASEWEB_API_KEY", ""), "give me a leasewebApiKey")
	flag.StringVar(&port, "port", getEnv("DYNDNS_PORT_NUMBER", "80"), "give me a port number")
}

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.Handle("/", dyndnsHandler(leasewebApiKey, username, password))

	log.Printf("Starting up on port %s", port)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func dyndnsHandler(leasewebApiKey string, username string, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		u, p, ok := req.BasicAuth()
		if !ok {
			fmt.Println("Error parsing basic auth.")
			w.WriteHeader(401)
			return
		}
		if u != username {
			fmt.Printf("Username provided is incorrect: %s.\n", u)
			w.WriteHeader(401)
			return
		}
		if p != password {
			fmt.Printf("Password provided is incorrect: %s.\n", p)
			w.WriteHeader(401)
			return
		}
		fmt.Printf("Username and password provided is correct.\n")

		parsedUrl, _ := url.Parse(req.URL.String())
		hostname := parsedUrl.Query().Get("hostname")
		if len(hostname) == 0 {
			fmt.Printf("Must provide hostname query parameter.\n")
			w.WriteHeader(400)
			return
		}

		realIP := getRealIP(req)

		fmt.Printf("Hostname: %s.\n", hostname)
		fmt.Printf("RemoteAddr: %s.\n", req.RemoteAddr)
		fmt.Printf("RealIp: %s.\n", realIP)

		fmt.Printf("Will use %s as apiKey.\n", leasewebApiKey)
		provider := leaseweb.Provider{APIKey: leasewebApiKey}

		records := []libdns.Record{
			libdns.Record{
				Type:  "A",
				Name:  hostname,
				Value: realIP,
				TTL:   300 * time.Second,
			},
		}

		ctx := context.TODO()

		// TODO Do not do this without retrieving existing records.
		// its a literal SET, so overrides existing.
		records, err := provider.SetRecords(ctx, dyndnsZone, records)
		if err != nil {
			fmt.Printf("Something went wrong during contacting leaseweb.\n")
			fmt.Println(err.Error())
			w.WriteHeader(500)
			return

		}

		for _, record := range records {
			fmt.Printf("%s %v %s %s\n", record.Name, record.TTL.Seconds(), record.Type, record.Value)
		}

		// <remote_IP_address> - [<timestamp>] "<request_method> <request_path> <request_protocol>" -
		// log.Printf("%s - - [%s] \"%s %s %s\" - -", r.RemoteAddr, time.Now().Format("02/Jan/2006:15:04:05 -0700"), r.Method, r.URL.Path, r.Proto)
	})
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getRealIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-IP")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarder-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
