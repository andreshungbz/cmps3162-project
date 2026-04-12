document.addEventListener('DOMContentLoaded', () => {
  const video = document.getElementById('bg-video');
  if (!video) return;

  const shouldLoadVideo =
    window.innerWidth > 1000 &&
    !window.matchMedia('(prefers-reduced-motion: reduce)').matches;

  if (shouldLoadVideo) {
    const source = document.createElement('source');
    source.src = '/static/media/videos/ocean.mp4';
    source.type = 'video/mp4';

    video.appendChild(source);
    video.autoplay = true;
    video.load();
  }
});
