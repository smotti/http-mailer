package main

import (
  "crypto/tls"
  "log"
  "net/http"
  _ "smtp" // See how to do that here: https://godoc.org/net/smtp
  "time"

  "github.com/vharitonsky/iniflags"
)

var (
  from = flag.String("from", "", "From-Email address")
  to = flag.String("to", "", "To-Email address")
  pw = flag.String("pw", "", "From-Email account password")
  user = flag.String("user", "", "From-Email account username, if not provided From-Email address will be used instead")
  url = flag.String("url", "/contact", "URL on which to receive contact form data")
  addr = flag.String("addr", "127.0.0.1:443", "Listen address of HTTP server as IP:PORT")
  mail_srv = flag.String("smtp", "", "Address of the mail server for From-Email")
  certFile = flag.String("cert", "./cert.crt", "SSL/TLS certificate (x509 pem)")
  keyFile = flag.String("key", "./cert.key", "Private key")
)

func init() {
  iniflags.Parse()

  if len(from) == 0 {
    log.Fatal("No From-Email address supplied")
  }
  if len(to) == 0 {
    log.Fatal("No To-Email address supplied")
  }
  if len(pw) == 0 {
    log.Println("Warning no password supplied")
  }
  if len(user) == 0 {
    log.Printf("No username supplied, will use %s instead", *from)
  }
  if len(mail_srv) == 0 {
    log.Fatal("No mail server address supplied")
  }
}

func parseContactForm(form string) (map[string]string, error) {
  return map[string]string{}, nil
}

func sendMail(mail map[string]string) error {
  return nil
}

func handleContactForm(w http.ResponseWrite, r *http.Request) {
  log.Println("Send mail")
}

func main() {
  s := &http.Server{
    Addr: *addr,
    ReadTimeout: time.Second * 10,
    WriteTimeout: time.Second * 10,
  }

  http.Handle(*url, handleContactForm)
  log.Fatal(s.ListenAndServeTLS(*certFile, *keyFile))
}
