package main

import (
  "crypto/tls"
  "flag"
  "fmt"
  _ "io/ioutil"
  "log"
  "net"
  "net/http"
  "net/mail"
  "net/smtp" // See how to do that here: https://godoc.org/net/smtp
  "net/url"
  "strings"
  "time"

  "github.com/vharitonsky/iniflags"
)

var (
  from = flag.String("from", "", "From-Email address")
  to = flag.String("to", "", "To-Email address")
  pw = flag.String("pw", "", "From-Email account password. If auth is CRAMMD5 pw is used as secret.")
  user = flag.String("user", "", "From-Email account username, if not provided From-Email address will be used instead")
  endpoint = flag.String("url", "/contact", "URL on which to receive contact form data")
  addr = flag.String("addr", "127.0.0.1:443", "Listen address of HTTP server as IP:PORT")
  mail_srv = flag.String("smtp", "", "Address of the mail server for From-Email as IP:PORT")
  certFile = flag.String("cert", "./cert.crt", "SSL/TLS certificate (x509 pem)")
  keyFile = flag.String("key", "./cert.key", "Private key")
  mail_auth = flag.String("auth", "PLAIN", "SMTP authentication method: CRAMMD5, PLAIN")
  mail_tls = flag.Bool("starttls", true, "Use STARTTLS for mail transport")
  useTLS = flag.Bool("tls", true, "Use HTTPS for included HTTP-server")
)

func init() {
  iniflags.Parse()

  if len(*from) == 0 {
    log.Fatal("No From-Email address supplied")
  }
  if len(*to) == 0 {
    log.Fatal("No To-Email address supplied")
  }
  if len(*pw) == 0 {
    log.Println("Warning no password supplied")
  }
  if len(*user) == 0 {
    log.Printf("No username supplied, will use %s instead", *from)
    *user = *from
  }
  if len(*mail_srv) == 0 {
    log.Fatal("No mail server address supplied")
  }
  if strings.Compare(*mail_auth, "CRAMMD5") != 0 && strings.Compare(*mail_auth, "PLAIN")  != 0 {
    log.Fatal("No valid authentication method for SMTP supplied")
  }
}

// sendMail takes the parsed contact form and sends it as an email.
// For this to work with gmail one needs to enable less secure apps:
// https://www.google.com/settings/security/lesssecureapps
func sendMail(form url.Values) {
  host, port, _ := net.SplitHostPort(*mail_srv)
  var a smtp.Auth
  if strings.Compare(*mail_auth, "CRAMMD5") == 0 {
    a = smtp.CRAMMD5Auth(*user, *pw)
  } else {
    a = smtp.PlainAuth("", *user, *pw, host)
  }

  f := mail.Address{"", *from}
  subj := "Contact"
  body := ""
  for k, v := range form {
    body += fmt.Sprintf("%s: %s\r\n", k, v[0])
  }

  headers := make(map[string]string)
//  headers["From"] = f.String()
//  headers["To"] = t.String()
  headers["Subject"] = subj

  message := ""
  for k, v := range headers {
    message += fmt.Sprintf("%s: %s\r\n", k, v)
  }
  message += "\r\n" + body

  var conn net.Conn
  var err error
  tlsConfig := &tls.Config{
//      InsecureSkipVerify: true,
    ServerName: host,
  }
  if strings.Compare(port, "465") == 0 {
    conn, err = tls.Dial("tcp", *mail_srv, tlsConfig)
    if err != nil {
      log.Printf("Failed to connect to mail server: %s", err)
      return
    }
    log.Printf("Use TLS connection")
  } else {
    conn, err = net.Dial("tcp", *mail_srv)
    if err != nil {
      log.Printf("Failed to connect to mail server: %s", err)
      return
    }
    log.Printf("Use unecrypted connection")
  }
  defer conn.Close()

  c, err := smtp.NewClient(conn, host)
  if err != nil {
    log.Printf("Failed to create mail client: %s", err)
    return
  }

  if err := c.Hello("localhost"); err != nil {
    log.Printf("Failed to send HELO or EHLO: %s", err)
    return
  }

  if *mail_tls {
    if ok, _ := c.Extension("STARTTLS"); ok {
      if err = c.StartTLS(tlsConfig); err != nil {
	log.Printf("Failed to initiate STARTTLS: %s", err)
	return
      }
    }
  }

  plain, _ := c.Extension("PLAIN")
  auth, _ := c.Extension("AUTH")
  if plain || auth {
    if err = c.Auth(a); err != nil {
      log.Printf("Failed to authenticate: %s", err)
      return
    }
  }

  if err = c.Mail(f.Address); err != nil {
    log.Printf("Failed to set from-address: %s", err)
    return
  }

  toAddresses := strings.Split(*to, ",")
  for _, t := range toAddresses {
    if err = c.Rcpt(t); err != nil {
      log.Printf("Failed to set to-address: %s, %s", t, err)
      return
    }
  }

  w, err := c.Data()
  if err != nil {
    log.Printf("Failed to send data command: %s", err)
    return
  }
  _, err = w.Write([]byte(message))
  if err != nil {
    log.Printf("Failed to send mail: %s", err)
    return
  }
  w.Close()

  c.Quit()

  log.Printf("Sent mail to: %s", *to)
}

// handleContactForm takes the request containing the contact form, parses
// it and sends an email containing the contact form data.
// TODO: Limit requests from same host. How to do that?
func handleContactForm(w http.ResponseWriter, r *http.Request) {
  err := r.ParseForm()
  if err != nil {
    log.Printf("Failed to parse form: %s", err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  log.Printf("Received contact form: %v", r.PostForm)

  go sendMail(r.PostForm)
  w.WriteHeader(http.StatusAccepted)
}

func main() {
  // Server configuration
  s := &http.Server{
    Addr: *addr,
    ReadTimeout: time.Second * 10,
    WriteTimeout: time.Second * 10,
  }

  // URL handler for the contact form.
  log.Printf("Use URL %s to receive contact form data", *endpoint)
  http.HandleFunc(*endpoint, handleContactForm)

  // Run the HTTP server.
  log.Printf("Start listening for contact forms on %s", *addr)
  if !*useTLS {
    log.Fatal(s.ListenAndServe())
  } else {
    log.Fatal(s.ListenAndServeTLS(*certFile, *keyFile))
  }
}
