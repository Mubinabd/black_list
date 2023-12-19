from pytube import YouTube
yt = YouTube('http://youtube.com/watch?v=Mup6uM2gdKk&list=RDMup6uM2gdKk&start_radio=1')
print(yt.title)
print(yt.thumbnail_url)
print(yt.streams)
print(yt.streams.filter(file_extension='mp4').first().url)
print(yt.streams.get_highest_resolution().download('video'))
video = yt.streams.get_by_itag(17)
video.download
