import { state } from '../../core/state.js';

import { getPaginationItems } from '../../shared/pagination.js';

// renderLoading displays loading text.
function renderLoading() {
  if (!state.guests.loading) return '';
  return `<div class="loading-box">Loading guests...</div>`;
}

// renderError displays an error message from the server.
function renderError() {
  if (!state.guests.error) return '';
  return `<div class="error-box">${state.guests.error}</div>`;
}

// renderFilters displays controls for filtering and sorting guests.
function renderFilters() {
  return `
    <div id="filters">
      <input id="filter-name" placeholder="Name" value="${state.guests.filters.name}" />
      <input id="filter-country" placeholder="Country" value="${state.guests.filters.country}" />

      <select id="filter-page-size">
        <option value="1" ${state.guests.filters.page_size == 1 ? 'selected' : ''}>1</option>
        <option value="5" ${state.guests.filters.page_size == 5 ? 'selected' : ''}>5</option>
        <option value="10" ${state.guests.filters.page_size == 10 ? 'selected' : ''}>10</option>
      </select>

      <select id="filter-sort">
        <option value="name" ${state.guests.filters.sort === 'name' ? 'selected' : ''}>
          Name (Ascending)
        </option>
        <option value="-name" ${state.guests.filters.sort === '-name' ? 'selected' : ''}>
          Name (Descending)
        </option>
        <option value="passport_number" ${state.guests.filters.sort === 'passport_number' ? 'selected' : ''}>
          Passport (Ascending)
        </option>
        <option value="-passport_number" ${state.guests.filters.sort === '-passport_number' ? 'selected' : ''}>
          Passport (Descending)
        </option>
      </select>

      <button id="apply-filters">Apply</button>
    </div>
  `;
}

// renderGuests displays guests in an unordered list.
function renderGuests() {
  // check loading or error
  if (state.guests.loading || state.guests.error) return '';

  // no entries
  if (state.guests.data.length === 0) {
    return `<p>No guests found.</p>`;
  }

  return `
    <ul>
      ${state.guests.data.map((g) => `<li>${g.name} (${g.country})</li>`).join('')}
    </ul>
  `;
}

// renderPagination displays controls for navigating to different pages of a guests list.
function renderPagination() {
  // check state object
  const meta = state.guests.metadata;
  if (!meta?.last_page) return '';

  // construct pagination items
  const current = meta.current_page;
  const total = meta.last_page;
  const items = getPaginationItems(current, total);

  let html = '';

  // add previous button
  html += `<button data-page="${Math.max(1, current - 1)}">←</button>`;

  // add middle buttons, including for current
  for (const item of items) {
    if (item === '...') {
      html += `<span>...</span>`;
      continue;
    }

    html += `
      <button data-page="${item}" class="${item === current ? 'active' : ''}" ${item === current ? 'disabled' : ''}>
        ${item}
      </button>
    `;
  }

  // add next button
  html += `<button data-page="${Math.min(total, current + 1)}">→</button>`;

  return `<div id="pagination">${html}</div>`;
}

// render composes the html for guests.html.
export function render() {
  const app = document.querySelector('#app');

  app.innerHTML = `
    <h2>Guests</h2>
    ${renderFilters()}
    ${renderLoading()}
    ${renderError()}
    ${renderGuests()}
    ${renderPagination()}
  `;
}
