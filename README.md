# private-notes

If this application gets off the ground, it will support groups having long form private
conversations using E2E encryption with keys stored by the consuming clients.

## Development instructions

1. Install Go 1.17
2. Install PostgreSQL 13.1
3. `./scripts/setup_dev_dbs` to create development database
4. Run `cp .env.example .env`
5. Update .env for login (needs google creds + session keys)
6. If you want emails to work, you can create an [app password](https://myaccount.google.com/apppasswords) and use that in the .env.
7. Run `go run main.go`

## Generate Session Keys

- use gore https://github.com/motemen/gore (OR) golang playground https://play.golang.org to generate keys
- use golang crypto/rand.Read https://golang.org/pkg/crypto/rand/
- (OR) securecookie.GenerateRandomKey() https://godoc.org/github.com/gorilla/securecookie
- `fmt.Printf("SESSION_AUTH_KEY=%x \n", securecookie.GenerateRandomKey(64)) fmt.Printf("SESSION_ENCRYPTION_KEY=%x \n", securecookie.GenerateRandomKey(32))`

## Roadmap + TODOs

- [x] Migrate all unnecessary util code to github.com/xy-planning-network/trails
- [x] Move http to root directory + tmpl into http
- [x] Implement better redirect handling with flash messages
- [x] Add notifications for creating a note + creating a comment (plaintext email using basic auth)
- [x] Fix bug with trix not rendering content
- [x] Setup Procfile + run on heroku
- [x] Setup DNS on domain + add users to database for login
- [x] Update to latest version of tailwindcss
- [x] Add environment + embed in necessary structs
- [x] Add logging and update email failure (create note + create comment)
- [x] Add guide for running the app for dev
- [x] Render flash with alpinejs
- [x] Paginate lists (notes)
- [x] Add MVP of meetings
- [x] Fix wonky spacing in comments (either newlines aren't captured or not being rendered)
- [x] Add better sigterm handling
- [ ] Paginate meetings
- [ ] Add MVP of meeting review
- [ ] Update header to context based breadcrumb
- [ ] Implement application encryption for note content (E2E)
- [ ] Add draft feature for notes
- [ ] Add edit ability to group
- [ ] Add testing
- [ ] Improve rendering of partials (remove need for wrapper) and move to trails
