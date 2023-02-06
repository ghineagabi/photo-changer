# photo-changer
Adds sine waves as white noise to a photo and converts individual photos to a video

Required ffmpeg installed and accessible via env variables

ffmpeg command to generate the video from jpg's : ffmpeg -framerate 10 -i "converted/%03d.jpeg" -vf "pad=ceil(iw/2)*2:ceil(ih/2)*2" output.mp4

https://user-images.githubusercontent.com/102144337/216997426-c5351a6e-aa7a-4f94-a51a-bac58fa9a6b6.mp4

