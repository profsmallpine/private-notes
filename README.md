# private-notes

If this application gets off the ground, it will support groups having long form private
conversations using E2E encryption with keys stored by the consuming clients.

## Roadmap + TODOs

- [x] Migrate all unnecessary util code to github.com/xy-planning-network/trails
- [x] Move http to root directory + tmpl into http
- [x] Implement better redirect handling with flash messages
- [x] Add notifications for creating a note + creating a comment (plaintext email using basic auth)
- [x] Fix bug with trix not rendering content
- [ ] Setup Procfile + run on heroku
- [ ] Setup DNS on domain + add users to database for login
- [ ] Update to latest version of tailwindcss
- [ ] Add environment + embed in necessary structs
- [ ] Fix flash once trails adds update
- [ ] Render flash with alpinejs
- [ ] Update header to context based breadcrumb
- [ ] Implement application encryption for note content (E2E)
- [ ] Add draft feature for notes
- [ ] Add edit ability to group
- [ ] Add edit ability to note
- [ ] Add testing
- [ ] Add sentry notifications for bugs
- [ ] Add better sigterm handling
