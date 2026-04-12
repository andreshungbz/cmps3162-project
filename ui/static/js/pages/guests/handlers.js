import { emitter } from '../../core/event-emitter.js';
import { state } from '../../core/state.js';

// setupHandlers applies event listeners for guests pagination and filters.
export function setupHandlers() {
  const app = document.querySelector('#app');

  // Pagination Buttons (event delegation)
  app.addEventListener('click', (e) => {
    const btn = e.target.closest('[data-page]');
    if (!btn) return;

    emitter.emit('guests:pageChanged', Number(btn.dataset.page));
  });

  // Filter Apply Button
  app.addEventListener('click', (e) => {
    if (e.target.id !== 'apply-filters') return;

    // update name filter state
    state.guests.filters.name = document.querySelector('#filter-name').value;
    // update country filter state
    state.guests.filters.country =
      document.querySelector('#filter-country').value;
    // update page_size filter state
    state.guests.filters.page_size = Number(
      document.querySelector('#filter-page-size').value,
    );
    // update sort state
    state.guests.filters.sort = document.querySelector('#filter-sort').value;

    // reset page to 1 when applying filters
    state.guests.filters.page = 1;

    emitter.emit('guests:fetch');
  });
}
