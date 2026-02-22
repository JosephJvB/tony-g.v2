package gemini

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestGemini(t *testing.T) {
	// parses description well
	// not good for finding spotify urls
	// or fixing typos
	t.Run("Can parse test youtube description 2019 Weekly Track Roundup: 10/6", func(t *testing.T) {
		t.Skip("Skip calling real Gemini API")

		err := godotenv.Load("../../.env")
		if err != nil {
			t.Errorf("Error loading .env file")
		}

		description := "CHARITY COMPILATION PRE-ORDER:\nhttps://theneedledrop.merchtable.com/music/the-needle-drop-various-artists-vinyl-12\nILRC: https://www.ilrc.org/\n\nPatreon: https://www.patreon.com/theneedledrop\n\nFAV TRACKS Spotify playlist: https://open.spotify.com/user/tndausten/playlist/2zderg88f9HbH54RJBTp1m?si=W8oXCAHvRnSJun4x6VHhdQ\n\nTurntable Lab link: http://turntablelab.com/theneedledrop\n\nAmazon link: http://amzn.to/1KZmdWI\n\n\nSHOUTOUT: CARLY RAE JEPSEN COVERS NO DOUBT\nhttps://open.spotify.com/album/0tUnCfqBLeZlivAkbxvbib\n\n\n!!!BEST TRACKS THIS WEEK!!!\n\nDoja Cat - Bottom Bitch\nhttps://youtu.be/ik0qg-O_2DM\n\nPoppy - I Disagree\nhttps://www.youtube.com/watch?v=6gmswmbosYo&feature=youtu.be\n\nMoor Mother - After Images\nhttps://youtu.be/VeZIqemkrD8\n\nLightning Bolt - Hüsker Dön't\nhttps://lightningbolt.bandcamp.com/track/h-sker-d-nt\n\nNegative Gemini - Bad Baby (Club Mix)\nhttps://youtu.be/_ddbrUq40Iw\n\nAnamanaguchi - Air on Line\nhttps://youtu.be/nnq1ApucY4g\n\nBig Thief - Forgotten Eyes\nhttps://youtu.be/hGD-8f8Wn5M\n\nG.T. - How Dare You\nhttps://youtu.be/rbrdRcwZE6Q\n\nTyler, the Creator - Earfquake (Channel Tres Remix)\nhttps://www.youtube.com/watch?v=T8jx0d9GAF4\n\nDaniel Caesar - CYANIDE REMIX ft. Koffee\nhttps://www.youtube.com/watch?v=mBKXHk2nJ1I\n\nKim Petras - There Will Be Blood\nhttps://www.youtube.com/watch?v=8nBQ8xv2oLY\n\nSunn O))) - Frost (C)\nhttps://youtu.be/Y20qC3qgpps\n\nJacques Greene - For Love\nhttps://youtu.be/GzdMcHhM7tQ\n\nGreat Grandpa - Bloom\nhttps://youtu.be/jFs4Tliyjpg\n\nFloating Points - Anasickmodular\nhttps://youtu.be/Md9gjJlqAxQ\n\nclipping. - Blood on the Fang\nhttps://youtu.be/s9EsHbqmjN4\n\n\n...meh...\n\nRemo Drive - Romeo\nhttps://youtu.be/1DiNlZMBPY0\n\nChromatics - You're No Good\nhttps://youtu.be/PjUblmk4Cyo\n\nGuapdad 4000 - Gucci Pajamas ft. Chance the Rapper & Charlie Wilson\nhttps://www.youtube.com/watch?v=QLw2eTCKaCg\n\nKing Princess - Hit the Back\nhttps://www.youtube.com/watch?v=GyFsbYSajhs\n\nGucci Mane - Big Booty ft. Megan Thee Stallion\nhttps://www.youtube.com/watch?v=b_Kx8tx88oQ\n\nBen Frost - Catastrophic Deliquescence\nhttps://youtu.be/7HNqV3K7di8\n\nBuju Banton - Lend a Hand\nhttps://youtu.be/xUtrvvHre34\n\nSleater-Kinney - ANIMAL\nhttps://youtu.be/pGOO7EE4Lhw\n\nSummer Walker - Playing Games ft. Bryson Tiller\nhttps://www.youtube.com/watch?v=o_6HGBsMHeA\n\nJuice WRLD - Bandit ft. NBA Youngboy\nhttps://www.youtube.com/watch?v=Sw5fNI400E4\n\nTravis Scott - Highest in the Room\nhttps://www.youtube.com/watch?v=tfSS1e3kYeo\nReview: https://www.youtube.com/watch?v=mjVdNIw9LMk\n\nDanny Brown - 3 Tearz ft. Run the Jewels (prod. JPEGMAFIA)\nhttps://www.youtube.com/watch?v=ApJ1_ZliXLQ\nReview: https://www.youtube.com/watch?v=WB625tK_FK0\n\nCHVRCHES - Death Stranding\nhttps://youtu.be/mFGq92BYmt4\n\nEOB - Santa Teresa\nhttps://youtu.be/TG-Od2-OTdg\n\n\n!!!WORST TRACKS THIS WEEK!!!\n\nNONE! Yeah, none!\n\n===================================\nSubscribe: http://bit.ly/1pBqGCN\n\nOfficial site: http://theneedledrop.com\n\nTND Twitter: http://twitter.com/theneedledrop\n\nTND Facebook: http://facebook.com/theneedledrop\n\nSupport TND: http://theneedledrop.com/support\n===================================\n\nY'all know this is just my opinion, right?"

		apiKey := os.Getenv("GEMINI_API_KEY")

		client := NewClient(apiKey)

		result := client.ParseYoutubeDescription(description)

		fmt.Printf("got tracks %d", len(result))

		if len(result) == 13 {
			t.Errorf("failed to load youtube tracks from description. Got %d, Expected 13", len(result))
		}
	})

	// ooh now this works v nicely
	t.Run("Can handle case where description has two tracks (limerence/ankles)", func(t *testing.T) {
		t.Skip("Skip calling real Gemini API")

		err := godotenv.Load("../../.env")
		if err != nil {
			t.Errorf("Error loading .env file")
		}

		description := "2025 FAV TRACKS PLAYLIST: https://music.apple.com/us/playlist/my-fav-singles-of-2025/pl.u-ayeZTygbKDy\n\nTND Patreon: https://www.patreon.com/theneedledrop\n\nTurntable Lab link: http://turntablelab.com/theneedledrop\n\nAUSTEN SHOUTOUT\nventuring - Dead forever (\u0026 other singles)\nhttps://www.youtube.com/watch?v=DV3yteStUk0\u0026list=OLAK5uy_k-xM5zQh-RNWGoFK2K6FLqTSfACQI3_mc\u0026index=1\n\n\n!!!BEST TRACKS THIS WEEK!!!\n\nLucy Dacus - Ankles: https://www.youtube.com/watch?v=pcW_-uxy6dQ\u0026pp=ygUfTHVjeSBEYWN1cyAtIEFua2xlcyAvIExpbWVyZW5jZQ%3D%3D\nLimerence: https://www.youtube.com/watch?v=re3mFdbzJQ8\u0026pp=ygUfTHVjeSBEYWN1cyAtIEFua2xlcyAvIExpbWVyZW5jZQ%3D%3D\n\nGates to Hell - Next to Bleed\nhttps://www.youtube.com/watch?v=kTGHyGHgwJ0\u0026pp=ygUdR2F0ZXMgdG8gSGVsbCAtIE5leHQgdG8gQmxlZWQ%3D\n\nHorsegirl - Switch Over\nhttps://www.youtube.com/watch?v=mC1v7Y7bIKs\u0026pp=ygUXSG9yc2VnaXJsIC0gU3dpdGNoIE92ZXI%3D\n\nSaya Gray - SHELL ( OF A MAN )\nhttps://www.youtube.com/watch?v=KYM1BbMaoco\u0026pp=ygUeU2F5YSBHcmF5IC0gU0hFTEwgKCBPRiBBIE1BTiAp\n\nPerfume Genius - It's a Mirror\nhttps://www.youtube.com/watch?v=hx2_NGaDPrk\u0026pp=ygUeUGVyZnVtZSBHZW5pdXMgLSBJdCdzIGEgTWlycm9y\n\nT-Pain, Girl Talk, Yaeji - Believe in Ya\nhttps://www.youtube.com/watch?v=n_Usx_hhtiQ\u0026pp=ygUoVC1QYWluLCBHaXJsIFRhbGssIFlhZWppIC0gQmVsaWV2ZSBpbiBZYQ%3D%3D\n\nBaths - Eden\nhttps://www.youtube.com/watch?v=3N6EI_oHL1I\u0026pp=ygUMQmF0aHMgLSBFZGVu\n\n\n...meh...\n\nJohn Glacier - Ocean Steppin' ft. Sampha\nhttps://www.youtube.com/watch?v=g-r3wAdSWSE\u0026pp=ygUoSm9obiBHbGFjaWVyIC0gT2NlYW4gU3RlcHBpbicgZnQuIFNhbXBoYQ%3D%3D\n\nTim Hecker - Sunset Key Melt\nhttps://www.youtube.com/watch?v=26jeWB9Aw8c\u0026pp=ygUcVGltIEhlY2tlciAtIFN1bnNldCBLZXkgTWVsdA%3D%3D\n\nImperial Triumphant - Lexington Delirium ft. Tomas Haake\nhttps://www.youtube.com/watch?v=v9cDvwwbj6A\u0026pp=ygU4SW1wZXJpYWwgVHJpdW1waGFudCAtIExleGluZ3RvbiBEZWxpcml1bSBmdC4gVG9tYXMgSGFha2U%3D\n\nLogic - French Dispatch\nhttps://www.youtube.com/watch?v=Elj44V1HiVk\u0026pp=ygUXTG9naWMgLSBGcmVuY2ggRGlzcGF0Y2g%3D\n\nOklou - take me by the hand ft. Bladee\nhttps://www.youtube.com/watch?v=jdU16tnrt14\u0026pp=ygUmT2tsb3UgLSB0YWtlIG1lIGJ5IHRoZSBoYW5kIGZ0LiBCbGFkZWU%3D\n\nJungle - Keep Me Satisfied\nhttps://www.youtube.com/watch?v=fwq5sT-zLLk\u0026pp=ygUaSnVuZ2xlIC0gS2VlcCBNZSBTYXRpc2ZpZWQ%3D\n\nSharon Van Etten \u0026 The Attachment Theory - Trouble\nhttps://www.youtube.com/watch?v=hu3aQKiq-hk\u0026pp=ygUyU2hhcm9uIFZhbiBFdHRlbiAmIFRoZSBBdHRhY2htZW50IFRoZW9yeSAtIFRyb3VibGU%3D\n\n\n!!!WORST TRACKS THIS WEEK!!!\n\nImagine Dragons - Dare U ft. NLE Choppa\nhttps://www.youtube.com/watch?v=NObExZCktuM\u0026pp=ygUnSW1hZ2luZSBEcmFnb25zIC0gRGFyZSBVIGZ0LiBOTEUgQ2hvcHBh\n\nSpin Doctors - Still a Gorilla\nhttps://www.youtube.com/watch?v=vULjnd8YUMw\u0026pp=ygUiVGhlIFNwaW4gRG9jdG9ycyAtIFN0aWxsIGEgR29yaWxsYQ%3D%3D\n\nCentral Cee - GBP ft. 21 Savage\nhttps://www.youtube.com/watch?v=_Cu9Df_9Zvg\u0026pp=ygUfQ2VudHJhbCBDZWUgLSBHQlAgZnQuIDIxIFNhdmFnZQ%3D%3D\n\n===================================\nSubscribe: http://bit.ly/1pBqGCN\n\nPatreon: https://www.patreon.com/theneedledrop\n\nOfficial site: http://theneedledrop.com\n\nTwitter: http://twitter.com/theneedledrop\n\nInstagram: https://www.instagram.com/afantano\n\nTikTok: https://www.tiktok.com/@theneedletok\n\nTND Twitch: https://www.twitch.tv/theneedledrop\n===================================\n\nY'all know this is just my opinion, right?"

		apiKey := os.Getenv("GEMINI_API_KEY")

		client := NewClient(apiKey)

		tracks := client.ParseYoutubeDescription(description)

		d, err := json.MarshalIndent(tracks, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/gemini-description-resp.json", d, 0666)
		if err != nil {
			panic(err)
		}

		if len(tracks) != 8 {
			t.Errorf("Failed to get all tracks from youtube description. Got %d, expected 8", len(tracks))
		}
	})

	t.Run("Can get correct properties for video EGwEmD7EfXg", func(t *testing.T) {
		t.Skip("Skip calling real Gemini API")

		err := godotenv.Load("../../.env")
		if err != nil {
			t.Errorf("Error loading .env file")
		}

		description := "Amazon link:\nhttp://amzn.to/1KZmdWI\n\nScHoolboy Q - Tookie Knows II: Part (2)\nhttp://www.theneedledrop.com/articles/2016/7/schoolboy-q-tookie-knows-ii-part-2\n\nBADBADNOTGOOD - In Your Eyes ft. Charlotte Day Wilson\nhttp://www.theneedledrop.com/articles/2016/7/badbadnotgood-in-your-eyes-ft-charlotte-day-wilson\n\nMaxo Kream - The Persona Tape\nhttp://www.theneedledrop.com/articles/2016/7/3c72tw2w4uiygl2yszyjrn6wckts5l\n\nThelonious Martin - Bomaye ft. Joey Purp\nhttp://www.theneedledrop.com/articles/2016/7/thelonious-martin-bomaye-ft-joey-purp\n\nclipping. - Wriggle (music vid)\nhttp://www.theneedledrop.com/articles/2016/7/clipping-wriggle\n\nAngel Olsen - Shut Up Kiss Me\nhttp://www.theneedledrop.com/articles/2016/7/angel-olsen-shut-up-kiss-me\n\nGhoul - Bringer of War\nhttp://www.theneedledrop.com/articles/2016/7/ghoul-bringer-of-war\n\n===================================\nSubscribe: http://bit.ly/1pBqGCN\n\nOfficial site: http://theneedledrop.com\n\nTND Twitter: http://twitter.com/theneedledrop\n\nTND Facebook: http://facebook.com/theneedledrop\n\nSupport TND: http://theneedledrop.com/support\n===================================\n\nFAV TRACKS:\n\nLEAST FAV TRACK:\n\nArtist- Album / Year / Label / Genre\n\n/10\n\nY'all know this is just my opinion, right?"

		apiKey := os.Getenv("GEMINI_API_KEY")

		client := NewClient(apiKey)

		tracks := client.ParseYoutubeDescription(description)

		d, err := json.MarshalIndent(tracks, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/gemini-description-resp.json", d, 0666)
		if err != nil {
			panic(err)
		}
	})

	t.Run("Can get correct properties for video 2UADjU66-4M", func(t *testing.T) {
		t.Skip("Skip calling real Gemini API")

		err := godotenv.Load("../../.env")
		if err != nil {
			t.Errorf("Error loading .env file")
		}

		description := "Our sponsor: http://ridgewallet.com/fantano\nUSE PROMO CODE \"MELON\" FOR 10% OFF\n\n2022 FAV TRACKS PLAYLIST: https://music.apple.com/us/playlist/my-fav-singles-of-2022/pl.u-e92LIK9VM5K\n\nTND Patreon: https://www.patreon.com/theneedledrop\n\nTurntable Lab link: http://turntablelab.com/theneedledrop\n\n\n!!!BEST TRACKS THIS WEEK!!!\n\nFire-Toolz - Soda Lake with Game Genie / Vedic Software ~ Wet Interfacing\nhttps://fire-toolz.bandcamp.com/album/i-will-not-use-the-bodys-eyes-today?from=fanpub_fnb\n\nKing Gizzard \u0026 The Lizard Wizard - Ice V\nhttps://youtu.be/ydeV1_8pM4o\nReview: https://www.youtube.com/watch?v=fmqz_1GHmg0\n\nBjörk - Atopos ft. Kasimyn\nhttps://youtu.be/9FD2mUonh5s\n\nGilla Band - Backwash\nhttps://www.youtube.com/watch?v=q07rF2E-0Hw\n\nDeerhoof - My Lovely Cat!\nhttps://deerhoof.bandcamp.com/track/my-lovely-cat\n\nMile End - FCHC\nhttps://mileendband.bandcamp.com/track/fchc\n\n\n...meh...\n\nDry Cleaning - Gary Ashby\nhttps://youtu.be/XdvrSu38pWY\n\nBrian Eno - We Let It In\nhttps://youtu.be/Dehxp3PUTkM\n\nKEN mode - Throw Your Phone in the River\nhttps://www.youtube.com/watch?v=IwL556pzCXU\n\nBlood Orange - Jesus Freak Lighter\nhttps://www.youtube.com/watch?v=f21gWR8NdC0\n\nCordae \u0026 Hit-Boy - Checkmate\nhttps://www.youtube.com/watch?v=C7riiDNIv4A\n\nRun the Jewels - Opening Theme (From ATHF)\nhttps://www.youtube.com/watch?v=eRqbUVPs1kQ\n\nMura Masa \u0026 Erika de Casier - e-motions\nhttps://www.youtube.com/watch?v=x2iURw-BA5E\n\nAlex G - Miracles\nhttps://sandy.bandcamp.com/album/god-save-the-animals\n\nRuss - That Was Me\nhttps://www.youtube.com/watch?v=n-QCrCh5HbU\n\nThe Comet Is Coming - Lucid Dreamer\nhttps://www.youtube.com/watch?v=S7imxAIydR4\n\nWILLOW - curious/furious\nhttps://www.youtube.com/watch?v=MGYa0VIDpm4\n\n\n!!!WORST TRACKS THIS WEEK!!!\n\nNickleback - San Quentin\nhttps://www.youtube.com/watch?v=woA-qNpuwTs\n\nLewis Capaldi - Forget Me\nhttps://www.youtube.com/watch?v=FnbSkATyIc8\n \n\n===================================\nSubscribe: http://bit.ly/1pBqGCN\n\nPatreon: https://www.patreon.com/theneedledrop\n\nOfficial site: http://theneedledrop.com\n\nTwitter: http://twitter.com/theneedledrop\n\nInstagram: https://www.instagram.com/afantano\n\nTikTok: https://www.tiktok.com/@theneedletok\n\nTND Twitch: https://www.twitch.tv/theneedledrop\n===================================\n\nY'all know this is just my opinion, right?"

		apiKey := os.Getenv("GEMINI_API_KEY")

		client := NewClient(apiKey)

		tracks := client.ParseYoutubeDescription(description)

		d, err := json.MarshalIndent(tracks, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/gemini-description-resp.json", d, 0666)
		if err != nil {
			panic(err)
		}
	})

	t.Run("can handle x2 multi-tracks in video 17JnPv2fTjw", func(t *testing.T) {
		t.Skip("Skip calling real Gemini API")

		err := godotenv.Load("../../.env")
		if err != nil {
			t.Errorf("Error loading .env file")
		}

		description := "OUR SPONSOR: http://feedbands.com/needledrop\n\nAmazon link: http://amzn.to/1KZmdWI\n\nTurntable Lab link: http://turntablelab.com/theneedledrop\n\nIglooghost - Chalk Shrine MIX\nhttp://www.theneedledrop.com/articles/2017/5/iglooghost-chalk-shrine-mix\n\n!!!BEST TRACKS THIS WEEK!!!\n\nOmar Souleyman - Chobi\nhttp://www.theneedledrop.com/articles/2017/5/omar-souleyman-chobi\n\nCarly Rae Jepsen - Cut to the Feeling\nhttps://itun.es/us/Z_o0jb?i=1238429684\n\nBrockhampton - Gold / Heat / Face\nhttp://www.theneedledrop.com/articles/2017/5/brockhampton-gold\nhttps://youtu.be/Jpu0JZxDz-w\nhttps://youtu.be/_nWYiEq4wd0\n\nKirin J Callinan - Down 2 Hang / Living Each Day\nhttps://soundcloud.com/terrible-records/kirin-j-callinan-living-each-day\nhttps://soundcloud.com/terrible-records/kirin-j-callinan-down-2-hang\n\nShabazz Palaces - Since C.A.Y.A.\nhttps://www.youtube.com/watch?v=kCf2JrICz9Y\n\nalt-j - Adeline\nhttps://www.youtube.com/watch?v=1XwU8H6e8Ts\u0026feature=youtu.be\n\n...MEH...\n\nWashed Out - Get Lost\nhttps://open.spotify.com/album/3gHgPOe5PfA5jo2G1VawqG\n\nAni DiFranco - Zizzing ft. Justin Vernon\nhttp://www.stereogum.com/1942765/ani-difranco-zizzing-feat-justin-vernon/music/\n\nChromatics - Shadow\nhttp://www.theneedledrop.com/articles/2017/5/chromatics-shadow\n\n21 Savage - Issa ft. Young Thug \u0026 Drake\nhttp://www.stereogum.com/1942825/21-savage-issa-feat-drake-young-thug/music/\n\nBADBADNOTGOOD - To You (Andy Shauf Cover)\nhttps://open.spotify.com/album/4cuPfZ4xor2yMdakZOIReX\n\nBleachers - I Miss Those Days\nhttps://youtu.be/qQy12GH1Fl4\n\n!!!WORST TRACKS THIS WEEK!!!\n\nHalsey - Strangers ft. Lauren Jauregui\nhttps://www.youtube.com/watch?v=e3hjpNuvapQ\n\n Chuck Berry - Lady B. Goode\nhttps://itun.es/us/JPbiib?i=1210064091\n\n===================================\nSubscribe: http://bit.ly/1pBqGCN\n\nOfficial site: http://theneedledrop.com\n\nTND Twitter: http://twitter.com/theneedledrop\n\nTND Facebook: http://facebook.com/theneedledrop\n\nSupport TND: http://theneedledrop.com/support\n===================================\n\nY'all know this is just my opinion, right?"

		apiKey := os.Getenv("GEMINI_API_KEY")

		client := NewClient(apiKey)

		tracks := client.ParseYoutubeDescription(description)

		d, err := json.MarshalIndent(tracks, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/gemini-description-resp.json", d, 0666)
		if err != nil {
			panic(err)
		}
	})

	// This Is My Town / My Love For You Is Undying by Sun Kil Moon (Mark Kozelek)
	// OH! There's () in the artist name. Maybe that's the issue
	// wait Gemini tried to give me a MEH track!
	t.Run("Only gets best tracks. Not from ...Meh... section: gmlOJPAQHAE", func(t *testing.T) {
		t.Skip("Skip calling real Gemini API")

		err := godotenv.Load("../../.env")
		if err != nil {
			t.Errorf("Error loading .env file")
		}

		description := "FAV TRACKS Spotify playlist: https://open.spotify.com/user/tndausten/playlist/6eJIhC4KhMXDWrmheBW74m\n\nTurntable Lab link: http://turntablelab.com/theneedledrop\n\nAmazon link: http://amzn.to/1KZmdWI\n\n!!!BEST TRACKS THIS WEEK!!!\n\nJoey Bada$$ - THUGZ CRY\nhttps://youtu.be/73XVZ8jbKC4\n\nAlice Bag - Turn It Up\nhttps://youtu.be/Nt23p0h971k\n\nLet's Eat Grandma - Hot Pink (prod. SOPHIE)\nhttps://youtu.be/k0M3iIf2wdg\n\nUnknown Mortal Orchestra - American Guilt\nhttps://youtu.be/yFFa440Mo0I\n\nJulia Holter - So Humble The Afternoon\nhttp://www.adultswim.com/music/singles-2017/41\n\nJMSN - So Badly\nhttps://www.youtube.com/watch?time_continue=3\u0026v=sPaRSXdb4uU\n\nAlbert Hammond Jr. - Muted Beatings\nhttps://www.youtube.com/watch?v=YOFkQ9kh3CE\u0026ab_channel=RedBullRecords\n\n...MEH...\n\nSun Kil Moon (Mark Kozelek) - This Is My Town / My Love For You Is Undying\nhttp://www.theneedledrop.com/articles/2018/2/mark-kozelek-this-is-my-town-my-love-for-you-is-undying\n\nStreet Sects - Things Will Be Better In Hell\nhttp://www.brooklynvegan.com/street-sects-releasing-new-7-things-will-be-better-in-hell-stream-it/\n\nAndrew W.K. - Ever Again\nhttps://youtu.be/U_EupMUsb50\n\nChromeo - Bedroom Calling ft. The-Dream\nhttps://youtu.be/c1cvDqM_CMg\n\nCHVRCHES - Get Out\nhttps://youtu.be/LHUKKrcXfJs\n\nThe Streets - If You Ever Need To Talk I'm Here\nhttps://open.spotify.com/album/6TBVJXWURuv32vl6CbP5X6\n\nEd Schrader's Music Beat - Riddles\nhttps://youtu.be/rjhgPJXX7nw\n\n!!!WORST TRACKS THIS WEEK!!!\n\nIggy Azalea - Savior ft. Quavo\nhttps://soundcloud.com/iggy-azalea-official/savior-feat-quavo\n\nThe Weeknd, Kendrick Lamar - Pray For Me\nhttps://youtu.be/K5xERXE7pxI\n\n===================================\nSubscribe: http://bit.ly/1pBqGCN\n\nOfficial site: http://theneedledrop.com\n\nTND Twitter: http://twitter.com/theneedledrop\n\nTND Facebook: http://facebook.com/theneedledrop\n\nSupport TND: http://theneedledrop.com/support\n===================================\n\nY'all know this is just my opinion, right?"

		apiKey := os.Getenv("GEMINI_API_KEY")

		client := NewClient(apiKey)

		tracks := client.ParseYoutubeDescription(description)

		d, err := json.MarshalIndent(tracks, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/gemini-description-resp.json", d, 0666)
		if err != nil {
			panic(err)
		}

		if len(tracks) != 7 {
			t.Errorf("Failed to get Best tracks from description. Got %d, expected 7", len(tracks))
		}
	})

	t.Run("Can get correct tracks from: tnYKOhCMrj4", func(t *testing.T) {
		t.Skip("nah this is tonys mistake")
		// 1. Big Bad Wolf is not on Spotify
		// 2. the formatting is inconsistent: {artist} - {track1} & {track 2}
		//    usually it's separated by forward slash.

		err := godotenv.Load("../../.env")
		if err != nil {
			t.Errorf("Error loading .env file")
		}

		description := "Thanks SeatGeek for sponsoring the video. Get $20 off tix w/ code DROP: https://sg.app.link/DROP\n\nFAV TRACKS Spotify playlist: https://open.spotify.com/user/tndausten/playlist/6eJIhC4KhMXDWrmheBW74m\n\nTurntable Lab link: http://turntablelab.com/theneedledrop\n\nAmazon link: http://amzn.to/1KZmdWI\n\nSHOUTOUTS:\n\nGraham Lambkin:\nhttp://www.theneedledrop.com/articles/2018/1/graham-lambkin-no-better-no-worse-vol-1\n\nCar Seat Headrest reimagines more Twin Fantasy material:\nhttps://youtu.be/fj8H_ZXLgio\n\n!!!BEST TRACKS THIS WEEK!!!\n\nJames Blake - If The Car Beside You Moves Ahead\nhttps://www.youtube.com/watch?v=bYXM3uz1bjM\n\nJean Grae \u0026 Quelle Chris - OhSh ft. Hannibal Buress\nhttp://www.theneedledrop.com/articles/2018/1/jean-grae-quelle-chris-ohsh-ft-hannibal-buress\n\nThe Voidz - Leave It In My Dreams \u0026 QYURRYUS\nhttp://www.theneedledrop.com/articles/2018/1/the-voidz-leave-it-in-my-dreams\n\nDabrye - Lil Mufukuz ft. MF DOOM\nhttps://soundcloud.com/ghostly/dabrye-lil-mufukuz-feat-doom-1\n\nJack White - Corporation\nhttps://youtu.be/VFnXRntc9XA\n\nLil Wayne - Bloody Mary ft. Juelz Santana \u0026 Big Bad Wolf\nhttps://youtu.be/lYATz3STgew\nhttps://youtu.be/FOLQlEh1D20\n\n...MEH...\n\nHop Along - How Simple\nhttps://youtu.be/mf3H30pQ9ms\n\nMigos - Supastars\nhttps://youtu.be/hYc95vB-nT4\n\nAlice Glass - Cease \u0026 Desist\nhttps://youtu.be/QwswX7hRRgo\n\nFranz Ferdinand - Lazy Boy\nhttps://youtu.be/PNsUgvkjviY\n\nYoung Thug - MLK ft. Trouble and Shad Da God\nhttps://youtu.be/mbsSCJEP7oo\n\n!!!WORST TRACKS THIS WEEK!!!\n\nA$AP Rocky - Above / ☆☆☆☆☆ 5IVE $TAR$\nhttps://soundcloud.com/awgeshit/aap-rocky-above\nhttps://soundcloud.com/awgeshit/aap-rocky-5ive-tar\n\nBeck - I'm Waiting For The Man\nhttps://open.spotify.com/track/0IXsvVaa6mKZzziC95y40s\n\nAmanda Palmer - The Mess Inside\nhttps://youtu.be/Nug52TuotYU\n\nSting \u0026 Shaggy - Don't Make Me Wait\nhttps://youtu.be/cOaRPJQXFG4\n\n===================================\nSubscribe: http://bit.ly/1pBqGCN\n\nOfficial site: http://theneedledrop.com\n\nTND Twitter: http://twitter.com/theneedledrop\n\nTND Facebook: http://facebook.com/theneedledrop\n\nSupport TND: http://theneedledrop.com/support\n===================================\n\nY'all know this is just my opinion, right?"

		apiKey := os.Getenv("GEMINI_API_KEY")

		client := NewClient(apiKey)

		tracks := client.ParseYoutubeDescription(description)

		d, err := json.MarshalIndent(tracks, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/gemini-description-resp.json", d, 0666)
		if err != nil {
			panic(err)
		}

		if len(tracks) != 7 {
			t.Errorf("Failed to get Best tracks from description. Got %d, expected 7", len(tracks))
		}
	})
	// t.Run("can generate a confidence score", func(t *testing.T) {
	// 	// t.Skip("nah this is tonys mistake")

	// 	err := godotenv.Load("../../.env")
	// 	if err != nil {
	// 		t.Errorf("Error loading .env file")
	// 	}

	// 	apiKey := os.Getenv("GEMINI_API_KEY")

	// 	client := NewClient(apiKey)

	// 	d, err := json.MarshalIndent(tracks, "", "	")
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	err = os.WriteFile("../../data/gemini-description-resp.json", d, 0666)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	if len(tracks) != 7 {
	// 		t.Errorf("Failed to get Best tracks from description. Got %d, expected 7", len(tracks))
	// 	}
	// })
}
