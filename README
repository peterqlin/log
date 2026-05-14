# **Daily Logger CLI (log)**

A lightweight Go CLI tool for logging daily activities. It appends timestamped entries to a daily markdown file (e.g., 2026.05.12.md), making it the perfect companion for Obsidian daily notes or general journaling.

If you add entries out of order (by backdating them), the tool uses binary search to automatically insert the new entry into the exact chronological spot within the file.

## **Features**

* **Daily Markdown Files**: Automatically creates or appends to a YYYY.MM.DD.md file.  
* **Chronological Sorting**: Uses binary search to insert backdated entries into the correct line based on time.  
* **12-Hour Format**: Timestamps are formatted as 03:04 PM.  
* **Time Offsets**: Log events that happened in the past using \--m-ago (minutes) and \--h-ago (hours).  
* **Fuzzy Logging**: Use \--around to apply a realistic, random \+/- 5-minute offset to your logs.  
* **Clean Output**: Ensures exactly one newline between all entries, cleaning up any messy whitespace.

## **Installation (Windows)**

1. Make sure you have [Go installed](https://go.dev/dl/).  
2. Clone or download the source code (main.go).  
3. Open your terminal in the directory containing main.go and compile the executable:  
   go build \-o log.exe main.go

4. **Add to PATH**: Move log.exe to a folder that is in your System's PATH (e.g., C:\\Program Files\\Go\\bin or a custom C:\\tools folder) so you can run the log command from anywhere.

## **Configuration**

Before logging your first activity, you need to tell the tool where to save your daily markdown files.

log \--set-dest "C:\\Users\\YourName\\Documents\\Obsidian\\Daily Notes"

*This will automatically create a configuration file at \~\\.config\\log\\config.toml to store your path.*

## **Usage**

Simply type log followed by your activity.

### **Basic Logging**

Logs the activity at the current exact time.

log went to the gym

*(Output: 02:15 PM \- went to the gym)*

### **Backdating**

Forgot to log something earlier? Use \--m-ago (minutes) and/or \--h-ago (hours).

log \--m-ago 30 finished reading chapter 4  
log \--h-ago 2 ate lunch at the cafe  
log \--h-ago 1 \--m-ago 15 had a meeting

### **Fuzzy Logging ("Around")**

If you don't want your logs to look too perfectly timed to the exact minute you ran the command, use the \--around flag. It applies a random \+/- 5-minute offset to the calculated time.

log \--m-ago 60 \--around woke up

*(If it's currently 09:00 AM, the log time will randomly land somewhere between 07:55 AM and 08:05 AM).*

## **Example Output**

If you check your destination folder, you'll see a file named 2026.05.14.md looking exactly like this:

08:02 AM \- woke up  
09:15 AM \- had a meeting  
12:00 PM \- ate lunch at the cafe  
02:15 PM \- went to the gym  
