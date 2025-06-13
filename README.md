# Video Renamer

A command-line utility to rename video files by appending their duration in a concise format.

## Features

- Renames video files by appending their duration in `[XhYm]` or `[Ym]` format.
- Supports various video formats: `.mp4`, `.mov`, `.avi`, `.mkv`, `.wmv`, `.flv`, `.webm`.
- Optionally process a single file or a directory of video files.
- Recursive directory processing with the `-R` option.
- Helpful command-line options for usage information.

## Prerequisites

- Go 1.21 or later
- FFmpeg installed on your system

## Installation

1. Clone the repository:
   ```sh
   git clone <repository-url>
   cd video-renamer
   ```

2. Build the executable:
   ```sh
   go build -o video-renamer
   ```

## Usage

- Process the current directory:
  ```sh
  ./video-renamer
  ```

- Process a specific file:
  ```sh
  ./video-renamer /path/to/video.mp4
  ```

- Process a directory:
  ```sh
  ./video-renamer /path/to/directory
  ```

- Process a directory recursively:
  ```sh
  ./video-renamer /path/to/directory -R
  ```

- Display help message:
  ```sh
  ./video-renamer -h
  ```

## Example

If you have a video file named `my_video.mp4` that is 1 hour, 23 minutes, and 45 seconds long, it will be renamed to `my_video_[1h24m].mp4`.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 