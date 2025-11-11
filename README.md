# 🍅 Pomodoro CLI

A command-line Pomodoro timer to help you stay focused and productive. Work in focused intervals, take breaks, and optionally block distracting websites during your work sessions.

## ✨ Features

- ⏰ **Customizable work and break intervals** - Set your own timing
- 🔄 **Multiple Pomodoro cycles** - Run several rounds automatically
- 🚫 **Website blocking** - Block distracting sites during work time (automatically unblocks during breaks)
- 🔊 **Audio feedback** - Windup sound at start, ticking during work, and ding when complete
- 🤫 **Silent mode** - Disable ticking sound if you prefer quiet focus
- 📊 **Visual progress bar** - countdown bar with progress indicator
- ⌨️ **Graceful interruption** - Cancel anytime with Ctrl+C and sites will be automatically unblocked

## 🎯 What is the Pomodoro Technique?

The [Pomodoro Technique](https://en.wikipedia.org/wiki/Pomodoro_Technique) is a time management method that breaks work into focused intervals (traditionally 25 minutes) separated by short breaks (usually 5 minutes). This helps maintain high levels of focus and prevents burnout.

## 📋 Requirements

- **macOS** (this tool is designed for macOS only)
- **Go 1.16+** (for building from source)
- **sudo privileges** (only needed if using the website blocking feature)

## 🚀 Installation

### Option 1: Build from source

```bash
# Clone the repository
git clone https://github.com/maiku1008/pomodoro-cli.git
cd pomodoro-cli

# Build the binary
go build -o pomodoro

# Move to your PATH (optional)
sudo mv pomodoro /usr/local/bin/
```

### Option 2: Direct go install

```bash
go install github.com/maiku1008/pomodoro-cli@latest
```

## 📖 Usage

### Basic Examples

**Standard 25-minute Pomodoro with 5-minute break:**
```bash
./pomodoro
```

**Custom timing - 45 minutes work, 10 minutes break:**
```bash
./pomodoro -timer 45 -break 10
```

**Run 4 Pomodoro cycles in a row:**
```bash
./pomodoro -interval 4
```

**Silent mode (no ticking sound):**
```bash
./pomodoro -silent
```

**Block distracting websites during work:**
```bash
sudo ./pomodoro -blocklist "twitter.com,reddit.com,youtube.com"
```

**Full-featured session:**
```bash
sudo ./pomodoro -timer 25 -break 5 -interval 4 -blocklist "twitter.com,facebook.com,reddit.com"
```

### 🎮 Real-World Examples

**Deep work session (2 hours of focused work):**
```bash
sudo ./pomodoro -timer 50 -break 10 -interval 2 -blocklist "twitter.com,reddit.com,news.ycombinator.com"
```

**Quick task sprint (15 minutes):**
```bash
./pomodoro -timer 15 -break 3
```

**Study session with social media blocked:**
```bash
sudo ./pomodoro -timer 25 -break 5 -interval 4 -blocklist "instagram.com,tiktok.com,twitter.com,facebook.com"
```

**Late night coding (silent mode):**
```bash
./pomodoro -timer 30 -break 5 -silent -blocklist "youtube.com,reddit.com"
```

## ⚙️ Configuration Options

| Flag | Default | Description |
|------|---------|-------------|
| `-timer` | `25` | Work session duration in minutes |
| `-break` | `5` | Break duration in minutes |
| `-interval` | `1` | Number of Pomodoro cycles to complete |
| `-blocklist` | `""` | Comma-separated list of websites to block during work (e.g., "twitter.com,reddit.com") |
| `-hosts` | `/etc/hosts` | Path to hosts file (only change if you know what you're doing) |
| `-silent` | `false` | Disable the ticking sound during work sessions |

## 🔒 Website Blocking

The website blocking feature works by temporarily adding entries to your system's `/etc/hosts` file, redirecting blocked sites to `127.0.0.1` (localhost).

### Important Notes:

- **Requires sudo:** You must run the command with `sudo` when using `-blocklist`
- **Automatic cleanup:** Sites are automatically unblocked when:
  - The work session ends (for your break)
  - You press Ctrl+C to cancel
  - The program exits for any reason
- **Format:** Just provide the domain name, don't include `http://` or `www.` (both will be blocked automatically)
- **DNS cache:** Some browsers cache DNS. If a site isn't blocked immediately, try:
  - Chrome/Brave: Open a new incognito window
  - Firefox: Restart the browser
  - Safari: Usually respects hosts file immediately

### What gets blocked?

When you block `example.com`, both of these are blocked:
- `example.com`
- `www.example.com`

### Example blocked domains:
```bash
-blocklist "twitter.com,reddit.com,youtube.com,facebook.com,instagram.com,tiktok.com,news.ycombinator.com"
```

## 🎵 Sound Features

The app includes three embedded sounds (macOS only, uses `afplay`):
- **Windup** 🎬 - Plays at the start of each work session
- **Ticking** ⏱️ - Gentle ticking during work (can be disabled with `-silent`)
- **Ding** 🔔 - Notification when a session completes

All sounds are embedded in the binary, so no external files are needed!

## 💡 Tips & Best Practices

1. **Start with standard Pomodoros:** Try the default 25/5 split before experimenting
2. **Use website blocking wisely:** Block your biggest time-wasters during deep work
3. **Take your breaks seriously:** Step away from your computer, stretch, hydrate
4. **Experiment with timing:** Some tasks need longer focus periods (50 min), others benefit from shorter sprints (15 min)
5. **Silent mode for shared spaces:** Use `-silent` when working in libraries or offices
6. **Chain multiple cycles:** Use `-interval 4` for a full work session without having to restart

## 🛠️ Troubleshooting

### "Permission denied" error when using blocklist
**Problem:** You're trying to modify `/etc/hosts` without sufficient permissions.

**Solution:** Run the command with `sudo`:
```bash
sudo ./pomodoro -blocklist "twitter.com,reddit.com"
```

### Sites aren't being blocked
**Possible solutions:**
1. Make sure you're using `sudo` with the `-blocklist` flag
2. Try flushing your DNS cache:
   - **macOS:** `sudo dscacheutil -flushcache; sudo killall -HUP mDNSResponder`
   - **Linux:** `sudo systemd-resolve --flush-caches` (if using systemd)
3. Open an incognito/private browsing window
4. Check that the sites are actually in your hosts file: `cat /etc/hosts | grep Pomodoro`

### Sites remain blocked after cancelling
**This shouldn't happen**, but if it does:
1. Check your hosts file: `cat /etc/hosts`
2. Manually remove the Pomodoro section between these markers:
   ```
   ### Pomodoro CLI - Begin Blocked sites ###
   ...
   ### Pomodoro CLI - End Blocked sites ###
   ```
3. Save the file and flush DNS cache (see commands above)

### No sound playing
1. Check your system volume
2. Make sure no other audio is playing
3. Verify `afplay` works: `afplay /System/Library/Sounds/Ping.aiff`

## 🎨 What You'll See

```
🍅 Starting 2 Pomodoro cycle(s)

═══════════════════════════════════════
🍅 Pomodoro 1 of 2
═══════════════════════════════════════

⏰ Work session (25 minutes)
Blocking distracting sites...
🍅 [████████████████░░░░░░░░░░░░░░] 12:34 remaining

✅ Work session complete!
Unblocking sites for break...

☕ Break time! (5 minutes)
Sites are now unblocked. Take a break!
☕ [██████████████████████████████] 00:05 remaining

⏰ Break finished!

✨ Pomodoro 1 complete!

═══════════════════════════════════════
🍅 Pomodoro 2 of 2
═══════════════════════════════════════
...
```

## 🤝 Contributing

Issues and pull requests are welcome! Feel free to:
- Report bugs
- Suggest new features
- Improve documentation
- Share your Pomodoro workflows

## 📝 License

MIT License - feel free to use this however you'd like!

## 🙏 Acknowledgments

This program is inspired by [TomatoBar](https://github.com/ivoronin/TomatoBar), a neat Pomodoro timer for the macOS menu bar.

Timer sounds are licensed from buddhabeats.

Built with Go and inspired by the timeless Pomodoro Technique created by Francesco Cirillo.

---

**Happy focusing! 🍅✨**

*Remember: The best productivity tool is the one you actually use. Keep it simple, stay consistent, and take those breaks!*
