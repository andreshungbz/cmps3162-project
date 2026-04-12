// state holds all the data that may be retrieved from data service layer.
export const state = {
  // guest entity
  guests: {
    data: [],
    metadata: {},
    filters: {
      page: 1,
      page_size: 5,
      name: '',
      country: '',
      sort: 'name',
    },

    loading: false,
    error: null,
  },
};
