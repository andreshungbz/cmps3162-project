import { emitter } from './event-emitter.js';
import { state } from './state.js';

const API_BASE = '/v1';

// DataService is the layer that interacts with the server API to interact
// with the database.
export const DataService = {
  // fetchGuests retrieves multiple guest records.
  // Authorization: Front Desk
  async fetchGuests(payload) {
    try {
      // get URL parameters
      const params = new URLSearchParams(state.guests.filters);

      // fetch guests
      const res = await fetch(`${API_BASE}/guests?${params}`, {
        method: 'GET',
        headers: {
          Authorization: `Bearer UZPHDJXKVJORSLLQD7AH3QDCU3`, // test token
        },
      });
      if (!res.ok) {
        throw new Error(`Server Error: ${res.status} ${res.statusText}`);
      }

      // read JSON data
      const data = await res.json();

      // emit event
      emitter.emit('guests:fetched', data);
    } catch (err) {
      emitter.emit('guests:error', err.message);
    }
  },
};
