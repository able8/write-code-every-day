package main

import "time"

type Post struct {
	Date    time.Time
	Title   string
	Content string
}

var posts = []Post{
	{
		Date:  time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		Title: "Happy New Year!",
		Content: `New Year is a widely celebrated occasion in the United Kingdom, marking the end of one year and the beginning of another.

Top New Year Activities in the UK include:

* Attending a Hogmanay celebration in Scotland
* Taking part in a local First-Foot tradition in Scotland and Northern England
* Setting personal resolutions and goals for the upcoming year
* Going for a New Year's Day walk to enjoy the fresh start
* Visiting a local pub for a celebratory toast and some cheer
`,
	},
	{
		Date:  time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC),
		Title: "May Day",
		Content: `May Day is an ancient spring festival celebrated on the first of May in the United Kingdom, embracing the arrival of warmer weather and the renewal of life.

Top May Day Activities in the UK:

* Dancing around the Maypole, a traditional folk activity
* Attending local village fetes and fairs
* Watching or participating in Morris dancing performances
* Enjoying the public holiday known as Early May Bank Holiday
`,
	},
}
