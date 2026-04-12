// getCountryFlag uses the ISO 3166-1-alpha-2 country codes to convert to
// the appropriate flag emoji.
export function getCountryFlag(code) {
  if (!code) return '🌍';

  return code
    .toUpperCase()
    .replace(/./g, (char) => String.fromCodePoint(127397 + char.charCodeAt()));
}

// validateGuestFilters performs frontend validation of guest filter values.
export function validateGuestFilters(filters) {
  // country code check
  const country = filters.country.trim();
  if (country && !/^[a-zA-Z]{2}$/.test(country)) {
    return 'Country code must contain exactly 2 letters.';
  }

  // page size validation
  const MAX_PAGE_SIZE = 60;
  const MIN_PAGE_SIZE = 1;
  if (
    Number.isNaN(filters.page_size) ||
    filters.page_size < MIN_PAGE_SIZE ||
    filters.page_size > MAX_PAGE_SIZE
  ) {
    return `Page size must be between ${MIN_PAGE_SIZE} and ${MAX_PAGE_SIZE}.`;
  }

  // sort validation (whitelist)
  const ALLOWED_SORTS = [
    'name',
    '-name',
    'passport_number',
    '-passport_number',
  ];
  if (!ALLOWED_SORTS.includes(filters.sort)) {
    return 'Invalid sort option selected.';
  }

  return null; // valid
}
