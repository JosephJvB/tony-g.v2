# Tony G

rebuild Tony with Go

`sam build && sam deploy --profile joe`

Gonna end up deploying him to AWS as per. So I'll use this repo as ref https://github.com/JosephJvB/spotify-users-backend

need:

- aws parameter store api
- youtube api
- google sheets api
- spotify api

steps:

- load secrets from parameter store
  - currently other clients are reading from os.Getenv()
  - I could have param store do os.Setenv() if I wanted
  - but I think I'd rather pass from paramClient into other clients on creation
- load youtube items
- load google sheets
- load existing playlists
  - don't need their items immediately
  - only need to know their items if we are adding to those playlists

https://edu.anarcho-copy.org/Programming%20Languages/Go/Concurrency%20in%20Go.pdf

go clean -testcache: expires all test results

hey actually. Tony and his mates keep an up to date Apple playlist no?
https://music.apple.com/us/playlist/my-fav-singles-of-2025/pl.u-ayeZTygbKDy

Why not use that as the source, rather than youtube video descriptions?

Better flow:
playlist description says it gets updated every friday
so every saturday
get tony's playlists from apple music
find the current one by title (current year)
get all the playlists songs from apple music
(filter for just the recent ones?)

get my spotify playlist for the same year (create if not exists)
get all songs currently in my playlist

find the tracks that need to be added (song_artist_album)
find those tracks in spotify
add them to playlist
Seems easy enough?

nah wait apple music api is garbage

lets do it this way:
scrape https://theneedledrop.com/loved-list/${year}

- years >= 2022 use apple playlists
- years < 2022 use spotify or nothing at all
- with this one, I'm gonna just handle future playlists
  - maybe recreate playlists from 2022 onwards too with new name and decommish the old service

find the apple music link in html:
https://embed.music.apple.com/us/playlist/my-fav-singles-of-2024/pl.u-e2ZmtK9VM5K?wmode=opaque
https://music.apple.com/us/playlist/my-fav-singles-of-2024/pl.u-e2ZmtK9VM5K

scrape the apple music playlist page for tracks: songname, artist, album

get google sheets tracks that I've already tried to add

find those tracks in spotify that I haven't tried to add

get my spotify playlists and their items

search in spotify for those tracks

rather than scraping the playlist url every week should I save it in google sheets?
We'll see

deleted tony2 stack cos it would have found the new Tony #2 playlists due to prefix overlap
I can change the prefix maybe but should be fine now that the old fn is not running

just noticed too actually that the curated Apple Playlists are a lot shorter than the ones I made
see: https://music.apple.com/us/playlist/my-fav-singles-of-2024/pl.u-e2ZmtK9VM5K
vs: https://open.spotify.com/playlist/3cIeEpjP3PhNiFD6aKfyD6?si=a277b1e20a8043f7
ie: Apple Playlists don't have every liked song from Youtube Best Tracks
eep!
Maybe I should keep my guy running!

1. turn tony2 back online (delete web assets stack tho that's garbage)
2. tony-g needs a new prefix "2onys 2op 2racks 2024" ?

- also I should rename tony-g to tony-g that's way better
- then I can have Tony-g2 if I wanna redo the JS service

I wanna split the service into two lambdas.
Then I could have a generic api that can turn Apple Playlist into Spotify Playlist.
But My current service handles checking for existing songs and updating an existing playlist. So maybe that's a separate thing entirely.
So let's not yet.

### todo:

- [x] deploy go lambda to run on chron. can accept payload { "year": int }. Create now that's what I call melon music playlists
- [] continue tony-g2 to replace tony2 service scraping from youtube. Should I do this one? Or just resurrect old service? Nah let's keep pushing I reckon. that original service is ah, not good!
  - do I keep using same google sheets data tho? Yeah I guess even tho I don't like it.
  - I wish I had kept a similar sheet where all tracks get added in, but when I first did it I only added missing tracks.
  - I do think tracking parsed videos is OK tho.
- [] make api gateway lambda that can turn apple playlist into spotify playlist.

I wanna re-parse all youtube descriptions
so I can include all the tracks that were found as well as the missing ones
However, old tony-2 had loads of those "replacers" to fix typos
Is it LLM time?

Google gemini is pretty good at getting the Best Tracks from description
But it's really bad at finding Spotify URLs

- maybe I can use Google Search for that?
- Or as long as Gemini can fix the typos I can continue to use Spotify Search

It seems to handle typos
Can it handle two songs with slash? it can! it can VERY well

so, the steps:

1. Gemini creates a BestTrackList from youtube description
2. Try find all those songs in Spotify
3. If I can't find a song I have the following options
   1. google search api for spotify link (should be more forgiving than spotify search api): yeah this works well
   2. Gemini to do a google search per missing song (seems to work for individual songs? Or it did on chatgpt once, but not since)
   3. ~~scrape google search results~~ doesnt work

I dont think husky precommit is working on this machine :(

Or maybe this:

1. gemini gets tracklist from description(s)
2. Search in spotify for each
3. Add found to Playlist
4. Add all rows to spreadsheet (found and not)

Separate lambda on cron:

1. called once per day
2. lookup all google sheet rows
3. Find all not found rows
4. limit to 100 - daily customsearch rate limit
5. search for missing tracks
6. add found to playlist
7. update all rows

columns might need update

- found by spotify search
- found by google search

  - so that way I don't keep retrying songs that weren't found first time

  But that's like quite a complicated thing
  So maybe it's better to just try to find the missing ones at the time
  There should never really be a time when I need to look up more than 100 in a day right so it's kinda dumb.

how to I wanna launch this

1. Apple Playlist:
   . delete current playlists? Or at least rename them so they don't get in the way
   . test at least one year locally first
   . Deploy as lambda
   . Use Lambda to go thru backlog - using test input year

1. Youtube Playlists
   . using a new playlist prefix so existing lists won't be affected
   . Test a few videos locally first
   . Use Lambda to go thru backlog - using test input video ids

Could everything be in package main?? Have i been overcomplicating the package names
I strongly suspect ya yes yup

check i'm not going over my allowed search quota
https://console.cloud.google.com/apis/api/customsearch.googleapis.com/metrics?authuser=1&inv=1&invt=AbxFsQ&project=tnd-best-tracks&pageState=(%22duration%22:(%22groupValue%22:%22P1D%22,%22customValue%22:null))

run without google search on loads/all of the videos?

Review results:

- spot check a few videos to make sure Gemini is correctly pulling tracklist
- spot check some tracks where track url is missing:
  - Can spotify FindTrack be improved? (a la ft. and feat.)
    - should I make a CleanArtistName method?
  - Would that track have been found with google search (harder to check now)
  - Is the track just missing from Spotify Generally.

It's gonna be hard to run this migration task with the Google Search limit. I guess just patience, a set of videos a day

Maybe lets not run any tonight, so then by tomorrow I can run it at any hour

It's really hard to manage this migration with google searh limit

Should be so fine, for regular weekly runs. But not for this case!

Be patient. 100 vids every 24 hours I think
After 100 videos today
Review, how many google search requests was that?

Deploy lambda and do the last few videos from manual invoke!!

omg I ran outta quota JUST for first 100 videos

got to 892 of 950 tracks

There were more tracks than I expected
I guess cos I filtered out some invalid videos today there were more tracks to query
I also just fixed a bug where it was including ...meh... tracks. So maybe it was searching for more cos of that too.
Lets run in full and see total tracks
yeah that's lower now
and I fixed a bug meaning more songs would not be found in Spotify so would need google search
Confident that another run with limit 100 would WORK

golang lambda size much larger cos it's including golang runtime too
is lambda size important anyway? Surely not for my purposes anyway
nodejs ~1mb
golang ~10mb

Noticed that Joey Badass Google Searched returned a totally random spotify track:
THE FINALS Joey Bada$$ -> https://open.spotify.com/track/6xU31XJHAwgAjfj8jBxf2k
Perhaps tracks found by google search should be marked so they can be manually reviewed
I feel like it wouldn't be too bad to do that via migration
Go thru Google Sheet
If you can't find the track from spotify.findTrack(), mark it
Then replace rows
Another job for another time, when I feel like it.

I should use customsearch results "title": "symbol - song and lyrics by Adrianne Lenker | Spotify"

But then what's the column structure I want?
Title(from youtube description)
Artist(from youtube description)
Source("", Spotify, GoogleSearch) - I think that's still OK
QueryResult: either Spotify Track data or Google Custom Search Data
....

I think that's good!

Ah but migrating and including custom search data means RE-SEARCHING everything, Jaysus
I guess that's OK....

TODO finish migration tomorrow! run again

test google search on this Cold Blood HEALTH x LAMB OF GOD
cos it worked once and then failed a second time?

Stray Dog Popcaan worked first round
failed first migration
then worked on second migration run? Since it wasn't excluded when script re-ran

I think I've manually reviewed all historic Google Search results now
Still keen to run a script to compare old playlist items to new ones and see the difference

and should definitely review those items which I was able to manually resolve that Google Search couldn't handle. See if I can improve the google search. Cos it's not a great hitrate really, maybe 50%? I wanna find that number too
32 manually resolved
42 found with google search

So just over 50%. It's better than nothing but not gr8 for sure

set up CI to build & test on commits?
https://github.com/JosephJvB/tony-g/new/main?filename=.github%2Fworkflows%2Fgo.yml&workflow_template=ci%2Fgo
https://docs.github.com/en/actions/automating-builds-and-tests/building-and-testing-go

v2
Can I update Youtube Playlists instead of Spotify Playlists?

Then I could move to using Youtube Music > Spotify

would need

1. new sheets in google sheets
2. new youtube apis

   - get playlists
   - search
   - create playlist
   - update playlist

3. remove googlesearch fallback since that was just to be a bandaid over spotify search
4. rm spotify apis
5. make a new repo? Or just a branch
   nah I think it makes sense to have a new repo
6. deploy as a brand new aws service I think so I can keep using spotify for now
   - would that have any issues with like api quota?
     - ie if both services are running and using the same API keys
     - gemini quota should be fine
     - only risk is if I was migrating videos to youtube and trying to parse loads of video descriptions as the spotify service was running

Wait is SoundCloud an option? Maybe it's the one?

Youtube would be super easy. Tony's links these days are all Youtube ones, I wouldn't even have to search for the track, he's got the video url and Id right there, it's v darn simple.
just get all videos, get all youtube links from description from best section, get id from url, make request to youtube. (would still need youtube search as a fallback in case the link isn't for a youtube video)

maybe even just add a new lambda here?

- tnd -> apple -> spotify (exists)
  - prob won't recreate this one (as youtube playlists) it has fewer songs so it's not as interesting!
- youtube -> spotify (exists)
  - becomes: youtube -> youtube (new)

I need new google auth since I'll be creating and updating youtube playlists that needs access/refresh_token stuff
https://developers.google.com/youtube/v3/guides/auth/server-side-web-apps
Specifically, content owners can use service accounts to call API methods that support the onBehalfOfContentOwner request parameter.

ugh I still need to do an oauth flow so smelly
i need this scope I think
https://www.googleapis.com/auth/youtube

https://accounts.google.com/o/oauth2/v2/auth
&client_id=
&redirect_uri=http://localhost:8080
&response_type=code
&scope=https://www.googleapis.com/auth/youtube

POST https://oauth2.googleapis.com/token
&client_id=
&redirect_uri=http://localhost:8080
&code=xxx

- Move to YT music
  - make sure new Youtube API methods are working
    - setAccessToken
      - turns out my existing refresh token had a 7 day expiry. Got a non-expiring one now! Wicky wings.
    - loadAllPlaylists (my playlists)
  - create new Youtube API methods
    - createPlaylist
    - addItemsToPlaylist
      - consider which order items are added
  - New Google Sheets Data
  - Shall I keep sorting with the Google Script? Or do it in code. I think it's better in code right. I'll do that later tho lol.

made a new spreadsheet for fear of messing up live one.
  - need a new apps script
  - I should really back this data up too since if I delete the sheet or w/e I'm stuffed
    - save to CSV? Probably...