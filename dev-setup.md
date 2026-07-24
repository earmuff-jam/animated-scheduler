# How to get Facebook token

Verify existing token debugger - https://developers.facebook.com/tools/debug/accesstoken/?utm_source=chatgpt.com

Step 1: Get a long-lived user token (60 days)

Open this URL in your browser, replacing the placeholders:

https://graph.facebook.com/v25.0/oauth/access_token?grant_type=fb_exchange_token&client_id=1052274500625787&client_secret=YOUR_APP_SECRET&fb_exchange_token=YOUR_CURRENT_TOKEN

You need:

client_id: 1052274500625787
client_secret: from Meta App Dashboard → Settings → Basic
fb_exchange_token: the current token from Graph API Explorer

Step 2: Exchange that for a Page token

Now take the new 60-day user token and call:
https://graph.facebook.com/v25.0/me/accounts?access_token=LONG_LIVED_USER_TOKEN

Copy the access_token for Homehive Solutions.

That is the token you should store in Netlify.

Step 3: Verify it

Go back to the Access Token Debugger and paste the new Page token.
