# aproxygo

Aproxygo is a minimal reverse proxy written in go.

It supports ssl out of the box with automatic certificate renewal (thanks to acme/autocert!).

It comes with no configuration and is meant as inspiration and has not been thoroughly tested in production. Use at your own risk.

# Building
`go get github.com/Sketchground/aproxygo`

# Deploying
1) Copy compiled binary (assuming server and computer you built binary on is using the same processor architecture).
2) Edit service file to your needs and install it on your server
3) run `sudo systemctl start aproxygo` (possibly `sudo systemctl enable aproxygo`)

# Purpose
If you have a bunch of go project on the same host and don't want to bother with apache/nginx or the like, maybe this is useful.

If you have more complex configuration needs nothing stops you from extending it to be configured over etcd or the like.

It scores an A rating on ssllabs.com out of the box https://www.ssllabs.com/ssltest/analyze.html?d=blog.sketchground.dk&hideResults=on
