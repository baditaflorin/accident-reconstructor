# Privacy

The GitHub Pages frontend collects no analytics in v1.

User-selected videos are not uploaded until the user presses Run. When Run is pressed, videos are sent to the configured backend URL shown in the UI.

Browser storage:

- Last successful reconstruction artifact is cached in IndexedDB.
- API base URL is stored in localStorage.

Server storage:

- Uploaded videos and generated artifacts are written to the backend case volume.
- Operators are responsible for retention, deletion, and backups.

Support link:

https://www.paypal.com/paypalme/florinbadita
