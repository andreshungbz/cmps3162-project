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

  // hotel image showcase
  showcase: {
    data: [],

    uploadForm: {
      file: null,
      filter: 'grayscale',
      isOpen: true,
    },

    ui: {
      formOpen: false,
      modalOpen: false,
      selectedJob: null,
    },

    loading: false,
    error: null,
  },
};
