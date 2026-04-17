import { emitter } from '../../core/event-emitter.js';
import { state } from '../../core/state.js';
import { render } from './render.js';

export function setupHandlers() {
  const app = document.querySelector('#app');

  // Upload Media Toggle Button
  app.addEventListener('click', (e) => {
    const toggleBtn = e.target.closest('#toggle-upload-form');

    if (toggleBtn) {
      state.showcase.ui.formOpen = !state.showcase.ui.formOpen;
      render();
    }
  });

  // Upload Media Submit Button (STUB)
  app.addEventListener('submit', (e) => {
    if (e.target.id !== 'upload-form') return;
    e.preventDefault();

    const file = document.querySelector('#media-file').files[0];
    const filter = document.querySelector('#media-filter').value;

    if (!file) return;

    const fakeJob = {
      id: Date.now(),
      status: 'pending',
      filter,
      input_url: URL.createObjectURL(file),
      output_url: null,
    };

    emitter.emit('showcase:uploadCreated', fakeJob);
  });

  // Open Image Modal
  app.addEventListener('click', (e) => {
    const card = e.target.closest('[data-job-id]');
    if (!card) return;

    const job = state.showcase.data.find(
      (j) => j.id === Number(card.dataset.jobId),
    );

    state.showcase.ui.modalOpen = true;
    state.showcase.ui.selectedJob = job;

    render();
  });

  // Close Modal
  app.addEventListener('click', (e) => {
    const isCloseBtn = e.target.id === 'modal-close';
    const isBackdrop = e.target.classList.contains('modal-overlay');

    if (isCloseBtn || isBackdrop) {
      state.showcase.ui.modalOpen = false;
      state.showcase.ui.selectedJob = null;
      render();
    }
  });
}
