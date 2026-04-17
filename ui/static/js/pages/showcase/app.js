import { emitter } from '../../core/event-emitter.js';
import { state } from '../../core/state.js';

import { render } from './render.js';
import { setupHandlers } from './handlers.js';

// showcase:fetch event -> fetch all media jobs (stub for later /media endpoint)
emitter.on('showcase:fetch', () => {
  state.showcase.loading = true;
  state.showcase.error = null;
  render();

  // STUB: later replace with real API call
  emitter.emit('showcase:fetched', []);
});

// showcase:fetched event -> update showcase data and rerender
emitter.on('showcase:fetched', (data) => {
  state.showcase.loading = false;
  state.showcase.data = data;
  render();
});

// showcase:uploadCreated event -> triggered after POST /showcase returns job_id
emitter.on('showcase:uploadCreated', (job) => {
  // 1. add job to state
  state.showcase.data.unshift(job);

  // 2. open modal immediately
  state.showcase.ui.modalOpen = true;
  state.showcase.ui.selectedJob = job;

  render();
});

// showcase:jobUpdated event -> observer pattern hook (polling later)
emitter.on('showcase:jobUpdated', (updatedJob) => {
  const idx = state.showcase.data.findIndex((j) => j.id === updatedJob.id);
  if (idx !== -1) {
    state.showcase.data[idx] = updatedJob;
  }
  render();
});

// showcase:error event -> set error message
emitter.on('showcase:error', (err) => {
  state.showcase.loading = false;
  state.showcase.error = err;
  render();
});

// INITIALIZATION
setupHandlers();
emitter.emit('showcase:fetch');
