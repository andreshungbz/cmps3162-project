import { DataService } from '../../core/data-service.js';
import { emitter } from '../../core/event-emitter.js';
import { state } from '../../core/state.js';

import { render } from './render.js';
import { setupHandlers } from './handlers.js';

// guests:fetch event -> call DataService
emitter.on('guests:fetch', () => {
  state.guests.loading = true;
  state.guests.error = null;
  render();

  DataService.fetchGuests();
});

// guests:fetched event -> update guests data and metadata then rerender
emitter.on('guests:fetched', (data) => {
  state.guests.loading = false;
  state.guests.data = data.guests;
  state.guests.metadata = data.metadata;
  render();
});

// guests:pageChanged event -> update guests filters page then call DataService again
emitter.on('guests:pageChanged', (page) => {
  state.guests.filters.page = page;
  emitter.emit('guests:fetch');
});

// guests:error event -> set error message
emitter.on('guests:error', (error) => {
  state.guests.loading = false;
  state.guests.error = error;
  render();
});

// INITIALIZATION

setupHandlers();
emitter.emit('guests:fetch');
