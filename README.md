# 🎵 YouTube Audio Extractor Bot

Telegram bot in Go for extracting audio from YouTube videos. 
Users can send links to videos, and the bot will return audio files and save their history.

## Features
- Download audio from YouTube via link
- Save query information in PostgreSQL
- View history of downloaded audio files
- Fast query processing

## Stack
- **Programming language**: Go 1.21+
- **Database**: PostgreSQL
- **YouTube library**: [yt-dlp](https://github.com/yt-dlp/yt-dlp)
- **Audio processing**: FFmpeg
