Cloudflare Dynamic DNS updater tool.
====================================
This tool updates a specific record in your cloudflare DNS zone with the local IP of the machine it's running on (Not the public IP).


Set environment variables: 
* `CF_API_TOKEN` -- your API token
* `CF_ZONE_NAME` -- the zone name `example.org` or whatever
* `CF_RECORD_PREFIX` -- the name for the record `dynamic.example.org` for example. 
* `CF_PROXY_MODE` -- "true" or "false" as to whether to enable CF proxy mode and all that jazz. 

run the application on a cronjob or something.

