import { emitter } from '../../core/event-emitter.js';
import { state } from '../../core/state.js';

import { validateGuestFilters } from './helpers.js';
import { render } from './render.js';

// setupHandlers applies event listeners for guests pagination and filters.
export function setupHandlers() {
  const app = document.querySelector('#app');

  // Pagination Buttons (event delegation)
  app.addEventListener('click', (e) => {
    const btn = e.target.closest('[data-page]');
    if (!btn) return;

    emitter.emit('guests:pageChanged', Number(btn.dataset.page));
  });

  // Filters Toggle Button
  app.addEventListener('click', (e) => {
    const toggleBtn = e.target.closest('#toggle-filters');
    if (toggleBtn) {
      state.guests.ui.filtersOpen = !state.guests.ui.filtersOpen;
      render();
    }
  });

  // Filters Apply Button
  app.addEventListener('click', (e) => {
    if (e.target.id !== 'apply-filters') return;

    const newFilters = {
      name: document.querySelector('#filter-name').value,
      country: document.querySelector('#filter-country').value,
      page_size: Number(document.querySelector('#filter-page-size').value),
      sort: document.querySelector('#filter-sort').value,
    };

    // run validation
    const error = validateGuestFilters(newFilters);

    if (error) {
      state.guests.error = error;
      render();
      return;
    }

    // clear previous errors
    state.guests.error = null;

    // commit filters
    state.guests.filters = {
      ...state.guests.filters,
      ...newFilters,
      page: 1,
    };

    emitter.emit('guests:fetch');
  });
}
