import { state } from '../../core/state.js';

import { renderLoading, renderError } from '../../shared/render.js';
import { getPaginationItems } from '../../shared/pagination.js';

import { getCountryFlag } from './helpers.js';

// renderFilters displays controls for filtering and sorting guests.
function renderFilters() {
  return `
    <div id="filters">
      <div class="filters-header">
        <h4>Filters</h4>
        <button id="toggle-filters">Toggle</button>
      </div>

      <div class="controls ${state.guests.ui.filtersOpen ? '' : 'hidden'}">
        <div class="text-search">
          <div class="text-group">
            <label for="filter-name">Name</label>
            <input id="filter-name" type="text" placeholder="e.g. Smith" value="${state.guests.filters.name}" />
          </div>
          <div class="text-group">
            <label for="filter-country">Country Code</label>
            <input id="filter-country" type="text" maxlength="2" placeholder="e.g. BZ" value="${state.guests.filters.country}" />
          </div>
        </div>

        <div class="select-group">
          <label for="filter-page-size">Items Per Page</label>
          <select id="filter-page-size">
            <option value="1" ${state.guests.filters.page_size == 1 ? 'selected' : ''}>1</option>
            <option value="6" ${state.guests.filters.page_size == 6 ? 'selected' : ''}>6</option>
            <option value="12" ${state.guests.filters.page_size == 12 ? 'selected' : ''}>12</option>
          </select>
        </div>

        <div class="select-group">
          <label for="filter-sort">Sort</label>
          <select id="filter-sort">
            <option value="name" ${state.guests.filters.sort === 'name' ? 'selected' : ''}>Name (Ascending)</option>
            <option value="-name" ${state.guests.filters.sort === '-name' ? 'selected' : ''}>Name (Descending)</option>
            <option value="passport_number" ${state.guests.filters.sort === 'passport_number' ? 'selected' : ''}>Passport (Ascending)</option>
            <option value="-passport_number" ${state.guests.filters.sort === '-passport_number' ? 'selected' : ''}>Passport (Descending)</option>
          </select>
        </div>

        <button id="apply-filters">Apply</button>
      </div>
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
    <section class="guests-section">
      <div class="guest-grid">
        ${state.guests.data
          .map(
            (g) => `
              <div class="guest-card">
                <div class="guest-header">
                  <h3>${g.name}</h3>
                  <span class="guest-country">${getCountryFlag(g.country)} ${g.country}</span>
                </div>

                <div class="guest-body">
                  <p><strong>Gender:</strong> ${g.gender}</p>
                  <p><strong>Passport:</strong> ${g.passport_number}</p>
                  <p><strong>Email:</strong> ${g.contact_email}</p>
                  <p><strong>Phone:</strong> ${g.contact_phone}</p>
                </div>

                <div class="guest-footer">
                  <small>${g.street}, ${g.city}</small>
                </div>
              </div>
            `,
          )
          .join('')}
      </div>
    </section>
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

  return `
  <div class="pagination-wrapper">
    <div id="pagination">
      <button class="page-btn" data-page="${Math.max(1, current - 1)}">←</button>

      <div class="page-numbers">
        ${items
          .map((item) => {
            if (item === '...') return `<span class="dots">...</span>`;

            return `
            <button
              class="page-btn ${item === current ? 'active' : ''}"
              data-page="${item}"
              ${item === current ? 'disabled' : ''}
            >
              ${item}
            </button>
          `;
          })
          .join('')}
      </div>

      <button class="page-btn" data-page="${Math.min(total, current + 1)}">→</button>
    </div>
  </div>
`;
}

// render composes the html for guests.html.
export function render() {
  const app = document.querySelector('#app');

  let resultsContent = '';

  // choose the last component rendered
  if (state.guests.loading) {
    resultsContent = renderLoading('Loading...');
  } else if (state.guests.error) {
    resultsContent = renderError(state.guests.error);
  } else {
    resultsContent = renderGuests();
  }

  app.innerHTML = `
    <h2>Guests</h2>
    ${renderFilters()}
    ${renderPagination()}
    ${resultsContent}
  `;
}
