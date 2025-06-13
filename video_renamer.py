#!/usr/bin/env python3

import os
import subprocess
import re
from pathlib import Path
import math

def get_video_duration(file_path):
    """Get video duration using ffprobe."""
    try:
        cmd = [
            'ffprobe',
            '-v', 'error',
            '-show_entries', 'format=duration',
            '-of', 'default=noprint_wrappers=1:nokey=1',
            str(file_path)
        ]
        result = subprocess.run(cmd, capture_output=True, text=True)
        duration = float(result.stdout.strip())
        return duration
    except Exception as e:
        print(f"Error getting duration for {file_path}: {e}")
        return None

def format_duration(duration):
    """Format duration in seconds to [XhYm] format, rounding up minutes."""
    total_minutes = math.ceil(duration / 60)
    hours = total_minutes // 60
    minutes = total_minutes % 60
    
    if hours > 0:
        return f"[{hours}h{minutes}m]"
    else:
        return f"[{minutes}m]"

def clean_filename(filename):
    """Remove all content within square brackets from filename."""
    return re.sub(r'\[.*?\]', '', filename).strip(' _')

def process_video_files():
    """Process all video files in the current directory."""
    video_extensions = {'.mp4', '.mov', '.avi', '.mkv', '.wmv', '.flv', '.webm'}
    
    for file_path in Path('.').iterdir():
        if file_path.suffix.lower() in video_extensions and file_path.is_file():
            duration = get_video_duration(file_path)
            if duration is not None:
                formatted_duration = format_duration(duration)
                # Clean the filename and add the new duration
                clean_name = clean_filename(file_path.stem)
                new_name = f"{clean_name}_{formatted_duration}{file_path.suffix}"
                try:
                    file_path.rename(new_name)
                    print(f"Renamed: {file_path.name} -> {new_name}")
                except Exception as e:
                    print(f"Error renaming {file_path.name}: {e}")

if __name__ == "__main__":
    process_video_files() 