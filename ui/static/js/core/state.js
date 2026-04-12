// state holds all the data that may be retrieved from data service layer.
export const state = {
  // guest entity
  guests: {
    data: [],
    metadata: {},
    filters: {
      page: 1,
      page_size: 6,
      name: '',
      country: '',
      sort: 'name',
    },

    ui: {
      filtersOpen: false,
    },

    loading: false,
    error: null,
  },
};
