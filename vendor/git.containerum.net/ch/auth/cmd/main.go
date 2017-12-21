package main

import (
	"flag"
	"fmt"
	"os"

	"git.containerum.net/ch/auth/utils"
	"git.containerum.net/ch/grpc-proto-files/auth"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var serverAddress string

var (
	userAgent   string
	fingerprint string
	userIp      string
)

func init() {
	flag.StringVar(&serverAddress, "s", "127.0.0.1:8888", "Auth server address")
	flag.StringVar(&userAgent, "ua", "Mozilla/5.0 (X11; Linux x86_64; rv:56.0) Gecko/20100101 Firefox/56.0",
		"User agent")
	flag.StringVar(&fingerprint, "ufp", "101924019824", "User fingerprint")
	flag.StringVar(&userIp, "uip", "127.0.0.1", "User IP")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s [options] <cmd> [cmdparams]
cmd:
	issue
		Issues access and refresh tokens
	check
		Cshecks access token. Requires access token as parameter (cmdparams)
options:
`, os.Args[0])
		flag.PrintDefaults()
	}
}

func chkErr(err error) {
	if err != nil {
		log.Errorf("Error: %v", err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(), // disable transport security
	}

	log.Infof("Setup connection to %v", serverAddress)
	conn, err := grpc.Dial(serverAddress, opts...)
	chkErr(err)
	client := auth.NewAuthClient(conn)
	switch flag.Arg(0) {
	case "issue":
		log.Infoln("Issuing token")
		resp, err := client.CreateToken(context.Background(), &auth.CreateTokenRequest{
			UserAgent:   userAgent,
			Fingerprint: fingerprint,
			UserId:      utils.NewUUID(),
			UserRole:    auth.Role_USER,
			RwAccess:    true,
			Access:      &auth.ResourcesAccess{},
			PartTokenId: nil,
		})
		chkErr(err)
		log.Printf("Got response %+v", resp)
	case "check":
		accessToken := flag.Arg(1)
		log.Infoln("Checking access token %v", accessToken)
		resp, err := client.CheckToken(context.Background(), &auth.CheckTokenRequest{
			AccessToken: accessToken,
			UserAgent:   userAgent,
			FingerPrint: fingerprint,
			UserIp:      userIp,
		})
		chkErr(err)
		log.Printf("Got response %+v", resp)
	default:
		log.Errorln("Invalid command specified")
		flag.Usage()
		os.Exit(1)
	}
}
