# Automato

Create automated youtube videos out of Reddit posts.

This is for educational purposes only, read more about the journey [here](https://blog.mrg.sh/side-project-post-automated-youtube-videos-from-reddit)

## Build

```bash
./buildit
```

## Prerequisites

Since I have decided to use `dotenv` to store configuration options, you will have to make sure that you have a `.env` file under `./build` directory that has these options:

```
DATABASE_URI=mongodb+srv://user:password@localhost/?retryWrites=true&w=majority
POSTS_LIMIT=5
BASE_DIR=/tmp
OUTPUT_PATH=/tmp/final.mp4
DOWNLOADER_THREADS=3
PROCESSOR_THREADS=1

SUBREDDIT=TikTokCringe
POSTS_TYPE=hot

TOKEN_FILE=.credentials.json
SECRETS_FILE=client_secrets.json
OUTPUT_FILE=/tmp/final.mp4
PRIVACY_STATUS=unlisted
VIDEO_INFO_TYPE=random
VIDEO_TITLE=Top 10 best/cringiest Tiktoks today
VIDEO_DESCRIPTION=Prepare yourself for the definitely not automated dose of Tiktoks
```

for publishing Youtube video, you'll need to have `./build/client_secrets.json` file that you download from [Google cloud API credentials console](https://console.cloud.google.com/apis/credentials)

## Run

```bash
./runit
```
