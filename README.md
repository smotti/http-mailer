# Description

This is a simple go app that comes with it's own HTTP server that accepts
POST requests and sends it of to the provided from e-mail address to the
given to e-mail address.

# Usage

```
# ./http-mailer --help
Usage of ./http-mailer:
  -addr string
          Listen address of HTTP server as IP:PORT (default "127.0.0.1:443")
  -auth string
	  SMTP authentication method: CRAMMD5, PLAIN (default "PLAIN")
  -cert string
	  SSL/TLS certificate (x509 pem) (default "./cert.crt")
  -config string
	  Path to ini config for using in go flags. May be relative to the current executable path.
  -configUpdateInterval duration
	  Update interval for re-reading config file set via -config flag. Zero disables config file re-reading.
  -dumpflags
	  Dumps values for all flags defined in the app into stdout in ini-compatible syntax and terminates the app.
  -from string
	  From-Email address
  -key string
	  Private key (default "./cert.key")
  -pw string
	  From-Email account password. If auth is CRAMMD5 pw is used as secret.
  -smtp string
	  Address of the mail server for From-Email as IP:PORT
  -starttls
	  Use STARTTLS for mail transport (default true)
  -tls
	  Use HTTPS for included HTTP-server (default true)
  -to string
	  To-Email address
  -url string
	  URL on which to receive contact form data (default "/contact")
  -user string
	  From-Email account username, if not provided From-Email address will be used instead
```

*Note that by default the HTTP Server uses HTTPS, thus you need to either change
the HTTP server listen address or provide a cert and the according private key.*

# Example

contact-form.html:

```
<html>
  <body>
    <form action="https://127.0.0.1:8888/contact" method="post">
      First name:<br>
      <input type="text" name="firstname"><br>
      Last name:<br>
      <input type="text" name="lastname"><br>
      E-Mail:<br>
      <input type="email" name="email"><br>
      Message:<br>
      <textarea name="message" rows="4" cols="50"></textarea><br>
      <input type="submit" value="Contact">
    </form>
  </body>
</html>
```
