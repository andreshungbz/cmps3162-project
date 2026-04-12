document.addEventListener('DOMContentLoaded', () => {
  const video = document.getElementById('bg-video');
  if (!video) return;

  let loaded = false;

  function tryLoadVideo() {
    if (loaded) return;

    const shouldLoadVideo =
      window.innerWidth > 1000 &&
      !window.matchMedia('(prefers-reduced-motion: reduce)').matches;

    if (!shouldLoadVideo) return;

    const source = document.createElement('source');
    source.src = '/static/media/videos/ocean.mp4';
    source.type = 'video/mp4';

    video.appendChild(source);
    video.autoplay = true;
    video.load();

    loaded = true; // ensure it only happens once
  }

  // run on load
  tryLoadVideo();

  // run on resize
  window.addEventListener('resize', tryLoadVideo);
});
