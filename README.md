# Resy-Snipy (Lite)

If you're like me and need a reservation at Carbone NY and prefer to be notified via Discord (webhook), then this might be for you. Or can be used as a template to scrape other sites and be notified via Discord.

Program reads the list of dates/restaurants in your restys.txt file and will post when an opening populates, checks every minute.  
## Usage
Go 1.1+  

```
go run main.go --wh=WEBHOOK_URL

# go mod download (to install dependencies, if needed)
```

### Different Restaurant?
Open the `resys.txt` and add the desired date, party size, and venueId in the following format:
`2022-03-01,2,6194`.  
Feel free to use [this](https://resy-api.vercel.app/api/v1/resy?location=ny&slug=carbone) 
by changing the location and slug params to get the venueId.  
Example: `https://resy.com/cities/la/cassia` => `https://resy-api.vercel.app/api/v1/resy?location=la&slug=cassia`

### Questions/Concerns?
Direct them to @Krev#0001 on Discord.  