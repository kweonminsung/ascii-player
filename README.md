# 🎬 Console Cinema

**Watch videos as real-time ASCII or Pixel art directly in your terminal!**

Console Cinema is a powerful tool that converts and plays local video files or YouTube videos into vibrant ASCII or Pixel art right in your command-line environment.

## ✨ Features

- **Real-time Conversion**: Plays videos by converting them to ASCII or Pixel art in real-time.
- **Local & YouTube Support**: Supports both local video files (MP4, AVI, etc.) and YouTube URLs.
- **Multiple Art Styles**: Offers two distinct art styles: `ascii` and `pixel`.
- **Simple to Use**: Designed with an intuitive command structure for easy operation.

## 🚀 Installation

If you have a Go environment set up, you can easily install it with the following command:

```bash
go install github.com/kweonminsung/console-cinema@latest
```

Alternatively, you can clone this repository and build it yourself.

```bash
git clone https://github.com/kweonminsung/console-cinema.git
cd console-cinema
go build
```

## 📖 Usage

### Playing Local Videos

Use the `play` command to play local video files.

```bash
# Play in the default mode (ascii)
./console-cinema play test.mp4

# Play in pixel mode
./console-cinema play test.mp4 --mode pixel

# Set frames per second (fps)
./console-cinema play test.mp4 --fps 30
```

### Playing YouTube Videos

Use the `youtube play` command to play a video from a YouTube link. Supports standard videos, shorts, and embed URLs.

```bash
# Play in the default mode (pixel)
./console-cinema youtube play "https://www.youtube.com/watch?v=your_video_id"

# Play a Short in ascii mode
./console-cinema youtube play "https://www.youtube.com/shorts/your_short_id" --mode ascii
```

### Exploring YouTube Videos

Use the `youtube explore` command to search for videos on YouTube and select one to play.

```bash
./console-cinema youtube explore
```

This will open an interactive search interface in your terminal.

### Available Commands

```
Console Cinema - A real-time ASCII/Pixel art video player for the command line

Usage:
  console-cinema [command]

Available Commands:
  play        Play local video files (MP4, AVI, etc.)
  youtube     Play or explore YouTube videos
  help        Help about any command

Flags:
  -h, --help   help for console-cinema
```

## 📄 License

This project is licensed under the MIT License.
